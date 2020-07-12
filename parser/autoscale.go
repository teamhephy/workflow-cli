package parser

import (
	docopt "github.com/docopt/docopt-go"
	"github.com/teamhephy/workflow-cli/cmd"
	"github.com/teamhephy/workflow-cli/executable"
)

// Autoscale displays all relevant commands for `hephy autoscale`.
func Autoscale(argv []string, cmdr cmd.Commander) error {
	usage := executable.Render(`
Valid commands for autoscale:

autoscale:list   list autoscale options of an application
autoscale:set    turn on autoscale for an app
autoscale:unset  turn off autoscale for an app

Use '{{.Name}} help [command]' to learn more.
`)

	switch argv[0] {
	case "autoscale:list":
		return autoscaleList(argv, cmdr)
	case "autoscale:set":
		return autoscaleSet(argv, cmdr)
	case "autoscale:unset":
		return autoscaleUnset(argv, cmdr)
	default:
		if printHelp(argv, usage) {
			return nil
		}

		if argv[0] == "autoscale" {
			argv[0] = "autoscale:list"
			return autoscaleList(argv, cmdr)
		}

		PrintUsage(cmdr)
		return nil
	}
}

func autoscaleList(argv []string, cmdr cmd.Commander) error {
	usage := executable.Render(`
Prints a list of autoscale options for the application.

Usage: {{.Name}} autoscale:list [options]

Options:
  -a --app=<app>
    the uniquely identifiable name for the application.
`)

	args, err := docopt.Parse(usage, argv, true, "", false, true)

	if err != nil {
		return err
	}

	return cmdr.AutoscaleList(safeGetValue(args, "--app"))
}

func autoscaleSet(argv []string, cmdr cmd.Commander) error {
	usage := executable.Render(`
Set autoscale option per process type for an app.

Usage: {{.Name}} autoscale:set <process-type> --min=<min> --max=<max> --cpu-percent=<percent> [options]

Arguments:
  <process-type>
    the process type to add to the application's autoscale settings.
  --min=<min>
	minimum replicas to keep around
  --max=<max>
	max replicas to scale up to
  --cpu-percent=<cpu-percent>
	target CPU utilization

Options:
  -a --app=<app>
    the uniquely identifiable name of the application.
`)

	args, err := docopt.Parse(usage, argv, true, "", false, true)

	if err != nil {
		return err
	}

	processType := args["<process-type>"].(string)
	app := safeGetValue(args, "--app")
	min := safeGetInt(args, "--min")
	max := safeGetInt(args, "--max")
	CPUPercent := safeGetInt(args, "--cpu-percent")

	return cmdr.AutoscaleSet(app, processType, min, max, CPUPercent)
}

func autoscaleUnset(argv []string, cmdr cmd.Commander) error {
	usage := executable.Render(`
Unset autoscale per process type for an app.

Usage: {{.Name}} autoscale:unset <process-type> [options]

Arguments:
  <process-type>
    the process type to remove from the application's autoscale settings.

Options:
  -a --app=<app>
    the uniquely identifiable name of the application.
`)

	args, err := docopt.Parse(usage, argv, true, "", false, true)

	if err != nil {
		return err
	}

	processType := args["<process-type>"].(string)
	app := safeGetValue(args, "--app")

	return cmdr.AutoscaleUnset(app, processType)
}
