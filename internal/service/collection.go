package service

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/Tolfx/audioglyph/internal/api/discogs"
	"github.com/Tolfx/audioglyph/internal/db/models"
	"gorm.io/gorm"
)

// AddToCollection fetches full release details from Discogs, maps them to the
// database models and persists everything in a single transaction.
func (s *Service) AddToCollection(releaseID int, barcode string) error {

	detail, err := s.Discogs.GetRelease(releaseID)
	if err != nil {
		return fmt.Errorf("fetching release details: %w", err)
	}

	s.Logger.Info("fetched release details", "title", detail.Title, "year", detail.Year)

	return s.Database.DB.Transaction(func(tx *gorm.DB) error {
		// --- Artist ---
		artistName := primaryArtistName(detail)
		var artist models.Artist
		if err := tx.Where(models.Artist{Name: artistName}).FirstOrCreate(&artist).Error; err != nil {
			return fmt.Errorf("upserting artist: %w", err)
		}

		// --- Release ---
		var release models.Release
		if err := tx.Where(models.Release{ArtistID: artist.ID, Title: detail.Title}).
			FirstOrCreate(&release, models.Release{
				ArtistID: artist.ID,
				Title:    detail.Title,
				Year:     detail.Year,
			}).Error; err != nil {
			return fmt.Errorf("upserting release: %w", err)
		}

		// --- PhysicalMedia ---
		catalogNumber := primaryCatalogNumber(detail)
		physical := models.PhysicalMedia{
			ReleaseID:     release.ID,
			Format:        mapFormat(detail.Formats),
			Barcode:       barcode,
			CatalogNumber: catalogNumber,
			IssuedYear:    detail.Year,
			InCollection:  true,
		}
		if err := tx.Create(&physical).Error; err != nil {
			return fmt.Errorf("creating physical media: %w", err)
		}

		// --- Tracks ---
		for _, t := range detail.Tracklist {
			// Skip side/disc headings
			if strings.EqualFold(t.Type, "heading") {
				continue
			}
			// Use track-level artist when present, otherwise fall back to the release artist.
			trackArtist := artist.Name
			if len(t.Artists) > 0 {
				trackArtist = trackArtistName(t.Artists)
			}
			track := models.Track{
				PhysicalMediaID: physical.ID,
				Artist:          trackArtist,
				Position:        t.Position,
				Title:           t.Title,
				DurationSeconds: parseDuration(t.Duration),
			}
			if err := tx.Create(&track).Error; err != nil {
				return fmt.Errorf("creating track %q: %w", t.Title, err)
			}
		}

		s.Logger.Info("saved to collection",
			"artist", artist.Name,
			"release", release.Title,
			"format", physical.Format,
			"tracks", len(detail.Tracklist),
		)
		return nil
	})
}

// primaryArtistName returns the first artist name from a release detail.
func primaryArtistName(detail *discogs.ReleaseDetail) string {
	if len(detail.Artists) > 0 {
		name := detail.Artists[0].Name
		// Discogs appends " (N)" suffixes to disambiguate duplicate artist names.
		if idx := strings.LastIndex(name, " ("); idx != -1 {
			name = name[:idx]
		}
		return strings.TrimSpace(name)
	}
	return "Unknown Artist"
}

// trackArtistName joins multiple track-level artists with " & ".
func trackArtistName(artists []discogs.ReleaseArtist) string {
	names := make([]string, 0, len(artists))
	for _, a := range artists {
		name := a.Name
		if idx := strings.LastIndex(name, " ("); idx != -1 {
			name = name[:idx]
		}
		names = append(names, strings.TrimSpace(name))
	}
	return strings.Join(names, " & ")
}

// primaryCatalogNumber returns the catalog number from the first label entry.
func primaryCatalogNumber(detail *discogs.ReleaseDetail) string {
	if len(detail.Labels) > 0 {
		return detail.Labels[0].CatNo
	}
	return ""
}

// mapFormat converts Discogs format descriptions to the internal PhysicalType.
func mapFormat(formats []discogs.ReleaseFormat) models.PhysicalType {
	for _, f := range formats {
		switch strings.ToLower(f.Name) {
		case "vinyl":
			return models.PhysicalTypeVinyl
		case "cd", "cdr":
			return models.PhysicalTypeCD
		case "cassette":
			return models.PhysicalTypeTape
		}
		for _, desc := range f.Descriptions {
			switch strings.ToLower(desc) {
			case "lp", "ep", "single":
				return models.PhysicalTypeVinyl
			}
		}
	}
	return models.PhysicalTypeVinyl // sensible default
}

// parseDuration converts a "m:ss" string from Discogs to total seconds.
func parseDuration(s string) int {
	s = strings.TrimSpace(s)
	if s == "" {
		return 0
	}
	parts := strings.Split(s, ":")
	if len(parts) != 2 {
		return 0
	}
	minutes, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0
	}
	seconds, err := strconv.Atoi(parts[1])
	if err != nil {
		return 0
	}
	return minutes*60 + seconds
}
