package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	api "flip/internal/api/handler"
	"flip/internal/model"
	bus "flip/internal/repository/bus/service"
	csvparser "flip/internal/repository/csvparser/service"
	storage "flip/internal/repository/storage/service"
	worker "flip/internal/repository/worker/service"
	"flip/util"
)

func main() {
	util.Log = slog.New(slog.NewJSONHandler(os.Stdout, nil))

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	storageSvc := storage.NewStorage()
	eventBus := bus.NewDefault()

	handler := func(ctx context.Context, ev model.Event) *worker.HandleResult {
		_ = storageSvc.AddSuccess(ctx, ev.UploadID, ev.Tx)
		return &worker.HandleResult{Success: true}
	}

	w := worker.New(eventBus.Channel(), handler)
	w.Start(ctx, 10)
	defer w.Stop()
	parser := csvparser.NewParser(storageSvc, eventBus, w)

	h := api.New(storageSvc, parser, eventBus)
	mux := http.NewServeMux()
	h.Register(mux)

	srv := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	util.Log.Info("server started", "addr", srv.Addr)

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			util.Log.Error("server error", "err", err)
		}
	}()

	<-ctx.Done()
	util.Log.Info("shutdown initiated")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		util.Log.Error("shutdown failed", "err", err)
	}

	util.Log.Info("shutdown finished")

	util.Perf(context.Background(), "shutdown")
}
