package main

import (
	"fmt"
	"log"
	"os"

	"github.com/MikeRoetgers/deploi/cmd/client/cmd"
	"github.com/MikeRoetgers/deploi/cmd/client/config"
	"github.com/MikeRoetgers/deploi/protobuf"
	"google.golang.org/grpc"
)

var (
	deploiConfiguration *config.Configuration
)

func main() {
	deploiConfiguration = getConfig()
	grpcConn, err := grpc.Dial(deploiConfiguration.Host, grpc.WithInsecure())
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer grpcConn.Close()
	cmd.DeploiClient = protobuf.NewDeploiServerClient(grpcConn)
	cmd.Execute()
}

func getConfig() *config.Configuration {
	confPath := os.Getenv("DEPLOI_CONFIG")
	if confPath != "" {
		conf, err := config.NewConfigFromFile(confPath)
		if err != nil {
			log.Fatal(err)
		}
		return conf
	}

	if _, err := os.Stat(config.GetDefaultConfLocation()); !os.IsNotExist(err) {
		conf, err := config.NewConfigFromFile(config.GetDefaultConfLocation())
		if err != nil {
			log.Fatal(err)
		}
		return conf
	}
	return config.NewConfig()
}
