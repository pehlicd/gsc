/*
Copyright © 2024 Furkan Pehlivan <furkanpehlivan34@gmail.com>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/

package main

import (
	"errors"
	"flag"
	"fmt"
	"github.com/jwalton/gchalk"
	"github.com/pehlicd/gsc/internal"
	"github.com/pehlicd/gsc/internal/git"
	"github.com/pehlicd/gsc/internal/gitlab"
	"github.com/pehlicd/gsc/internal/logger"
	"github.com/rs/zerolog"
	"os"
)

const helpMsg = `
gsc - GitLab Structured Cloner

gsc is a tool to help you clone all the repositories from a Gitlab group in a way that you see the repositories in the same structure as they are in the Gitlab group.

usage: gsc [flags]
`

var (
	versionString string
	buildDate     string
	buildCommit   string

	red = gchalk.Stderr.Red

	fl       = flag.NewFlagSet("gsc", flag.ContinueOnError)
	host     = fl.String("host", "https://gitlab.com", "GitLab hostname, default is https://gitlab.com")
	username = fl.String("username", "", "GitLab username for authentication")
	token    = fl.String("token", "", "GitLab token for authentication")
	insecure = fl.Bool("insecure", false, "Allow insecure connection to your GitLab instance, default is false")

	all         = fl.Bool("all", true, "Clone all projects, default is true")
	concurrency = fl.Int("concurrency", 10, "Number of concurrent workers, default is 10")
	group       = fl.Int("group", 0, "Clone projects from the given group ID")
	matcher     = fl.String("matcher", "", "Clone projects that match the given regex")
	recursive   = fl.Bool("recursive", false, "Clone projects recursively, default is false")
	verbose     = fl.Bool("verbose", false, "Verbose output, default is false")
	version     = fl.Bool("version", false, "Print version information")

	logLevel zerolog.Level
)

func printErrAndExit(msg string) {
	_, _ = fmt.Fprintf(os.Stderr, "%s%s\n", red("error: "), msg)
	os.Exit(1)
}

func main() {
	if err := fl.Parse(os.Args[1:]); err != nil {
		if errors.Is(err, flag.ErrHelp) {
			fmt.Printf("%s\n", helpMsg)
			os.Exit(0)
		}
		printErrAndExit(err.Error())
	}

	if *version {
		fmt.Printf("gsc version: %s\n", versionString)
		fmt.Printf("Build date: %s\n", buildDate)
		fmt.Printf("Build commit: %s\n", buildCommit)
		os.Exit(0)
	}

	if *verbose {
		logLevel = zerolog.DebugLevel
	} else {
		logLevel = zerolog.InfoLevel
	}

	Application := internal.Application{
		Auth: &internal.Auth{
			Host:     host,
			Username: username,
			Token:    token,
			Insecure: insecure,
		},
		All:         all,
		Concurrency: concurrency,
		Group:       group,
		Matcher:     matcher,
		Recursive:   recursive,
	}

	log := logger.NewLogger(logLevel)
	log.Debug().Msgf("recursive: %t", *recursive)
	Application.Log = log

	client, err := gitlab.NewClient(Application.Auth)
	if err != nil {
		log.Fatal().Err(err).Msg("couldn't create GitLab client")
	}

	Application.Client = client

	log.Debug().Msg("GitLab client successfully created")

	projects, err := gitlab.GetGroupProjects(gitlab.Application{Application: Application})
	if err != nil {
		log.Fatal().Err(err).Msgf("couldn't get projects for group with ID %d", *group)
	}

	if err := git.Clone(
		git.Application{Application: Application},
		projects,
	); err != nil {
		log.Fatal().Err(err).Msg("couldn't clone projects")
	}
}
