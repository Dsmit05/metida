package utils

import (
	"context"
	"sync"
	"time"
)

type App interface {
	Stop(ctx context.Context)
}

func Shutdown(ctx context.Context, apps ...App) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	wg := sync.WaitGroup{}

	for _, app := range apps {
		wg.Add(1)
		app := app

		go func() {
			defer wg.Done()
			app.Stop(ctx)
		}()
	}

	wg.Wait()
}
