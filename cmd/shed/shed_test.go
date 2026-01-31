package main

import (
	"bytes"
	"testing"
	"context"

	"github.com/stretchr/testify/assert"
	cli "github.com/urfave/cli/v3"
)

func testApp() *cli.Command {
	tapp := *app
	return &tapp
}

func AppRunTest(t *testing.T) {
	var out, stderr bytes.Buffer

	shed := testApp()

	shed.Writer = &out
	shed.ErrWriter = &stderr

	err := shed.Run(context.Background(), []string{
		"shed",
		"-h",
	})
	assert.NoError(t, err)

	output := out.String()
	assert.Contains(t, output, "shed - The Shed toolbox.")
	assert.NotContains(t, output, "error")
	assert.NotContains(t, output, "panic")
	
	assert.Empty(t, stderr.String()) 
}
