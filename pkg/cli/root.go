package cli

import (
	"errors"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var Version = "dev"

var rootCmd = &cobra.Command{
	Use:     "scout",
	Short:   "🪖 Scout turns Daily Papers into deterministic intelligence briefings",
	Long:    rootLong(),
	Version: Version,
}

func init() {
	rootCmd.SetVersionTemplate("scout version {{.Version}}\n")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		var exitErr *cliExitError
		if errors.As(err, &exitErr) {
			if !exitErr.silent {
				fmt.Fprintln(os.Stderr, exitErr.Error())
			}
			os.Exit(exitErr.code)
		}
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

type cliExitError struct {
	err    error
	code   int
	silent bool
}

func newCLIExitError(err error, code int, silent bool) *cliExitError {
	if code == 0 {
		code = 1
	}
	return &cliExitError{err: err, code: code, silent: silent}
}

func (e *cliExitError) Error() string {
	if e == nil || e.err == nil {
		return ""
	}
	return e.err.Error()
}

func rootLong() string {
	return `Scout runs a deterministic research-intelligence pipeline over Hugging Face Daily Papers.

It ingests the daily feed, retrieves paper metadata and markdown for transient scoring and summarization, persists compact curated records, and renders daily and monthly executive briefings.`
}
