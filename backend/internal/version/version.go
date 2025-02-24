package version

import (
	"fmt"
	"strconv"
	"strings"
)

type ReleasModes string

const (
	Release ReleasModes = "release"
	Dev     ReleasModes = "dev"
)

type Version struct {
	Major int
	Minor int
	Patch int
}

func (v *Version) String() string {
	return fmt.Sprintf("%d.%d.%d", v.Major, v.Minor, v.Patch)
}

func NewVersion(versionString string) (*Version, error) {
	if versionString == "" {
		return nil, fmt.Errorf("back version is empty")
	}

	versionStringSeperated := strings.Split(versionString, ".")
	if len(versionStringSeperated) != 3 {
		return nil, fmt.Errorf("version %s is wrong", versionString)
	}
	versionIntSeperated := []int{}
	for _, numberString := range versionStringSeperated {
		number, err := strconv.Atoi(numberString)
		if err != nil {
			return nil, err
		}
		versionIntSeperated = append(versionIntSeperated, number)

	}

	newVersion := &Version{
		Major: versionIntSeperated[0],
		Minor: versionIntSeperated[1],
		Patch: versionIntSeperated[2],
	}

	return newVersion, nil
}
