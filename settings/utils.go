package settings

import (
	"os"
	"path/filepath"
	"regexp"

	"github.com/teamhephy/workflow-cli/executable"
)

var filepathRegex = regexp.MustCompile(`^.*[/\\].+\.json$`)

var EnvName = executable.Env()

func locateSettingsFile(cf string) string {
	if cf == "" {
		if v, ok := os.LookupEnv(EnvName); ok {
			cf = v
		} else {
			cf = "client"
		}
	}

	// if path appears to be a filepath (contains a separator and ends in .json) don't alter the path
	if filepathRegex.MatchString(cf) {
		return cf
	}

	return filepath.Join(FindHome(), executable.Config(), cf+".json")
}
