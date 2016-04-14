package main

import (
	"os"
	"testing"
)

func TestMargeXlsx1(t *testing.T) {
	os.Args = []string{
		"marge-xlsx",
		"-vv",
	}

	app := NewApp()
	if err := app.Run(os.Args); err != nil {
		t.Error(err)
	}
}
