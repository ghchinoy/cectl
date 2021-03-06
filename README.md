# Cloud Elements Control

Command-line interface for Cloud Elements' Platform APIs.

[![CircleCI](https://circleci.com/gh/ghchinoy/cectl.svg?style=shield)](https://circleci.com/gh/ghchinoy/cectl) [![Go Report Card](https://goreportcard.com/badge/github.com/ghchinoy/cectl)](https://goreportcard.com/report/github.com/ghchinoy/cectl)

```
Cloud Elements control is a CLI that manages the platform

Usage:
  cectl [command]

Available Commands:
  branding          Manage Branding of the Platform
  elements          Manage Elements on the Platform
  executions        Manage Formula Instance Executions
  formula-instances Manage Formula Instances
  formulas          Manage formulas on the platform
  help              Help about any command
  hubs              Hub management
  info              Information about this account
  instances         Manage Instances of Elements on the Platform
  jobs              Manage jobs on the platform
  profiles          Manage profiles
  resources         Manage common resources
  transformations   Manage Transformations on the Platform
  users             Manage users on the platform
  version           version of cectl

Flags:
  -h, --help   help for cectl

Use "cectl [command] --help" for more information about a command.
```


## Install

OS X, Linux, and Windows binaries are available.

### Option 1 - Use a package manager (preferred)

**On OS X, with Homebrew**

```
brew update
brew install ghchinoy/ce/cectl
```

Then, to initialize a blank configuration file and then create a profile via a login:

```
cectl profiles init
cectl profiles add --login
```

**On Windows, with [`scoop`](http://scoop.sh/)**


Add the bucket:

```
scoop bucket add ce https://github.com/ghchinoy/scoop-ce
scoop bucket list
scoop search cectl
```

Install `cectl`:

```
scoop install cectl
```

Then, to initialize a blank configuration file and then create a profile via a login:

```
cectl profiles init
cectl profiles add --login
```


Or you can manually create a config file needed to use cectl - see below for format.

Example, shown using `nano` (which can be installed via `scoop install nano`):

```
new-item -path ~\.config\ce -type directory
nano ~\.config\ce\cectl.toml
```

Refer to [scoop.sh](http://scoop.sh/) for more info on scoop.


### Option 2 - Download a release from GitHub

View the [releases](https://github.com/ghchinoy/cectl/releases) page for the `cectl` GitHub project and find the appropriate archive for your operating system and architecture.  For OS X systems, remember to use the Darwin archive.

### Option 3 - Build from source

If you have a Go environment configured, install from source like so:

```
go get github.com/ghchinoy/cectl
go install github.com/ghchinoy/cectl
```

# Config file

`cectl` expects a valid configuration file in [TOML](https://github.com/toml-lang/toml) format. The configuration file's default location is in `$HOME/.config/ce/cectl.toml`.

You can use `profiles init` to create a blank cectl.toml file in the default location and then `profiles add --login` to add a profile to the configuration with a Cloud Elements login.

Example config file:

```
[default]
base="https://api.cloud-elements.com/elements/api-v2"
user="USER-HASH-HERE"
org="ORG-HASH-HERE"

[snapshot]
base="https://snapshot.cloud-elements.com/elements/api-v2"
user="USER-HASH-HERE"
org="ORG-HASH-HERE"
```

You may define multiple CE environment targets with different TOML blocks.

Utilize profiles by adding the profile flag, ex. `--profile snapshot`
