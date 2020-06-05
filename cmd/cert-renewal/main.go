package main

import (
	"fmt"
	"os"

	"github.com/jduepmeier/certrenewal"

	"github.com/jessevdk/go-flags"
	"github.com/sirupsen/logrus"
)

type opts struct {
	Config    string `short:"c" long:"config" default:"config.yaml" description:"path to config"`
	Insecure  bool   `short:"k" long:"insecure" description:"do not validate ca-certificate"`
	LogFormat string `short:"f" long:"log-format" description:"output format for logging"`
	LogOutput string `short:"o" long:"log-output" description:"output file for logging (- is stderr)"`
	LogLevel  string `short:"l" long:"log-level" description:"level to log"`
}

func run() (returnCode int, err error) {
	opts := opts{
		LogFormat: "text",
		LogOutput: "-",
		LogLevel:  "INFO",
	}

	_, err = flags.Parse(&opts)
	if err != nil {
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

	return certrenewal.Run(config)
}

func main() {
	returnCode, err := run()
	if err != nil {
		logrus.Error(err)
	}
	os.Exit(returnCode)
}
