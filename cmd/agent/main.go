package main

import (
	"fmt"
	"os"
	"time"

	"github.com/MikeRoetgers/deploi/protobuf"
	logging "github.com/op/go-logging"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

var log = logging.MustGetLogger("app")

func main() {
	setupConfig()
	options := []grpc.DialOption{}
	if viper.GetBool("TLS.useTLS") {
		creds, err := credentials.NewClientTLSFromFile(viper.GetString("TLS.certFile"), "")
		if err != nil {
			fmt.Printf("Failed to read TLS cert file: %s", err)
			os.Exit(1)
		}
		options = append(options, grpc.WithTransportCredentials(creds))
	} else {
		options = append(options, grpc.WithInsecure())
	}
	grpcConn, err := grpc.Dial(viper.GetString("deploidHost"), options...)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer grpcConn.Close()
	deploiClient := protobuf.NewDeploiServerClient(grpcConn)
	var je JobExecutor
	switch viper.GetString("jobs.executor") {
	case "kubectl":
		je = &kubectlExecutor{}
	}
	a := newAgent(deploiClient, je)
	if err := a.validateEnvironment(); err != nil {
		fmt.Printf("Failed to validate environment with deploid: %s", err)
		os.Exit(1)
	}
	for {
		jobs, err := a.fetchJobs()
		if err != nil {
			log.Errorf("Failed to fetch jobs: %s", err)
		}
		if jobs != nil {
			for _, job := range jobs {
				if err := a.processJob(job); err != nil {
					log.Errorf("Failed to process job %s: %s", job.Id, err)
				}
			}
		}
		time.Sleep(time.Duration(viper.GetInt("jobs.checkInterval")) * time.Second)
	}
}

func setupConfig() {
	viper.SetConfigName("agent")
	viper.SetDefault("TLS.useTLS", false)
	viper.SetDefault("jobs.checkInterval", 10)
	if err := viper.BindEnv("DEPLOI_AGENT_CONFIG_PATH"); err != nil {
		fmt.Printf("Failed to handle environment variable: %s", err)
		os.Exit(1)
	}
	configPath := viper.GetString("DEPLOI_AGENT_CONFIG_PATH")
	if configPath != "" {
		viper.AddConfigPath(configPath)
	}
	viper.AddConfigPath("/etc/deploi-agent")
	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf("Failed to read config: %s", err)
		os.Exit(1)
	}
	if viper.GetString("environment") == "" {
		fmt.Println("The 'environment' setting is required to be present in the configuration file.")
		os.Exit(1)
	}
}
