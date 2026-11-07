package app

import (
	"context"
	"errors"
	"io"
	"log"
	"strings"
	"testing"
	"time"

	"biosky/internal/bio"
)

type fakeGenerator struct {
	text string
	slot int64
}

func (g *fakeGenerator) generate(slot int64) string {
	g.slot = slot
	return g.text
}

type fakeUpdater struct {
	loginErr   error
	updateErr  error
	changed    bool
	got        string
	loginCalls int
	putCalls   int
}

func (u *fakeUpdater) Login(context.Context) error {
	u.loginCalls++
	return u.loginErr
}

func (u *fakeUpdater) UpdateDescription(_ context.Context, description string) (bool, error) {
	u.putCalls++
	u.got = description
	return u.changed, u.updateErr
}

func newApp(gen *fakeGenerator, updater *fakeUpdater, now time.Time) *App {
	return New(
		gen.generate,
		updater,
		func() time.Time { return now },
		log.New(io.Discard, "", 0),
	)
}

func TestRunGeneratesForSlotAndUpdates(t *testing.T) {
	now := time.Date(2026, 11, 7, 12, 0, 0, 0, time.UTC)
	gen := &fakeGenerator{text: "hello"}
	updater := &fakeUpdater{changed: true}
	if err := newApp(gen, updater, now).Run(context.Background()); err != nil {
		t.Fatal(err)
	}
	if gen.slot != bio.Slot(now) {
		t.Fatalf("slot = %d, want %d", gen.slot, bio.Slot(now))
	}
	if updater.got != "hello" || updater.loginCalls != 1 || updater.putCalls != 1 {
		t.Fatalf("unexpected updater state: %+v", updater)
	}
}

func TestRunUnchangedIsNotAnError(t *testing.T) {
	now := time.Date(2026, 11, 7, 12, 0, 0, 0, time.UTC)
	updater := &fakeUpdater{changed: false}
	if err := newApp(&fakeGenerator{text: "same"}, updater, now).Run(context.Background()); err != nil {
		t.Fatal(err)
	}
	if updater.putCalls != 1 {
		t.Fatalf("putCalls = %d", updater.putCalls)
	}
}

func TestRunPropagatesLoginError(t *testing.T) {
	now := time.Date(2026, 11, 7, 12, 0, 0, 0, time.UTC)
	updater := &fakeUpdater{loginErr: errors.New("boom")}
	err := newApp(&fakeGenerator{text: "hello"}, updater, now).Run(context.Background())
	if err == nil || !strings.Contains(err.Error(), "login") {
		t.Fatalf("error = %v", err)
	}
	if updater.putCalls != 0 {
		t.Fatal("update was attempted after a login failure")
	}
}

func TestRunPropagatesUpdateError(t *testing.T) {
	now := time.Date(2026, 11, 7, 12, 0, 0, 0, time.UTC)
	updater := &fakeUpdater{updateErr: errors.New("boom")}
	err := newApp(&fakeGenerator{text: "hello"}, updater, now).Run(context.Background())
	if err == nil || !strings.Contains(err.Error(), "update") {
		t.Fatalf("error = %v", err)
	}
}

func TestRunRejectsEmptyBioBeforeLogin(t *testing.T) {
	now := time.Date(2026, 11, 7, 12, 0, 0, 0, time.UTC)
	updater := &fakeUpdater{}
	err := newApp(&fakeGenerator{text: "  "}, updater, now).Run(context.Background())
	if err == nil {
		t.Fatal("expected an error")
	}
	if updater.loginCalls != 0 {
		t.Fatal("login was attempted with an invalid bio")
	}
}

func TestRunRejectsNilGenerator(t *testing.T) {
	now := time.Date(2026, 11, 7, 12, 0, 0, 0, time.UTC)
	updater := &fakeUpdater{}
	err := New(nil, updater, func() time.Time { return now }, log.New(io.Discard, "", 0)).Run(context.Background())
	if err == nil {
		t.Fatal("expected an error")
	}
	if updater.loginCalls != 0 {
		t.Fatal("login was attempted with a nil generator")
	}
}

func TestValidateDescription(t *testing.T) {
	if err := ValidateDescription(""); err == nil {
		t.Error("empty bio accepted")
	}
	if err := ValidateDescription(strings.Repeat("a", 256)); err != nil {
		t.Errorf("256 characters rejected: %v", err)
	}
	if err := ValidateDescription(strings.Repeat("a", 257)); err == nil {
		t.Error("257 characters accepted")
	}
	if err := ValidateDescription(strings.Repeat("😀", 256)); err != nil {
		t.Errorf("256 emoji rejected: %v", err)
	}
	if err := ValidateDescription(strings.Repeat("😀", 257)); err == nil {
		t.Error("257 emoji accepted")
	}
}
