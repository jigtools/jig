package log

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"time"

	bugsnaglogrus "github.com/Shopify/logrus-bugsnag"
	logrus "github.com/Sirupsen/logrus"
	"github.com/SvenDowideit/jig/config"
	"github.com/SvenDowideit/jig/showuserlog"
	"github.com/SvenDowideit/jig/util"
	bugsnag "github.com/bugsnag/bugsnag-go"
)

var logFile *os.File

func StopLogging() {
	if logFile != nil {
		fmt.Errorf("CLOSING logfile pid %d\n", os.Getpid())
		// TODO: work ou thow to disable the bugsnag log too
		logFile.Close()
		logFile = nil
	}
	logrus.SetOutput(os.Stderr)
	logrus.Debugf("Stopped logging to file")
}

func InitLogging(logLevel logrus.Level, version string) {
	// TODO: i'm trusting that no-one has messed with it since we last made it..
	if _, err := os.Stat(config.LogDir); err != nil {
		if err := util.SudoRun("mkdir", "-p", config.LogDir); err != nil {
			logrus.Fatal(err)
		}
		if err := util.SudoRun("chmod", "777", config.LogDir); err != nil {
			logrus.Fatal(err)
		}
	}
	
	if logrus.StandardLogger().Out != os.Stderr {
		fmt.Printf("IDK\n")
		logrus.Debugf("IDK")
	}

	// Write all levels to a log file
	logrus.SetLevel(logrus.DebugLevel)
	if logFile == nil {
		filename := filepath.Join(config.LogDir, "verbose-"+time.Now().Format("2006-01-02T15.04-")+strconv.Itoa(os.Getpid())+".log")
		fmt.Printf("Debug log written to %s\n", filename)
		f, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Unable to open %s log file, %v", filename, err)
			logrus.Fatal(err)
		}
		logFile = f
	}
	logrus.SetOutput(logFile)

	// Filter what the user sees (info level, unless they set --debug)
	showuserHook, err := showuserlog.NewShowuserlogHook(logLevel)
	if err != nil {
		logrus.Fatal(err)
	}
	logrus.StandardLogger().Hooks.Add(showuserHook)

	// We'll get a bugsnag entry for Error, Fatal and Panic
	bugsnag.Configure(bugsnag.Configuration{
		APIKey:       "ad1003e815853e3c15d939709618d50e",
		AppVersion:   version,
		ReleaseStage: "initial",
		Synchronous:  true,
	})
	bugsnagHook, err := bugsnaglogrus.NewBugsnagHook()
	if err != nil {
		logrus.Fatal(err)
	}
	logrus.StandardLogger().Hooks.Add(bugsnagHook)

	pwd, err := os.Getwd()
	if err != nil {
		logrus.Fatal(err)
	}
	logrus.Debugf("START: %v in %s", os.Args, pwd)
}
