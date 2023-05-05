package main

import (
	"fmt"
	"log"
	"os"

	"github.com/chopper-c2-framework/c2-chopper/core"
	"github.com/chopper-c2-framework/c2-chopper/core/config"
	"github.com/chopper-c2-framework/c2-chopper/core/plugins"

	// Fix architecture and make it one import !
	server "github.com/chopper-c2-framework/c2-chopper/server"
	serverGrpc "github.com/chopper-c2-framework/c2-chopper/server/grpc"
	// "github.com/chopper-c2-framework/c2-chopper/server/grpc"
)

func setupCli() {

}
func x() error {
	return nil
}

func main() {

	configCommands := config.GetCommands()

	serverCommands := server.GetCommands(serverGrpc.ServerManager{})

	framework := core.CreateApp(configCommands, serverCommands)
	// frameworkConfiguration := config.ParseConfigFromPath()

	plugins, err := plugins.LoadPlugins()

	if err != nil {
		log.Fatalf("%s", err)
	}

	for _, plugin := range plugins {
		fmt.Println("[+]", plugin.Name)
	}

	// go grpc.NewgRPCServer(*frameworkConfiguration)

	err = framework.Run(os.Args)
	if err != nil {
		log.Panicln("Error occured while launching the framework", err)
	}
}
