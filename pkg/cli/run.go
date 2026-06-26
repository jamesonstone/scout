package cli

import (
	"context"
	"fmt"
	"time"

	"github.com/jamesonstone/scout/internal/config"
	"github.com/jamesonstone/scout/internal/hf"
	"github.com/jamesonstone/scout/internal/pipeline"
	"github.com/jamesonstone/scout/internal/storage"
	"github.com/spf13/cobra"
)

func init() {
	cmd := &cobra.Command{
		Use:   "run",
		Short: "Run the Scout daily pipeline",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runPipeline(cmd.Context(), cmd)
		},
	}
	cmd.Flags().String("date", "", "override run date (YYYY-MM-DD)")
	cmd.Flags().String("data-dir", "", "override scout data directory")
	cmd.Flags().String("base-url", "", "override Hugging Face base URL")
	cmd.Flags().Duration("timeout", 30*time.Second, "HTTP timeout")
	cmd.Flags().Int("retries", 3, "HTTP retry count")
	rootCmd.AddCommand(cmd)
}

func runPipeline(ctx context.Context, cmd *cobra.Command) error {
	cfg, err := config.FromEnv()
	if err != nil {
		return newCLIExitError(err, 1, false)
	}
	if cmd.Flags().Changed("date") {
		cfg.RunDate = cmd.Flag("date").Value.String()
	}
	if cmd.Flags().Changed("data-dir") {
		cfg.DataDir = cmd.Flag("data-dir").Value.String()
	}
	if cmd.Flags().Changed("base-url") {
		cfg.BaseURL = cmd.Flag("base-url").Value.String()
	}
	if cmd.Flags().Changed("timeout") {
		timeout, err := cmd.Flags().GetDuration("timeout")
		if err != nil {
			return newCLIExitError(err, 1, false)
		}
		cfg.Timeout = timeout
	}
	if cmd.Flags().Changed("retries") {
		retries, err := cmd.Flags().GetInt("retries")
		if err != nil {
			return newCLIExitError(err, 1, false)
		}
		cfg.Retries = retries
	}

	store := storage.New(cfg.DataDir)
	client := hf.NewClient(cfg)
	runner := pipeline.NewRunner(cfg, client, store)
	result, err := runner.Run(ctx)
	if err != nil {
		return newCLIExitError(err, 1, false)
	}

	fmt.Printf("Daily report: %s\n", result.DailyReportPath)
	fmt.Printf("Monthly report: %s\n", result.MonthlyReportPath)
	fmt.Printf("Processed %d paper(s); reused %d existing record(s).\n", result.ProcessedCount, result.ReusedCount)
	return nil
}
