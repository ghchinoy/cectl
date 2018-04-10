package main

import (
	"fmt"
)

const (
	version     = "0.14.0"
	versionName = "stratocumulus"
)

// Version returns the version of the app as a string
func Version() string {
	return fmt.Sprintf("%s %s", version, versionName)
}

/* in 1000s of ft
lenticular 		20  -130
cirrus 			20  - 39
cirrostratus 	20  - 39
cirrocumulus 	18  - 20
altostratus 	15  - 20
altocumulus 	7.9 - 20
cumulonumbus 	4.9 - 50
cumulus 		2   - 9.8
stratocumulus 	1.5 - 6.6
nimbostratus 	0   - 9.8
stratus 		0   - 6.8
*/
