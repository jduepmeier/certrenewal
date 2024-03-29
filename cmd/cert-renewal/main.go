package main

import (
	"fmt"
	"io"
	"os"

	"github.com/jduepmeier/certrenewal"

	"github.com/jessevdk/go-flags"
	"github.com/sirupsen/logrus"
)

var build = "dev"

type opts struct {
	Config    string `short:"c" long:"config" default:"config.yaml" description:"path to config"`
	Insecure  bool   `short:"k" long:"insecure" description:"do not validate ca-certificate"`
	LogFormat string `short:"f" long:"log-format" description:"output format for logging"`
	LogOutput string `short:"o" long:"log-output" description:"output file for logging (- is stderr)"`
	LogLevel  string `short:"l" long:"log-level" description:"level to log"`
	Version   bool   `short:"v" long:"version" description:"show version and exit"`
}

func run(args []string, stdout io.Writer, stderr io.Writer) (returnCode int, err error) {
	opts := opts{
		LogFormat: "text",
		LogOutput: "-",
		LogLevel:  "INFO",
	}

	_, err = flags.Parse(&opts)
	if err != nil {
		return returnCode, nil
	}

	if opts.Version {
		fmt.Printf("%s - %s\n", os.Args[0], build)
		return returnCode, nil
	}

	if opts.LogFormat == "json" {
		logrus.SetFormatter(&logrus.JSONFormatter{})
	}

	if opts.LogOutput != "-" {
		logfile, err := os.OpenFile(opts.LogOutput, os.O_APPEND|os.O_CREATE, 0640)
		if err != nil {
			return returnCode, fmt.Errorf("cannot open logfile: %s", err)
		}
		defer logfile.Close()

		logrus.SetOutput(logfile)
	} else {
		logrus.SetOutput(stderr)
	}

	level, err := logrus.ParseLevel(opts.LogLevel)
	if err != nil {
		return returnCode, err
	}

	logrus.SetLevel(level)

	config, err := certrenewal.ReadConfig(opts.Config)
	if err != nil {
		return returnCode, err
	}

	if opts.Insecure {
		config.Insecure = opts.Insecure
	}

	return certrenewal.Run(config)
}

func main() {
	returnCode, err := run(os.Args, os.Stdout, os.Stderr)
	if err != nil {
		logrus.Error(err)
	}
	os.Exit(returnCode)
}
