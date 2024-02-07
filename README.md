# shed

A dev toolbox for The Shed.

This is a monolith cli tool written in Go, containing everything I think I need to enhance my productivity.

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Report](https://goreportcard.com/badge/github.com/shanehull/shed)](https://goreportcard.com/report/github.com/shanehull/shed)

[![Test Workflow](https://github.com/shanehull/shed/actions/workflows/test.yaml/badge.svg)](https://github.com/shanehull/shed/actions/workflows/test.yaml/badge.svg)
[![Lint Workflow](https://github.com/shanehull/shed/actions/workflows/lint.yaml/badge.svg)](https://github.com/shanehull/shed/actions/workflows/lint.yaml/badge.svg)
[![Release Workflow](https://github.com/shanehull/shed/actions/workflows/release.yaml/badge.svg?branch=main)](https://github.com/shanehull/shed/actions/workflows/release.yaml/badge.svg?branch=main)

The tool itself is probably no use to you, but due to the way [`urfave/cli`](https://github.com/urfave/cli) defines commands, you can just copy/paste any command you like and use it in your own project.

They are all located [here](https://github.com/shanehull/shed/tree/main/cmd/shed) (e.g. `shed_awesomecmd.go`) and imported [here](https://github.com/shanehull/shed/blob/main/cmd/shed/main.go).
