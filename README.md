# gatekeeper

Authentication and authorization oriented tool allowing non-root users to ssh to a machine without giving them access to private keys.

## Setup

- The main file which drives gatekeeper goes into `configs/plan.json`. Have a look at sample [plan.json](https://github.com/agrim123/gatekeeper-cli/blob/master/configs/plan-sample.json).
    - It is a json file with `plan` key as an array of what we call **plans**.
    - Every plan has a key **name** which is the identifier of that plan.
    - Each plan has a set of options, with key as identifier and a field **type**, to take the required action when the option is called.
    - Usage:
    ```bash
    $ gatekeeper run-plan plan1 option1 # This gives you custom command line options
    ```
- Next important file is `servers.json`.
    - Each server has a set of instances which contain the username, ip and private key path required to ssh into the instance.
    - The keys can be put in relative path to gatekeeper binary in `keys` folder.
- Since, gatekeeper is entirely relies on authentication and authorization of user running the command, the `groups.json` and `users.json` are  critical configurations to gatekeeper working.
    - Every user is assigned some groups which in turn have `allowed_plans` that can be run by users in that group.
        - The allowed plans format is `plan.option`, for example, user_service.deploy.
        - Special cases:
            - Group `*` defines root privileges. This group has access to every plan, and can run any option.
            - Group `plan.*` gives access to all options of the plan.
    - Usernames are mapped to system users, so this gives us an extra security layer.
- IMP: gatekeeper cannot be run by root user. Instead we run the gatekeeper binary using `+s`.

### [Almost] Secured private keys

The main goal of gatekeeper is to run some commands on or provide access to a server without handing out private keys to all the users.
Ideal situation is to put all keys on bastion server, and have users access required server (if they are allowed) via gatekeeper.

We use `chmod +s gatekeeper` so that the non root user executing the binary, can use (not access, not read) the protected private key on behalf of the root user.

## Architecture

The roles are delegated inside gatekeeper for various tasks. The hierarchy is:
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
		runtime:  runtime.NewDefaultRuntime(),
		notifier: notifier.GetNotifier(),
		guard:    guard.NewGuard(),
	}
```

The guard is responsible for authentication and authorizing the user and the command it the user is requesting.

After guard verifies the user, the command is passed to runtime for execution. The required action is taken based on type of command.

After execution, whether success or failure, the status is returned back to gatekeeper, which then calls the notifier to inform the user of the result.

Every step is focused to be pluggable to provide ease of integrating your own methods.

### Pluggable modules

Gatekeeper provides basic authentication, authorization and notifier (default is stdout) modules. However, this can easily be customized by adding your own methods and passing them to gatekeeper after initialization initialization.

```golang
  gatekeeper := NewGatekeeper(context.Background())
  gatekeeper.WithNotifier(MyCustomNotifier)
```

### Examples

Checkout usage of gatekeeper [here](https://github.com/agrim123/gatekeeper-cli).

A sample run of gatekeeper
```bash
$ gatekeeper run-plan service1 shell
[SUCCESS]  | Authenticated as agrim
[SUCCESS]  | Authorized `agrim` to perform `service1 shell`
[INFO]     | Executing plan: service1 shell
[INFO]     | Spawning shell for <user>@<host>
[INFO] 🔐  | Reading private key
<spawned shell>
```

### Future prospects

Gatekeeper is not limited to only providing shell access, it can used to run deploy commands, as a proxy intermediary, currently the config is entirely file based but can be extended to a database for easy updates and more observability, can be used to run restricted commands on the local system which otherwise unprivileged user cannot run and many more.

## TODO

- [ ] see infra health (read-only)
- [ ] self update via git
    - [ ] https://github.com/go-git/go-git
- [ ] build image
- [ ] Log every ssh interaction
- [ ] Remove container support
- [ ] Check if private keys exist beforehand
