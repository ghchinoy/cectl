# CHANGELOG

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
