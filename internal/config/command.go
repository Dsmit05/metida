package config

import (
	"flag"
	"fmt"
	"os"
)

// CommandLine contains all params from os.Args.
type CommandLine struct {
	prod    bool // debag = true in mode dev
	logPath string
}

// NewCommandLine responsible for the launch settings of this project.
func NewCommandLine() (*CommandLine, error) {
	var prodMod bool
	var logPath string

	if len(os.Args[1:]) < 1 {
		CommandHelp()
		return &CommandLine{}, fmt.Errorf("The startup mode is not set")
	}

	modeProd := flag.NewFlagSet("prod", flag.ExitOnError)
	modeProd.StringVar(&logPath, "logPath", "logs.json", "file logging")

	modeDev := flag.NewFlagSet("dev", flag.ExitOnError)
	// Todo: fixed, it's not in use right now:
	modeDev.StringVar(&logPath, "logPath", "logs.json", "file logging")

	switch os.Args[1] {
	case "prod":
		prodMod = true

		if err := modeProd.Parse(os.Args[2:]); err != nil {
			fmt.Println("Arguments prod error")
		}
	case "dev":
		prodMod = false

		if err := modeDev.Parse(os.Args[2:]); err != nil {
			fmt.Println("Argument dev error")
		}

	case "help":
		CommandHelp()
		return &CommandLine{}, fmt.Errorf("select the startup mode")
	}

	return &CommandLine{prod: prodMod, logPath: logPath}, nil
}

// IfDebagOn return true if mod = dev.
func (o *CommandLine) IfDebagOn() bool {
	return !o.prod
}

// GetLogPath return file name for logs.
func (o *CommandLine) GetLogPath() string {
	return o.logPath
}

// CommandHelp shows information about flags.
func CommandHelp() {
	fmt.Println(`
List of main commands:
	dev: [development mode]
		flags:

	prod: [production mode]
		flags:
		-logPath [name log file] (string)
`)
}
