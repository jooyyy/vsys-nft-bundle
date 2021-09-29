package main

import (
	"github.com/mitchellh/cli"
	"log"
	"os"
	"vsys-nft-bundle/server"
)

func main() {
	c := cli.NewCLI("v.systems nft", "1.0")
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