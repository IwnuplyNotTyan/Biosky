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

	if *dryRun {
		slot := bio.Slot(time.Now())
		fmt.Printf("theme: %s\nslot: %d\nbio: %s\n", bio.Theme(), slot, bio.Generate(slot))
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
	application := app.New(bio.Generate, client, time.Now, log.Default())
	if err := application.Run(context.Background()); err != nil {
		log.Fatalf("%v", err)
	}
}
