package parser

import (
	docopt "github.com/docopt/docopt-go"
	"github.com/teamhephy/workflow-cli/cmd"
	"github.com/teamhephy/workflow-cli/executable"
)

// Labels displays all relevant commands for `hephy label`.
func Labels(argv []string, cmdr cmd.Commander) error {
	usage := executable.Render(`
Valid commands for labels:

labels:list   list application's labels
labels:set    add new application's label
labels:unset  remove application's label

Use '{{.Name}} help [command]' to learn more.
`)

	switch argv[0] {
	case "labels:list":
		return labelsList(argv, cmdr)
	case "labels:set":
		return labelsSet(argv, cmdr)
	case "labels:unset":
		return labelsUnset(argv, cmdr)
	default:
		if printHelp(argv, usage) {
			return nil
		}

		if argv[0] == "labels" {
			argv[0] = "labels:list"
			return labelsList(argv, cmdr)
		}

		PrintUsage(cmdr)
		return nil
	}
}

func labelsList(argv []string, cmdr cmd.Commander) error {
	usage := executable.Render(`
Prints a list of labels of the application.

Usage: {{.Name}} labels:list [options]

Options:
  -a --app=<app>
    the uniquely identifiable name for the application.
`)

	args, err := docopt.Parse(usage, argv, true, "", false, true)

	if err != nil {
		return err
	}

	return cmdr.LabelsList(safeGetValue(args, "--app"))
}

func labelsSet(argv []string, cmdr cmd.Commander) error {
	usage := executable.Render(`
Sets labels for an application.

A label is a key/value pair used to label an application. This label is a general information for {{.Name}} user.
Mostly used for administration/maintenance information, note for application. This information isn't send to scheduler.

Usage: {{.Name}} labels:set [options] <key>=<value>...

Arguments:
  <key> the label key, for example: "git_repo" or "team"
  <value> the label value, for example: "https://github.com/teamhephy/workflow" or "frontend"

Options:
  -a --app=<app>
    the uniquely identifiable name for the application.
`)

	args, err := docopt.Parse(usage, argv, true, "", false, true)
	if err != nil {
		return err
	}

	app := safeGetValue(args, "--app")
	tags := args["<key>=<value>"].([]string)

	return cmdr.LabelsSet(app, tags)
}

func labelsUnset(argv []string, cmdr cmd.Commander) error {
	usage := executable.Render(`
Unsets labels for an application.

Usage: {{.Name}} labels:unset [options] <key>...

Arguments:
  <key> the label key to unset, for example: "git_repo" or "team"

Options:
  -a --app=<app>
    the uniquely identifiable name for the application.
`)

	args, err := docopt.Parse(usage, argv, true, "", false, true)
	if err != nil {
		return err
	}

	app := safeGetValue(args, "--app")
	tags := args["<key>"].([]string)

	return cmdr.LabelsUnset(app, tags)
}
