# gatekeeper

Authentication and authorization oriented tool for allowing users to ssh a remote machine without giving them access to private keys.

## Setup

- The main file which drives gatekeeper goes into `configs/plan.json`. Have a look at sample `plan.json`.
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
- Since, gatekeeper is entirely relies on authentication and authorization of user running the command, the `roles.json` and `users.json` are  critical configurations to gatekeeper working.
    - Every user is assigned some roles which in turn have `allowed_plans` that can be run by users in that role.
        - The allowed plans format is `plan.option`, for example, user_service.deploy.
        - Special cases:
            - Role `*` defines root privileges. This role has access to every plan, and can run any option.
            - Role `plan.*` gives access to all options of the plan.
    - Usernames are mapped to system users, so this gives an extra security layer, since user cannot spoof its own system user name.
- IMP: gatekeeper cannot be run by root user. Instead we run the gatekeeper binary using `+s`.

### [Almost] Secured private keys

The main goal of gatekeeper is to run some commands on or provide access to a server without handing out private keys to all the users.
Ideal situation is to put all keys on bastion server, and have users access required server (if they are allowed) via gatekeeper.

We use `chmod +s gatekeeper` so that the non root user executing the binary, can use (not access, not read) the protected private key on behalf of the root user.

## TODO

- [ ] see infra health (read-only)
- [ ] self update via git
    - [ ] https://github.com/go-git/go-git
- [ ] build image
- [ ] Log every ssh interaction
- [ ] Remove container support