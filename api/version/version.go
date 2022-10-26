package version

import (
	"fmt"
	
	"admincheckapi/api/resource"
)

var (
	Version, Build, Revision string
)

func Set(version, build, revision string) {
	Version, Build, Revision = version, build, revision
}

//
// info composes version infor string using Version/Branch/Revision
// variables compiled into the binary executable
//
func OneLineInfo() string {
	return fmt.Sprintf("%s-%s-%s", Version, Build, Revision)
}

//
// Level provides formatted version info
//
func Level() resource.Version {
	v := resource.Version{
		Version: OneLineInfo(),
	}
	
	return v
}
