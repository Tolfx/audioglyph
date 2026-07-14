package service

import (
	"context"
	"log/slog"
	"os"

	"github.com/Tolfx/audioglyph/internal/api/discogs"
	"github.com/Tolfx/audioglyph/internal/db"
)

type Service struct {
	Context  *context.Context
	Logger   *slog.Logger
	Database *db.Database
	Discogs  *discogs.Client
}

func NewService(context *context.Context, logger *slog.Logger) *Service {
	return &Service{
		Context: context,
		Logger:  logger,
		Discogs: discogs.NewClient(logger, os.Getenv("DISCOGS_TOKEN")),
	}
}

func (service *Service) SetDatabase(database *db.Database) {
	service.Database = database
}
