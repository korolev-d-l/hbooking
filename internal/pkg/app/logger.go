package app

import "log/slog"

func (a *App) initLogger() {
	logger := slog.Default()
	a.logger = logger
}
