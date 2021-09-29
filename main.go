package main

import (
	"github.com/mitchellh/cli"
	"log"
	"os"
	"vsys-nft-bundle/server"
)

const VERSION = "0.1"

func main() {
	c := cli.NewCLI("v.systems nft", VERSION)
	c.Args = os.Args[1:]
	c.Commands = map[string]cli.CommandFactory{
		"api": server.CommandFactory,
	}

	exitStatus, err := c.Run()
	if err != nil {
		log.Println(err)
	}

	os.Exit(exitStatus)
}