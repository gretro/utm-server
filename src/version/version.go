package version

import _ "embed"

//go:embed version.txt
var versionDescribe string

func AppVersion() string {
	return versionDescribe
}
