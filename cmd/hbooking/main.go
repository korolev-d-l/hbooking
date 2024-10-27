package main

import (
	"log/slog"

	"github.com/spatecon/hbooking/internal/pkg/app"
)

func main() {
	ap, err := app.NewApp()
	if err != nil {
		slog.Error("failed to create app", "error", err)
		return
	}

	if err = ap.Run(); err != nil {
		slog.Error("failed to run app", "error", err)
	}
}
