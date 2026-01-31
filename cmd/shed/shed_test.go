package main

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	cli "github.com/urfave/cli/v2"
)

func testApp() *cli.App {
	tapp := *app
	return &tapp
}

func AppRun_Test(t *testing.T) {
	var out, err bytes.Buffer

	shed := testApp()

	shed.Writer = &out
	shed.ErrWriter = &err

	// check app launches up correctly
	assert.NoError(t, shed.Run([]string{
		"dcf",
		"-h",
	}))

	// check that the output looks more or less correct
	assert.Contains(t, out.String(), "shed - The Shed toolbox.")
	assert.NotContains(t, out.String(), "error")
	assert.NotContains(t, out.String(), "panic")
	assert.Empty(t, err.String())
}
