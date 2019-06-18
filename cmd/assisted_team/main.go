package main

import (
	"context"
	"fmt"
	"github.com/misterfaradey/assisted_team_test/internal/config"
	"github.com/misterfaradey/assisted_team_test/internal/controller"
	"github.com/misterfaradey/assisted_team_test/internal/server"
	"github.com/misterfaradey/assisted_team_test/internal/service"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {

	// Context
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Config
	cfg, err := config.Init()
	if err != nil {
		panic(err)
	}

	service := service.NewAssistedService()
	err = service.Update(cfg.MainConf.DataFilesReturn, cfg.MainConf.DataFilesOneWay)
	if err != nil {
		panic(err)
	}
	assistedController := controller.NewController(service)
	assistedServer := server.NewServer(
		assistedController,
		cfg.UserServerConf,
	)

	fmt.Println("listen", cfg.UserServerConf.Address)
	go runServer(cancel, assistedServer)

	// graceful shutdown
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGKILL)
	select {
	case <-signalChan:
	case <-ctx.Done():
	}

	ctx, cancel = context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	err = assistedServer.Shutdown(ctx)
	if err != nil {
		log.Println(err)

	}
}

func runServer(cancel context.CancelFunc, server server.Server) {

	err := server.Run()
	if err == http.ErrServerClosed {
		log.Println("Closed server")
	} else if err != nil {

		log.Println(err)
		cancel()
	}
}
