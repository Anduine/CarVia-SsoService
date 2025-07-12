package server

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func StartServer(log *slog.Logger, router http.Handler, port string, timeout time.Duration) {
	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", port),
		Handler: router,
	}

	// Канал для остановки сервера
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	// Запуск сервера в отдельной горутине
	go func() {
		log.Info(fmt.Sprintf("SSO service running on port %s", port))
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Error("Server error: ", slog.Any("Error", err))
		}
	}()

	// Ожидание сигнала завершения
	<-stop
	log.Info("Shutting down server...")

	// Завершение работы сервера
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	log.Info("Shutdown ", slog.Any("stopcode", server.Shutdown(ctx)))
}