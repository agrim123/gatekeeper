# gatekeeper

Authentication and authorization oriented tool allowing non-root users to ssh to a machine without giving them access to private keys.

## Architecture

The roles are delegated inside the gatekeeper for various tasks. The hierarchy is:
```
gatekeeper
 |__ guard
 |   |__ authentication
 |   |__ authorization
 |__ runtime
 |   |__ executes based on type  ->--status is returned-->|
 |__ notifier         <-----------------------------------|
     |__ defaults to stdout
```

Gatekeeper is reponsible for calling `guard`, `runtime` and `notifier`. After executing the requested instrctions, the returned status is then notified to the users via `notifier` module. Gatekeeper initializes all the three to default when it is initialized,
```golang
&GateKeeper{
    ctx:      ctx,
    runtime:  runtime.NewRuntime(ctx),
    notifier: notifier.NewDefaultNotifier(),
    guard:    guard.NewGuard(ctx),
    store:    store.Store,
}
```

The guard is responsible for authentication and authorizing the user and the command the user is requesting.

After the guard verifies the user, the command is passed to runtime for execution. The required action is taken based on the type of command.

After execution, whether success or failure, the status is returned to the gatekeeper, which then calls the notifier to inform the user of the result.

Every step is focused to be pluggable to provide ease of integrating your methods.

A detailed architecture
```
User triggers
+-------------+
| plan.option |
+-------------+
      ↓ (1)
+------------+       Loads store           +-------+
| gatekeeper | <-------------------------> | store |-------------------------------------|
+------------+        on startup           +-------+                                     |
  |   ↓ (2)                                                                              |
  |-+-------+                                                                            |
  | | guard |__________________    (fails if root user is running)                       |
  | +---                       |                  |                      +----+          |
  |     | authentication (3) <-|-- fetches user executing the command -> | OS |          |
  |     |       ↓              |                                         +----+          |
  |     | authorization  (4) <-|----------------------------------------------------------
  |     +-------|--------------+        checks which all plan's options are allowed
  |             |
  |             |
  |             | (5)
  |             | plan.option is finally
  |             ↓ sent to runtime to be executed
  |---------+---------+
  |         | runtime |
  |         +---------+
  |             | (6)
  |             ↓ status is sent to notifier
  |---------+----------+
            | notifier |
            +----------+
```

### Pluggable modules

Gatekeeper provides basic authentication, authorization, and notifier (default is stdout) modules. However, this can easily be customized by adding your methods and passing them to the gatekeeper after initialization.

```golang
  gatekeeper := NewGatekeeper(context.Background())
  gatekeeper.WithNotifier(MyCustomNotifier)
```

#### Notifier Module

Default notifier module logs to stdout. However, it can entirely be customized by creating you own module and injecting it to gattekeeper on initialization. 

`SlackNotifier` is also present but disabled by default. It can be enabled by using:
```golang
    gatekeeper := NewGatekeeper(context.Background())
    gatekeeper.WithNotifier(notifier.NewSlackNotifier("<SLACK_WEBHOOK_URL>"))
```

If any notifer fails, the default behaviour is to dump logs to stdout, so that you don't miss out any logs.
```bash
[SUCCESS]  | Authenticated as agrim
[SUCCESS]  | Authorized `agrim` to perform `service1 shell`
[INFO]     | Executing plan: service1 shell
[INFO]     | Spawning shell for <user>@<host>
[INFO] 🔐  | Reading private key
[ERROR]    | Notifier: slack failed. Fallback to default notifier
[NOTIFIER] | Plan `service1 shell` executed by `agrim` failed. Error: Failed to connect to <user>@<host>. Error: dial tcp 3.84.241.53:22: i/o timeout
```

## Setup

Four configs drive gatekeeper:
- `users.json` [Sample](examples/configs/users.json)
    - The system users are to be given access to a particular resource.
    - This is the first and foremost config that is loaded and used to authenticate users.
    - Every user belongs to some groups, which in turn are allowed to run only a subset of commands.
- `groups.json` [Sample](examples/configs/groups.json)
    - Groups are the ACL for the gatekeeper.
    - Every group has a set of `allowed_plans` that the user belonging to that group can execute.
    - This is crucial to the authorization step.
    - Privileged groups:
        - Group `*` defines root privileges. This group has access to every plan and can run any option.
        - Group `plan.*` gives access to all options of that plan.
    - Usernames are mapped to system users, so this gives us an extra security layer.

Since gatekeeper is entirely relying on authentication and authorization of user running the command, the `groups.json` and `users.json` are critical configurations to gatekeeper's working.

- `plan.json` [Sample](examples/configs/plan.json)
    - Plan can be considered as the master config that defines what all commands are available to users.
    - It is a JSON file with the `plan` key as an array of what we call **plans**.
    - Every plan has a key **name** which is the identifier of that plan.
    - [Options](#supported-options):
        - Each plan has a set of options, with a key as an identifier and a field **type**, to take the required action when the option is called.
    - Example Usage:
    ```bash
    $ gatekeeper run-plan plan1 option1 # This gives us custom command-line options
    ```
- `servers.json` [Sample](examples/configs/servers.json)
    - When doing something on remote instances, this config is responsible for storing the config of ssh hosts, including hostname, port, private key.
    - Each server has a set of instances that contain the username, IP, and absolute private key path required to ssh into the instance.

A little side note: gatekeeper cannot be run by the root user. Instead, we run the gatekeeper binary using `+s`.

### [Almost] Secured private keys

The main goal of the gatekeeper is to run some commands on or provide access to a server without handing out private keys to all the users.
The ideal situation is to put all keys on the bastion server and have users access the required server (if they are allowed) via gatekeeper.

We use `chmod +s gatekeeper` so that the non-root user executing the binary, can use (not access, not read) the protected private key on behalf of the root user.

## Supported options

Options as identified by **type**, available options are:

### `local`

For running commands on local system. Can be useful if user doesn't have permission to execute certain commands, and can run only those without giving any other access.

```json
"options": {
    "some_cmd_name": {
        "type": "local",
        "stages": [
            "ls"
        ]
    },
}
```

Note: Here `some_cmd_name` is the command that can be provided to user to run from cli. Options are actually identified by **type**.

### `shell`

Spawns a pseudo shell for the given server.

```json
"options": {
    "some_cmd_name": {
        "type": "shell",
        "server": "service1-server"
    },
}
```

### `remote`

Runs commands on a remote server. Can be useful to trigger deploy commands without giving ssh access to user.

```json
"options": {
    "some_cmd_name": {
        "type": "remote",
        "server": "service1-server",
        "stages": [
            "ls -a",
            "/usr/bin/whoami",
            "echo \"Hello from remote server\""
        ]
    },
}
```

### `container`

Can be used to spawn a docker container. Container is flexible enough to do anything, mount a volume, build something, open a remote shell or run commands etc.

```json
"options": {
    "some_cmd_name": {
        "type": "container",
        "server": "service1-server",
        "protected": false,
        "volumes": {
            "/host/path/to/volume": "/container/path/to/mount/to"
        },
        "stages": [
            {
                "command": [
                    "ssh",
                    "-i",
                    "/home/deploy/keys/service1.pem",
                    "ec2-user@host",
                    "ls -a"
                ],
                "privileged": false
            }
        ]
    }
}
```

This by default mounts the provided **server** private key to container. (this is yet to be fixed).

## Examples

Checkout usage of gatekeeper [here](https://github.com/agrim123/gatekeeper-cli).

A sample run of gatekeeper
```bash
$ gatekeeper run-plan service1 shell
[SUCCESS]  | Authenticated as agrim
[SUCCESS]  | Authorized `agrim` to perform `service1 shell`
[INFO]     | Executing plan: service1 shell
[INFO]     | Spawning shell for <user>@<host>
[INFO] 🔐  | Reading private key
[INFO]     | Shell Spawned. Press Ctrl+C to exit.
<spawned shell>
```

### Future prospects

Gatekeeper is not limited to only providing shell access, it can be used to run deploy commands, as a proxy intermediary, currently, the config is entirely file-based but can be extended to a database for easy updates and more observability, can be used to run restricted commands on the local system which otherwise unprivileged user cannot run and many more.

## TODO

- [ ] see infra health (read-only)
- [ ] self update via git
    - [ ] https://github.com/go-git/go-git
- [ ] build image
- [ ] Log every ssh interaction
- [ ] Remove container support
- [ ] Check if private keys exist beforehand
