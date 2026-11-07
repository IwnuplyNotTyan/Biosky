package app

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"
	"unicode/utf8"

	"biosky/internal/bio"
)

// MaxDescriptionLength is the Bluesky limit for a profile description.
const MaxDescriptionLength = 256

// Updater is the slice of the Bluesky client the app needs.
type Updater interface {
	Login(ctx context.Context) error
	UpdateDescription(ctx context.Context, description string) (bool, error)
}

// App runs one bio update.
type App struct {
	selection bio.Selection
	updater   Updater
	now       func() time.Time
	logger    *log.Logger
}

// New builds an App; a nil logger falls back to log.Default().
func New(selection bio.Selection, updater Updater, now func() time.Time, logger *log.Logger) *App {
	if now == nil {
		now = time.Now
	}
	if logger == nil {
		logger = log.Default()
	}
	return &App{selection: selection, updater: updater, now: now, logger: logger}
}

// Run generates the bio for the current slot and writes it if needed.
func (a *App) Run(ctx context.Context) error {
	if a.updater == nil {
		return errors.New("updater is nil")
	}
	slot := bio.Slot(a.now())
	text := a.selection.Generator.Generate(slot)
	if err := ValidateDescription(text); err != nil {
		return err
	}
	if err := a.updater.Login(ctx); err != nil {
		return fmt.Errorf("login: %w", err)
	}
	changed, err := a.updater.UpdateDescription(ctx, text)
	if err != nil {
		return fmt.Errorf("update description: %w", err)
	}
	if changed {
		a.logger.Printf("theme=%s slot=%d updated: %s", a.selection.Name, slot, text)
	} else {
		a.logger.Printf("theme=%s slot=%d unchanged", a.selection.Name, slot)
	}
	return nil
}

// ValidateDescription enforces the non-empty and length rules.
func ValidateDescription(s string) error {
	if strings.TrimSpace(s) == "" {
		return errors.New("generated bio is empty")
	}
	if n := utf8.RuneCountInString(s); n > MaxDescriptionLength {
		return fmt.Errorf("generated bio is %d characters; max %d", n, MaxDescriptionLength)
	}
	return nil
}
