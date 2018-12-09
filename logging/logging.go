// Package logging provides logging facilities for ALoLStats
package logging

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/sirupsen/logrus"
)

// Init the logging framework
// has to be called only once
func Init() {
	logrus.SetFormatter(new(myFormatter))
	logrus.SetLevel(logrus.DebugLevel)
	logrus.SetOutput(os.Stdout)
}

// Get a logger with prefix name
func Get(name string) *logrus.Entry {
	return logrus.WithField("name", name)
}

// SetLoggingLevel takes one of the strings
// panic, fatal, error, warn/warning, info or debug
// and sets the log level accordingly
func SetLoggingLevel(loggingLevel string) error {
	level, err := logrus.ParseLevel(loggingLevel)
	if err != nil {
		Get("logging").Warnln("Error setting log level to", loggingLevel)
		return err
	}

	logrus.SetLevel(level)
	Get("logging").Infoln("Setting log level to", loggingLevel)
	return nil
}

// SetLogFile enables logging to a log file in addition to stdout
func SetLogFile(logFile string) error {
	file, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err == nil {
		mw := io.MultiWriter(os.Stdout, file)
		logrus.SetOutput(mw)
		Get("logging").Infoln("Setting log file to", logFile)
		return nil
	} else {
		Get("logging").Warnln("Failed to log to file, using stdout only")
		logrus.SetOutput(os.Stdout)
		return err
	}
}

type myFormatter struct{}

func (f *myFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	b := &bytes.Buffer{}
	name, ok := entry.Data["name"]
	if !ok {
		name = "default"
	}
	fmt.Fprintf(b, "%s [%-5.5s] (%s): %s\n", entry.Time.Format("2006-01-02 15:04:05.000"), strings.ToUpper(entry.Level.String()), name, entry.Message)
	return b.Bytes(), nil
}
