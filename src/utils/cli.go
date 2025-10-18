package utils

import (
	"flag"
	"os"
)

type CLIArgs struct {
	Local     bool
	Install   bool
	Docker    bool
	Stage     string
	Project   string
	ExtraArgs []string
}

func ParseArgs() CLIArgs {
	args := CLIArgs{}
	flags := flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	flags.BoolVar(&args.Local, "l", false, "Run locally")
	flags.BoolVar(&args.Local, "local", false, "Run locally")
	flags.BoolVar(&args.Local, "d", false, "Run in Docker")
	flags.BoolVar(&args.Local, "docker", false, "Run in Docker")
	flags.BoolVar(&args.Install, "i", false, "Install dependencies")
	flags.BoolVar(&args.Install, "install", false, "Install dependencies")

	flags.Parse(os.Args[1:])
	remaining := flags.Args()
	if len(remaining) == 0 {
		PrintUsage()
		os.Exit(2)
	}
	args.Stage = remaining[0]
	if len(remaining) > 1 {
		args.Project = remaining[1]
	}
	if len(remaining) > 2 {
		args.ExtraArgs = remaining[2:]
	}
	return args
}
