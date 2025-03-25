package main

import (
	"TarantoolKV/internal/application/core/service"
	"TarantoolKV/internal/router/server"
	"TarantoolKV/internal/storage/tarantool"
	"context"
	"golang.org/x/sync/errgroup"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt, syscall.SIGTERM)

		<-c
		cancel()
	}()
	log.Println("initiating database...")
	db := tarantool.NewStorage()

	log.Println("creating app...")
	app := service.NewService(&db)

	log.Println("initializing server...")
	httpServer := server.SetupHTTPServer(app)

	g, gCtx := errgroup.WithContext(ctx)
	g.Go(func() error {
		return httpServer.ListenAndServe()
	})
	g.Go(func() error {
		<-gCtx.Done()
		log.Println("shutting down server...")
		return httpServer.Shutdown(gCtx)
	})
	g.Go(func() error {
		<-gCtx.Done()
		log.Println("closing database...")
		return db.Shutdown()
	})

	if err := g.Wait(); err != nil {
		log.Printf("exit reason: %s \n", err)
	}

	log.Println("app shut down")
}
