package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"biosky/internal/app"
	"biosky/internal/bio"
	"biosky/internal/bluesky"
	"biosky/internal/config"
	_ "biosky/internal/themes"
)

var version = "dev"

func main() {
	dryRun := flag.Bool("dry-run", false, "print the generated bio and exit")
	showVersion := flag.Bool("version", false, "print the version and exit")
	flag.Parse()

	if *showVersion {
		fmt.Println(version)
		return
	}

	selection, err := bio.Select()
	if err != nil {
		log.Fatalf("select theme: %v", err)
	}

	if *dryRun {
		slot := bio.Slot(time.Now())
		fmt.Printf("theme: %s\nslot: %d\nbio: %s\n", selection.Name, slot, selection.Generator.Generate(slot))
		return
	}

	cfg, err := config.Load(os.Getenv)
	if err != nil {
		log.Fatalf("config: %v", err)
	}

	client := bluesky.New(bluesky.Config{
		Handle:      cfg.Handle,
		AppPassword: cfg.AppPassword,
		PDS:         cfg.PDS,
		UserAgent:   "biosky/" + version,
	})
	application := app.New(selection, client, time.Now, log.Default())
	if err := application.Run(context.Background()); err != nil {
		log.Fatalf("%v", err)
	}
}
