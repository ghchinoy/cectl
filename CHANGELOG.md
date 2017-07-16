# CHANGELOG

## 20170711 0.3.0

* Retry Formula Instance Execution via `executions retry <execution_id>`
* Common Resources: `resources` - with `resources list`
* Profiles broken out into `profiles list` and `profiles add` from simply `profiles` (which listed existing profiles, only)

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
