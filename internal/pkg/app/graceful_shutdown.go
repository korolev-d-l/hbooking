package app

import (
	"os"
	"os/signal"
)

func (a *App) initGracefulShutdown() error {
	a.closeCh = make(chan os.Signal, 1)
	signal.Notify(a.closeCh, os.Interrupt)

	return nil
}

func (a *App) waitCloseSignal() {
	<-a.closeCh
}

func (a *App) stopWaitCloseSignal() {
	signal.Stop(a.closeCh)
	select {
	case <-a.closeCh:
	default:
		close(a.closeCh)
	}
}

func (a *App) shutdown() {
	for i := len(a.closers) - 1; i >= 0; i-- {
		err := a.closers[i]()
		if err != nil {
			a.logger.Error("failed to close resource", "i", i, "error", err)
		}
	}
}
