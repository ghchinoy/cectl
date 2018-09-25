# (next version)

NEW FEATURES:

BUG FIXES:

IMPROVEMENTS:

# v0.17.2

NEW FEATURES:

* Code of Conduct added ([#14](https://github.com/ghchinoy/cectl/issues/14))

BUG FIXES:

* `jobs delete --max` now checks to determine whether user-provided max is greater than existing jobs ([#73](https://github.com/ghchinoy/cectl/issues/73))

IMPROVEMENTS:


# v0.17.1

NEW FEATURES:

BUG FIXES:

IMPROVEMENTS:

* limit job deletion concurrency or max count with new flags, `curr` | `r` and `max` | `m` respectively

# v0.17.0

NEW FEATURES:

* Create a job with `jobs create --file <job.json>` given a JSON config file
* Delete a job with `jobs delete <jobid>`
* Delete all jobs with `jobs delete --all`

BUG FIXES:

IMPROVEMENTS:

# v0.16.0

NEW FEATURES:

* `profiles init` creates a new, blank cectl.toml config file if the file doesn't exist, in the default location
* `profiles add -l` or `profiles add --login` begins an interactive login session to create a cectl.toml file

BUG FIXES:

* `profiles list -l` now lists all profiles rather than a truncated list

IMPROVEMENTS:

# v0.15.3

NEW FEATURES:

BUG FIXES:

IMPROVEMENTS:

* Updates to latest ce-go client which shows user privileges with roles  in JSON output via `user list -r --json`

# v0.15.2

NEW FEATURES:

* Hidden `lbdocs` command for `elements` to retrieve LoopBack models for an Element

BUG FIXES:

IMPROVEMENTS:

# v0.15.1

NEW FEATURES:

* exposes `branding set` with only the `--file` flag feature

BUG FIXES:

IMPROVEMENTS:

* spelling and err handling corrections

# v0.15.0

NEW FEATURES:

* `instances delete` now takes either a list of IDs (1 or more) to delete multiple instances or the flag `--all` to delete all of them. Will provide an error if Element Instance is in use or not found.

BUG FIXES:

* `molecules export` variants, continued fix to exporting raw JSON ([PR#51](https://github.com/ghchinoy/cectl/pull/51) 

IMPROVEMENTS:

* `profiles env` now also adds exports for `CE_ORG` and `CE_USER` in addition to existing `CE_BASE` and `CE_AUTH`

# v0.14.0

NEW FEATURES:

* `profiles env` adds the ability to output environment variables for current profile, `CE_BASE` url and `CE_AUTH` which can be used in a curl header; can also do `source <(cectl profiles env)` to set env variables

BUG FIXES:

IMPROVEMENTS:

* CircleCI added
* Malformed Formulas are displayed in `formulas list` with `N/A` for success step, trigger type, and API ([#2](https://github.com/ghchinoy/cectl/issues/2))
* `formulas list` now shows config count per Formula ([#9](https://github.com/ghchinoy/cectl/issues/9))

# v0.13.0

NEW FEATURES:

BUG FIXES:

IMPROVEMENTS:

* `instances test` now has a flag to remove bad instances, `--remove` ([#1](https://github.com/ghchinoy/cectl/issues/1))

# v0.12.1

NEW FEATURES:

* Convenience output method for CE ROI calculator via flag `--roi` for `elements list` command, outputs JSON in a particular format

BUG FIXES:

IMPROVEMENTS:

# v0.12.0

NEW FEATURES:

* Initial release of CE UI Branding management via `branding` command - implemented `get` and `reset`, with `set` available, but hidden, as there's more work to do ([#10](https://github.com/ghchinoy/cectl/issues/10))

BUG FIXES:

* Listing of Formual Instance Executions now enforces requiring of an ID (`executions list <id>`) ([#40](https://github.com/ghchinoy/cectl/issues/40))

IMPROVEMENTS:

* Output of `formula-instances create` now shows created Formula Instance ID ([#36](https://github.com/ghchinoy/cectl/issues/36), [#10](https://github.com/ghchinoy/cectl/issues/10))

# v0.11.0

NEW FEATURES:

BUG FIXES:

* Exports of Elements, Transformations, Formulas, and Resource export raw object so as not to conflict with changes to the Platform API, rather than relying on cectl/ce-go's structures ([#26](https://github.com/ghchinoy/cectl/issues/26))
* Formula imports weren't working as expected. Updated Formula structure ([#32](https://github.com/ghchinoy/cectl/issues/32))

IMPROVEMENTS:

# v0.10.0

NEW FEATURES:

* `instances enable <instance_id>` - to enable an Element Instance
* `instances disable <instance_id>` - to disable an Element Instance
* `instances events-enable <instance_id> [true|false]` - a single command that enables or disables an Element Instances events, given a true|false (bool) option

BUG FIXES:

IMPROVEMENTS:

# v0.9.0

NEW FEATURES:

* Adds combined vdr (common object + transforms) file with `--combined` flag for `molecules export`; speed of `--combined` flag is a known issue

BUG FIXES:

* Corrects csv flag outputting only profile data in `profiles list --csv`

IMPROVEMENTS:

# v0.8.1

NEW FEATURES:

* Adds csv output for profiles and transformations with `--csv` flag

BUG FIXES:

IMPROVEMENTS:

# v0.8.0

NEW FEATURES:

BUG FIXES:

* Critical fix in Formula definition, improper capitialization of configuration keys results in invalid Formula json

IMPROVEMENTS:

# v0.7.2

NEW FEATURES:

BUG FIXES:

IMPROVEMENTS:

* A new flag to provide more details on the definition of a Resource, `resources definition <resource> -d` 

# v0.7.1

NEW FEATURES:

* `resources copy <resource> <new_resource>` with optional flag `--deep` to also associate new resource Transformations with Element; addresses [#16](https://github.com/ghchinoy/cectl/issues/16) 
* `transformations delete <resource> <element>`, to delete a Transformation association from an Element, which is a prerequisite if deleting a Resource (existing capability in `resources delete <resource>`
* incomplete `transformations associate <resource> <element>` as there's still more discussion ([#20](https://github.com/ghchinoy/cectl/issues/20)) 

BUG FIXES:

IMPROVEMENTS:

* copyright statement added to transformations.go source file
* deleting an resource would continue even if it resulted in a non-200 response, resulting in a confusing output 

## v0.7.0

NEW FEATURES

* `transformations` top level command, with one subcommand, `transformations list` ([#17](https://github.com/ghchinoy/cectl/issues/17))
* New transformations subcommands, `elements transformations <id>` and `instances transformations <id>` command, to list the Transformations associated with an Element ID/key or Element Instance ID, respectively ([#17](https://github.com/ghchinoy/cectl/issues/17))
* `molecules export` now exports Transformations per Element into a file in the created `transformations` directory ([#17](https://github.com/ghchinoy/cectl/issues/17))
* `profiles list` has a new flag to output a table of profiles available, `-l` or `--long`

BUG FIXES

IMPROVEMENTS

* isolated CE API to [`ce-go`](https://github.com/ghchinoy/ce-go) library ([#8](https://github.com/ghchinoy/cectl/issues/8))
* using [dep-style changelog format](https://github.com/golang/dep) going forwards
* Export Resources and Formulas methods are here now, instead of within [ce-go]((https://github.com/ghchinoy/ce-go))

## 20180126 0.6.1

* `molecules export [formulas|resources|all (default)]` initial version; `molecules` hidden

## 20180125 0.5.0

* Element import help text, shows help
* `resources add` hidden (still active, but deprecated)
* `resources import` replaces `resources add`

## 20180121 0.4.9.1

* restoring removed `formula-instances create <id> [name] [--configuration <configfile>]`

## 20180121 0.4.9

* user info table display correction incorporated
* display of empty resources list corrected
* `resources delete <name>` - allows deleting of common resource objects
* `elements import <element.json>` - imports an Element json to the Platform
* `instances delete <id>` - deletes an Element Instance
* `formula-instances delete <id>` - deletes a Formula Instance
* separation of CE go library to [ce-go](https://github.com/ghchinoy/ce-go)

## 20180121 0.4.8

* introduced `instances test` which hits the `/ping` endpoint for each Instance and reports non-200

## 20180120 0.4.7.1

* using [ce-go](https://github.com/ghchinoy/ce-go) library

## 20171128 0.4.7

* `elements list` now has `--csv` csv output flag

## 20171001 0.4.6.1

* minor changes to info output

## 20170930 0.4.6

* rudimentary info command

## 20170926 0.4.5

* added `elements export` to export full element model

## 20170823 0.4.4

* Changed Formula Instance Exection list limit flag from `--num` to `--top` to make it clear that the return is from the latest
* internal, consolidation of `execution` commands
* `instances docs [id]` requires an integer ID for the Element Instance ID

## 20170817 0.4.3

* Instance details added `instance details <instanceID>` returns details about an Instance including token
* Instance docs added `instance docs <instanceID>` returns OAI Spec for Instance
* Details for Instances of a Formula now provides configuration params for the Instance `formulas instances <formula_ID>`
* Added "extendable" to `elements list` output
* Added an optional key param to `elements list [keyfilter]`, which will return only the Elements with keys matching `keyfilter`

## 20170810 0.4.2 

* Formula commands now honor `--curl` flag
* Formula Instances create command now takes a the configuration json with the `--configuration` flag, in addition deprecated`--instance`/`--i` flags
* Users list via `users list`

## 20170806 0.4.1

* clarification to help for `elements list --order` options
* combining profile commands into one file, `profiles.go`
* Elements OAI retrieval via ID or Key `elements docs <ID|key>`
* Elements metadata retrieval via ID or Key `elements metadata <ID|key>` - formatted json output
* Elements instance retrieval via ID or Key `elements instances [ID|key]` - optional ID/key, will do an `instances list` if no ID/key provided
* Instances list `instances list` is available

## 20170805 0.4.0

* Initial version of adding and defining individual Resources - `resources list | define <name> | add <name> <json>`

## 20170719 0.3.0

* Retry Formula Instance Execution via `executions retry <execution_id>`
* Common Resources: `resources` - with `resources list`
* Profiles broken out into `profiles list` (list existing profiles), `profiles add` (add a new profile), and `profiles set` (set default profile) from simply `profiles` (which listed existing profiles, only) - note, still doesn't create a config file if one does not already exist
* Added `hubs`, with `hubs list`
* Activate and deactivate a Formula `formula activate`, `formula deactivate`

## 20170707 0.2.0

* `formulas list` now returns the API field for manual triggered APIs, with presence indicating a Formula as a Resource

## 20170509 0.1.0

*  `formulas list` with `--json`/`-j` json output option
  * honors `--profile` flag
  * also shows count of instances
* `formulas details` - outputs the basic info, triggers, and steps of a Formula template
* `formulas import TEMPLATEPATH` - imports a json TEMPLATEPATH Formula template
*  `formula-instances list` - lists all instances
* `formula-instances trigger ID [-d]` - triggers an instance with data object given in `-d`
* added `version` which outputs version
* a variety of stub commands added (`eb`, `elements`, `profiles`)
* `executions` is the top-level command for managing formula instance executions
