package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"time"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"

	"github.com/manifoldco/promptui"
	"github.com/urfave/cli/v2"
)

var (
	secondBrain string
	filename    string
	subDir      string
)

var vimCmd = "vim"

var defaultSubDir = "0-inbox"

var fileTemplate = `---
date: %s
tags:
  - 
---

# %s 


`

var zetCommand = &cli.Command{
	Name:        "zet",
	Description: "",
	Usage:       "Creates a note in the Zettelkasten system.",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:        "second-brain",
			Aliases:     []string{"b"},
			Value:       "",
			Usage:       "the name of the file (note) to create",
			Destination: &secondBrain,
		},
		&cli.StringFlag{
			Name:        "sub-dir",
			Aliases:     []string{"d"},
			Value:       "",
			Usage:       "the name of the file (note) to create",
			Destination: &subDir,
		},
		&cli.StringFlag{
			Name:        "file-name",
			Aliases:     []string{"f"},
			Value:       "",
			Usage:       "the name of the file (note) to create",
			Destination: &filename,
		},
	},
	Action: func(cCtx *cli.Context) error {
		if secondBrain == "" {
			secondBrainEnv, ok := os.LookupEnv("SECOND_BRAIN")
			if ok {
				secondBrain = secondBrainEnv
			} else {
				fmt.Println("Please set the SECOND_BRAIN environment variable or use the -b flag")
				os.Exit(1)
			}
		}

		if subDir == "" {
			subDir = defaultSubDir
		}

		if filename == "" {
			var err error

			path := fmt.Sprintf("%s/%s", secondBrain, subDir)

			filename, err = promptUniqueDashedFileName(path)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		}

		if _, err := os.Stat(fmt.Sprintf("%s/%s", secondBrain, subDir)); os.IsNotExist(err) {
			err := os.MkdirAll(fmt.Sprintf("%s/%s", secondBrain, subDir), 0o755)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		}

		fullFilePath := fmt.Sprintf("%s/%s/%s.md", secondBrain, subDir, filename)

		initialFileContents := fmt.Sprintf(
			fileTemplate,
			time.Now().Format("2006-01-02"),
			titleCase(strings.ReplaceAll(filename, "-", " ")),
		)

		if err := os.WriteFile(fullFilePath, []byte(initialFileContents), 0o644); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		cmd := exec.Command(vimCmd, "+normal G", "+startinsert!", fullFilePath)

		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout

		err := cmd.Run()
		if err != nil {
			log.Fatal(err)
		}

		return nil
	},
}

func titleCase(s string) string {
	containsIgnoreCase := func(s []string, str string) bool {
		for _, v := range s {
			if strings.EqualFold(v, str) {
				return true
			}
		}
		return false
	}

	excludedConjunctions := []string{
		"but",
		"and",
		"nor",
		"or",
		"for",
		"so",
		"as",
		"if",
		"yet",
		"to",
	}

	titleCaser := cases.Title(language.English)

	words := strings.Fields(s)

	for i, word := range words {
		if containsIgnoreCase(excludedConjunctions, word) {
			words[i] = strings.ToLower(word)
		} else {
			words[i] = titleCaser.String(word)
		}
	}

	return strings.Join(words, " ")
}

func promptUniqueDashedFileName(path string) (string, error) {
	var val string

	validate := func(input string) error {
		fullPath := fmt.Sprintf("%s/%s.md", path, input)

		if _, err := os.Stat(fullPath); err == nil {
			return errors.New("file already exists")
		}

		reg := regexp.MustCompile(`^[A-Za-z0-9-]+$`)
		if !reg.MatchString(input) {
			return errors.New("invalid filename")
		}

		val = input

		return nil
	}

	s := promptui.Prompt{
		Label:     "Filename",
		Validate:  validate,
		AllowEdit: true,
	}

	_, err := s.Run()
	if err != nil {
		return "", err
	}

	return val, nil
}
