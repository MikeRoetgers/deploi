# deploi

[![Build Status](https://travis-ci.org/MikeRoetgers/deploi.svg?branch=master)](https://travis-ci.org/MikeRoetgers/deploi)

deploi is a tool that manages the deployment of build artifacts into different cluster environments.

## Basic Operation

deploi consist of three parts: The **deploid** daemon, a number of **agents** and a **client**, that is used by your CI system as well as humans to interact with deploid.

### deploid

The `deploid` daemon is responsible for managing state. All data is stored within [boltdb](https://github.com/boltdb/bolt), so there is no need to run a DBMS or something comparable to operate `deploid`. Every new build relevant for `deploid` is announced to it by the build system (or alternatively manually through the client).

### deploi-agent

On every environment you want to use deploi to deploy artifacts, you have to run a `deploi-agent`. The `deploi-agent` pulls in regular intervals new jobs from `deploid` and executes them. The choice for pull instead of push was made to circumvent potential firewall and/or private network issues. After a job is completed, the agent reports back to `deploid` and the job is marked as done.

### deploi

The `deploi` client is used to interact with `deploid`. Build systems will use the client to announce new builds to the daemon. Humans will use the client to manually create deployments or define automations.

## Client Usage

After installing the client on a system, it should be initialized with:

```
$ deploi setup config --host deploid.my.domain:3375
```

This command creates a fresh config in `~/.deploi`. Afterwards you can login via `deploi login my@username.tld`. A successful login results in the daemon issuing a JWT token and sending it back to the client. The token is stored in the config file, so make sure to protect it appropriately (by default it is created in your home directory with permissions only for your user).

`deploi help` gives you a list of all available commands.

## Development

Dependencies are managed via [dep](https://github.com/golang/dep).
