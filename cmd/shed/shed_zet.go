package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"time"

	cli "github.com/urfave/cli/v3"
)

var (
	secondBrain string
	title       string
	subDir      string
)

var (
	vimCmd  = "nvim"
	vimInit = "~/.config/nvim/init.lua"
)

var defaultSubDir = "0-inbox"

var fileTemplate = `---
id: %s
aliases: []
tags:
  - change-me
date: "%s"
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
			Name:        "title",
			Aliases:     []string{"f"},
			Value:       "",
			Usage:       "the name of the file (note) to create",
			Destination: &title,
		},
	},
	Action: func(ctx context.Context, cmd *cli.Command) error {
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

		var noteTitle string
		var noteID string
		var err error

		path := fmt.Sprintf("%s/%s", secondBrain, subDir)

		if title != "" {
			noteTitle = title
			noteID = createSlug(noteTitle)

			if _, statErr := os.Stat(fmt.Sprintf("%s/%s.md", path, noteID)); statErr == nil {
				fmt.Printf("Error: Note with ID '%s' already exists.\n", noteID)
				os.Exit(1)
			}
		} else {
			noteTitle, noteID, err = getUniqueNoteDetails(path)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		}

		if _, err := os.Stat(path); os.IsNotExist(err) {
			if err := os.MkdirAll(path, 0o755); err != nil {
				return err
			}
		}

		fullFilePath := fmt.Sprintf("%s/%s.md", path, noteID)

		fmt.Println("Creating note:", noteID)

		initialFileContents := fmt.Sprintf(
			fileTemplate,
			noteID,
			time.Now().Format("2006-01-02"),
			noteTitle,
		)

		if err := os.WriteFile(fullFilePath, []byte(initialFileContents), 0o644); err != nil {
			return err
		}

		vimCmd := exec.Command(
			vimCmd,
			"-u",
			vimInit,
			"+normal G",
			"+startinsert!",
			fullFilePath,
		)

		vimCmd.Stdin = os.Stdin
		vimCmd.Stdout = os.Stdout

		if err := vimCmd.Run(); err != nil {
			log.Fatal(err)
		}

		return nil
	},
}

func getUniqueNoteDetails(path string) (string, string, error) {
	for {
		title, err := promptForTitle()
		if err != nil {
			return "", "", err
		}

		slug := createSlug(title)

		fullPath := fmt.Sprintf("%s/%s.md", path, slug)
		if _, err := os.Stat(fullPath); os.IsNotExist(err) {
			return title, slug, nil
		}

		fmt.Printf("A note with the ID '%s' already exists. Please choose a different title.\n", slug)
	}
}

func promptForTitle() (string, error) {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("Note Title: ")

		input, err := reader.ReadString('\n')
		if err != nil {
			return "", err
		}

		input = strings.TrimSpace(input)

		if input == "" {
			fmt.Println("Title cannot be empty")
			continue
		}

		return input, nil
	}
}

func createSlug(s string) string {
	s = strings.ToLower(s)

	reg := regexp.MustCompile(`[^a-z0-9\s]+`)
	s = reg.ReplaceAllString(s, "")

	s = strings.ReplaceAll(s, " ", "-")

	regDash := regexp.MustCompile(`-+`)
	s = regDash.ReplaceAllString(s, "-")

	return strings.Trim(s, "-")
}
