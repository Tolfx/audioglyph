package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/Tolfx/audioglyph/cmd/cli/cmds"
	"github.com/Tolfx/audioglyph/internal/db"
	"github.com/Tolfx/audioglyph/internal/service"
	"github.com/Tolfx/audioglyph/internal/utils"
	"github.com/spf13/cobra"
)

var Service *service.Service

var cmdRoot = &cobra.Command{
	Use:   "audioglyph",
	Short: "A tool for managing your physical music collection",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: utils.If(os.Getenv("DEBUG") == "true", slog.LevelDebug, slog.LevelInfo)}))
		ctx := context.Background()

		database, err := db.NewDatabase(os.Getenv("DATABASE_URL"))
		if err != nil {
			return err
		}

		Service = service.NewService(&ctx, logger)
		Service.SetDatabase(database)
		cmds.Service = Service
		return nil
	},
}

func init() {
	cmdRoot.AddCommand(cmds.CmdAdd)
	cmdRoot.AddCommand(cmds.CmdView)
}

func main() {
	if err := cmdRoot.Execute(); err != nil {
		os.Exit(1)
	}
}
