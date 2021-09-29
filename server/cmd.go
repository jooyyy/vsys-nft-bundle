package server

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/mitchellh/cli"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"vsys-nft-bundle/config"
)

type Command struct{}

func CommandFactory() (cli.Command, error) {
	return new(Command), nil
}

type Service struct {
	mode       string
	httpServer *http.Server
}

func NewService(args ...string) *Service {
	serviceMode := "debug"
	if len(args) > 0 && args[0] == "release" {
		serviceMode = "release"
	}

	gin.SetMode(serviceMode)

	handler := gin.Default()

	addr := fmt.Sprintf(":%d", config.Config.Port)
	server := &http.Server{
		Addr:    addr,
		Handler: handler,
	}

	service := &Service{
		mode:       serviceMode,
		httpServer: server,
	}

	service.initRouter(handler)

	return service
}

func (s *Service) Run() int {

	startRestService := func() {
		fmt.Println("Start rest api server", s.httpServer.Addr)

		if err := s.httpServer.ListenAndServe(); err != nil {
			fmt.Errorf("RestServer.Run %s", err)
		}
		fmt.Println("rest server shutdown.")
	}

	waitToStop := func() {
		signalChan := make(chan os.Signal, 1)
		signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
		<-signalChan
	}

	stopAllService := func() {
		s.httpServer.Shutdown(context.Background())
		fmt.Println("all service has stopped, exit.")
	}

	go startRestService()

	waitToStop()
	stopAllService()

	return 1
}

func (c *Command) Run(args []string) int {
	return NewService(args...).Run()
}

func (c *Command) Help() string {
	return help
}

func (c *Command) Synopsis() string {
	return synopsis
}

const synopsis = "v.systems NFT bundle"

const help = ``