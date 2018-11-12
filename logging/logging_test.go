package logging

import (
	"testing"
)

func TestGetLogger(t *testing.T) {
	logger := Get("TEST")
	if logger == nil {
		t.Fatal("Could not get a logger")
	}
}

func TestSetLogLevel(t *testing.T) {
	err := SetLoggingLevel("panic")
	if err != nil {
		t.Error("Could not succesfully set log level panic")
	}
	err = SetLoggingLevel("fatal")
	if err != nil {
		t.Error("Could not succesfully set log level fatal")
	}
	err = SetLoggingLevel("error")
	if err != nil {
		t.Error("Could not succesfully set log level error")
	}
	err = SetLoggingLevel("warn")
	if err != nil {
		t.Error("Could not succesfully set log level warn")
	}
	err = SetLoggingLevel("warning")
	if err != nil {
		t.Error("Could not succesfully set log level warning")
	}
	err = SetLoggingLevel("info")
	if err != nil {
		t.Error("Could not succesfully set log level info")
	}
	err = SetLoggingLevel("debug")
	if err != nil {
		t.Error("Could not succesfully set log level debug")
	}

	err = SetLoggingLevel("blub")
	if err == nil {
		t.Error("Setting invalid log level should have failed")
	}
}
