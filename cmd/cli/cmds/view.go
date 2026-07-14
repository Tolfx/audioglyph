package cmds

import (
	"fmt"
	"strings"

	"github.com/Tolfx/audioglyph/internal/db/models"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
)

// ── Styles ────────────────────────────────────────────────────────────────────

var (
	styleBold    = lipgloss.NewStyle().Bold(true)
	styleMuted   = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	styleAccent  = lipgloss.NewStyle().Foreground(lipgloss.Color("205")).Bold(true)
	styleHeader  = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("99"))
	styleDivider = lipgloss.NewStyle().Foreground(lipgloss.Color("237"))

	styleBox = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("99")).
			Padding(1, 2).
			MarginTop(1)

	styleTrackPos    = lipgloss.NewStyle().Foreground(lipgloss.Color("240")).Width(6)
	styleTrackTitle  = lipgloss.NewStyle().Width(38)
	styleTrackArtist = lipgloss.NewStyle().Foreground(lipgloss.Color("240")).Width(22)
	styleTrackDur    = lipgloss.NewStyle().Foreground(lipgloss.Color("240")).Align(lipgloss.Right)
)

// ── Command ───────────────────────────────────────────────────────────────────

var CmdView = &cobra.Command{
	Use:   "view",
	Short: "Browse your physical music collection",
	RunE: func(cmd *cobra.Command, args []string) error {
		var items []models.PhysicalMedia
		if err := Service.Database.DB.
			Preload("Release").
			Preload("Release.Artist").
			Preload("Tracks").
			Where("in_collection = ?", true).
			Order("created_at desc").
			Find(&items).Error; err != nil {
			return fmt.Errorf("loading collection: %w", err)
		}

		if len(items) == 0 {
			fmt.Println(styleMuted.Render("Your collection is empty. Use `audioglyph add <barcode>` to get started."))
			return nil
		}

		selected, err := pickItem(items)
		if err != nil {
			return err
		}

		fmt.Println(renderDetail(selected))
		return nil
	},
}

// pickItem shows a searchable list of all items and returns the chosen one.
func pickItem(items []models.PhysicalMedia) (*models.PhysicalMedia, error) {
	opts := make([]huh.Option[int], len(items))
	for i, it := range items {
		label := fmt.Sprintf("%-32s  %-28s  %s  %s",
			truncate(it.Release.Artist.Name, 32),
			truncate(it.Release.Title, 28),
			styleMuted.Render(fmt.Sprintf("%d", it.Release.Year)),
			styleMuted.Render(fmt.Sprintf("[%s]", it.Format)),
		)
		opts[i] = huh.NewOption(label, i)
	}

	var idx int
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[int]().
				Title("Your Collection — select an item to view details").
				Options(opts...).
				Value(&idx),
		),
	)
	if err := form.Run(); err != nil {
		return nil, fmt.Errorf("selection cancelled: %w", err)
	}
	return &items[idx], nil
}

// renderDetail builds the full lipgloss detail card for a physical media item.
func renderDetail(pm *models.PhysicalMedia) string {
	var b strings.Builder

	// ── Header ────────────────────────────────────────────────────────────────
	b.WriteString(styleAccent.Render(pm.Release.Artist.Name))
	b.WriteString("\n")
	b.WriteString(styleBold.Render(pm.Release.Title))
	b.WriteString("  ")
	b.WriteString(styleMuted.Render(fmt.Sprintf("(%d)", pm.Release.Year)))
	b.WriteString("\n\n")

	// ── Metadata grid ─────────────────────────────────────────────────────────
	rows := []struct{ label, value string }{
		{"Format", string(pm.Format)},
		{"Issued", fmt.Sprintf("%d", pm.IssuedYear)},
		{"Barcode", pm.Barcode},
	}
	if pm.CatalogNumber != "" {
		rows = append(rows, struct{ label, value string }{"Cat. No.", pm.CatalogNumber})
	}
	for _, r := range rows {
		b.WriteString(fmt.Sprintf("  %s  %s\n",
			styleHeader.Render(fmt.Sprintf("%-10s", r.label)),
			r.value,
		))
	}

	// ── Tracklist ─────────────────────────────────────────────────────────────
	if len(pm.Tracks) > 0 {
		b.WriteString("\n")
		b.WriteString(styleDivider.Render(strings.Repeat("─", 56)))
		b.WriteString("\n")
		b.WriteString(styleHeader.Render("  Tracklist"))
		b.WriteString("\n")
		b.WriteString(styleDivider.Render(strings.Repeat("─", 56)))
		b.WriteString("\n")

		for _, t := range pm.Tracks {
			b.WriteString(styleTrackPos.Render(t.Position))
			b.WriteString(styleTrackTitle.Render(truncate(t.Title, 38)))
			b.WriteString(styleTrackArtist.Render(truncate(t.Artist, 22)))
			if t.DurationSeconds > 0 {
				b.WriteString(styleTrackDur.Render(formatDuration(t.DurationSeconds)))
			}
			b.WriteString("\n")
		}
	}

	return styleBox.Render(b.String())
}

// ── Helpers ───────────────────────────────────────────────────────────────────

func truncate(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return s[:max-1] + "…"
}

func formatDuration(secs int) string {
	return fmt.Sprintf("%d:%02d", secs/60, secs%60)
}
