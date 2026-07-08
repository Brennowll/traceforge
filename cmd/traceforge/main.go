package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"path/filepath"
	"sort"
	"syscall"

	"github.com/brenno/traceforge/internal/output"
	"github.com/brenno/traceforge/internal/scenario"
	"github.com/brenno/traceforge/internal/simulation"
	"github.com/spf13/cobra"
)

var version = "dev"

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelWarn}))
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	cmd := newRootCommand(logger)
	cmd.SetContext(ctx)
	if err := cmd.Execute(); err != nil {
		logger.Error("command failed", "error", err)
		os.Exit(1)
	}
}

func newRootCommand(logger *slog.Logger) *cobra.Command {
	root := &cobra.Command{
		Use:   "traceforge",
		Short: "Distributed system simulator written in Go",
	}

	root.AddCommand(&cobra.Command{
		Use:   "version",
		Short: "Print TraceForge version",
		Run: func(cmd *cobra.Command, args []string) {
			_, _ = fmt.Fprintf(cmd.OutOrStdout(), "traceforge %s\n", version)
		},
	})
	root.AddCommand(newRunCommand(logger))
	return root
}

func newRunCommand(logger *slog.Logger) *cobra.Command {
	var entry string
	var requests int
	var concurrency int
	var seed int64
	var htmlPath string

	cmd := &cobra.Command{
		Use:   "run scenario.yml",
		Short: "Run a distributed system simulation",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			path := args[0]
			s, err := scenario.LoadFile(path)
			if err != nil {
				return err
			}
			if entry == "" {
				entry = defaultEntry(s)
			}
			if err := scenario.ValidateEntry(s, entry); err != nil {
				return err
			}

			seeded := cmd.Flags().Changed("seed")
			logger.Debug("running simulation", "scenario", path, "entry", entry, "requests", requests, "concurrency", concurrency, "seeded", seeded)
			batch, err := simulation.New(s).RunBatch(cmd.Context(), entry, simulation.BatchOptions{
				Requests:    requests,
				Concurrency: concurrency,
				Seed:        seed,
				Seeded:      seeded,
			})
			if err != nil {
				return err
			}

			if err := (output.TextRenderer{}).Render(cmd.OutOrStdout(), batch); err != nil {
				return err
			}
			if htmlPath != "" {
				data := output.HTMLReportData{ScenarioName: filepath.Base(path), Entry: entry, Batch: batch}
				if err := (output.HTMLRenderer{}).WriteFile(htmlPath, data); err != nil {
					return err
				}
				_, _ = fmt.Fprintf(cmd.OutOrStdout(), "\nHTML report: %s\n", htmlPath)
			}
			return nil
		},
	}

	cmd.Flags().StringVar(&entry, "entry", "", "entry service name (defaults to api when present, otherwise first service alphabetically)")
	cmd.Flags().IntVar(&requests, "requests", 1, "number of requests to simulate")
	cmd.Flags().IntVar(&concurrency, "concurrency", 1, "maximum concurrent requests")
	cmd.Flags().Int64Var(&seed, "seed", 0, "deterministic seed")
	cmd.Flags().StringVar(&htmlPath, "html", "", "write static HTML report to this path")
	return cmd
}

func defaultEntry(s *scenario.Scenario) string {
	if _, ok := s.Services["api"]; ok {
		return "api"
	}
	names := make([]string, 0, len(s.Services))
	for name := range s.Services {
		names = append(names, name)
	}
	sort.Strings(names)
	return names[0]
}
