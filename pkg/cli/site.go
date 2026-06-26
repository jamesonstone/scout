package cli

import (
	"fmt"

	"github.com/jamesonstone/scout/internal/site"
	"github.com/spf13/cobra"
)

func init() {
	siteCmd := &cobra.Command{
		Use:   "site",
		Short: "Build and validate the static Scout site",
	}

	buildCmd := &cobra.Command{
		Use:   "build",
		Short: "Build the static GitHub Pages site",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg := siteConfigFromFlags(cmd)
			result, err := site.Build(cfg)
			if err != nil {
				return newCLIExitError(err, 1, false)
			}
			fmt.Printf("Static site: %s\n", result.OutDir)
			fmt.Printf("Generated %d daily page(s), %d monthly page(s), and %d paper page(s).\n", result.DailyPages, result.MonthlyPages, result.PaperPages)
			return nil
		},
	}
	addSiteFlags(buildCmd)

	validateCmd := &cobra.Command{
		Use:   "validate",
		Short: "Validate generated static site output",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg := siteConfigFromFlags(cmd)
			result, err := site.Validate(cfg)
			if err != nil {
				return newCLIExitError(err, 1, false)
			}
			fmt.Printf("Validated %d page(s) and %d link(s).\n", result.CheckedPages, result.CheckedLinks)
			return nil
		},
	}
	addSiteFlags(validateCmd)

	siteCmd.AddCommand(buildCmd, validateCmd)
	rootCmd.AddCommand(siteCmd)
}

func addSiteFlags(cmd *cobra.Command) {
	cmd.Flags().String("data-dir", ".", "scout data directory containing data/ and reports/")
	cmd.Flags().String("out-dir", "public", "static site output directory")
	cmd.Flags().String("base-path", "/scout/", "GitHub Pages project-site base path")
}

func siteConfigFromFlags(cmd *cobra.Command) site.Config {
	return site.Config{
		DataDir:  cmd.Flag("data-dir").Value.String(),
		OutDir:   cmd.Flag("out-dir").Value.String(),
		BasePath: cmd.Flag("base-path").Value.String(),
	}
}
