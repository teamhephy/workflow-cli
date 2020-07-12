package parser

import (
	docopt "github.com/docopt/docopt-go"
	"github.com/teamhephy/workflow-cli/cmd"
	"github.com/teamhephy/workflow-cli/executable"
)

// Version displays the client version
func Version(argv []string, cmdr cmd.Commander) error {
	usage := executable.Render(`
Displays the client version.

Usage: {{.Name}} version [options]

Options:
  -a --all
    list api and controller versions
`)

	args, err := docopt.Parse(usage, argv, true, "", false, true)
	if err != nil {
		return err
	}

	return cmdr.Version(args["--all"].(bool))
}
