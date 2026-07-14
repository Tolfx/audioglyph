package cmds

import (
	"errors"
	"fmt"
	"strings"

	"github.com/Tolfx/audioglyph/internal/api/discogs"
	"github.com/Tolfx/audioglyph/internal/db/models"
	"github.com/charmbracelet/huh"
	"github.com/spf13/cobra"
	"gorm.io/gorm"
)

var CmdAdd = &cobra.Command{
	Use:   "add <barcode>",
	Short: "Add a record to the collection by barcode",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		barcode := args[0]

		var existing models.PhysicalMedia
		if err := Service.Database.DB.Where("barcode = ?", barcode).First(&existing).Error; err == nil {
			return fmt.Errorf("barcode %s is already in your collection (physical media ID %d)", barcode, existing.ID)
		} else if !errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("checking for existing barcode: %w", err)
		}

		Service.Logger.Info("searching for barcode", "barcode", barcode)

		results, err := Service.Discogs.SearchByBarcode(barcode)
		if err != nil {
			return fmt.Errorf("searching discogs: %w", err)
		}

		if len(results.Results) == 0 {
			return fmt.Errorf("no releases found for barcode %s", barcode)
		}

		var selected *discogs.Release

		if len(results.Results) == 1 {
			selected = &results.Results[0]
			Service.Logger.Info("single result found, auto-selecting", "title", selected.Title)
		} else {
			selected, err = selectRelease(results.Results)
			if err != nil {
				return err
			}
		}

		Service.Logger.Info("selected release", "title", selected.Title, "year", selected.Year, "id", selected.ID)

		if err := Service.AddToCollection(selected.ID, barcode); err != nil {
			return fmt.Errorf("adding to collection: %w", err)
		}

		fmt.Printf("✓ Added \"%s\" to your collection.\n", selected.Title)
		return nil
	},
}

// releaseLabel builds a human-readable label for a release shown in the selector.
func releaseLabel(r discogs.Release) string {
	parts := []string{r.Title}
	if r.Year != "" {
		parts = append(parts, r.Year)
	}
	if len(r.Label) > 0 {
		parts = append(parts, r.Label[0])
	}
	if len(r.Format) > 0 {
		parts = append(parts, strings.Join(r.Format, "/"))
	}
	if r.Country != "" {
		parts = append(parts, r.Country)
	}
	return strings.Join(parts, " · ")
}

// selectRelease presents an interactive list of releases and returns the one the user picks.
func selectRelease(releases []discogs.Release) (*discogs.Release, error) {
	options := make([]huh.Option[int], len(releases))
	for i, r := range releases {
		options[i] = huh.NewOption(releaseLabel(r), i)
	}

	var idx int
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[int]().
				Title("Multiple releases found — select the correct one").
				Options(options...).
				Value(&idx),
		),
	)

	if err := form.Run(); err != nil {
		return nil, fmt.Errorf("selection cancelled: %w", err)
	}

	return &releases[idx], nil
}
