package main

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/krypton-0x00/go-test-rest-api/internal/config"
)

func main() {
	// LOAD CONFIG
	cfg := config.MustLoad()

	// DB SETUP

	// SETUP ROUTER
	router := http.NewServeMux()

	router.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Testing"))
	})

	// SETUP SERVER
	server := http.Server{
		Addr:    cfg.Addr,
		Handler: router,
	}

	slog.Info("[+] Server Started At ", slog.String("Address", cfg.Addr))

	done := make(chan os.Signal, 1)

	// ON INTERUPT SEND MESSAGE TO THE DONE CHANNEL
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	// SERVER OF SEPERATE GOROUTINE
	go func() {
		err := server.ListenAndServe()
		if err != nil {
			log.Fatal("[x] Failed to start server.")
		}
	}()

	// BLOCK UNTIL SHUTDOWN SIG IS RECIEVED
	<-done

	slog.Info("[+] Shutting Down the server.")
	// GRACEFULL SHUTDOWN
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := server.Shutdown(ctx)
	if err != nil {
		slog.Error("[x] Failed to shutdown server", slog.String("error ", err.Error()))
	}

	slog.Info("[+] Server Shutdown Sucessfully ")
}
