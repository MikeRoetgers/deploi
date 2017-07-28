package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/MikeRoetgers/deploi/protobuf"
	logging "github.com/op/go-logging"
	"google.golang.org/grpc"
)

var environment = flag.String("environment", "", "Name of the environment")
var deploidHost = flag.String("host", "127.0.0.1:8000", "host:port of the deploid")
var namespaces []string
var log = logging.MustGetLogger("app")

func init() {
	ns := flag.String("namespaces", "", "Comma-separated list of namespaces in the environment")
	namespaces = strings.Split(*ns, ",")
}

func main() {
	flag.Parse()
	if *environment == "" {
		fmt.Printf("An environment name has to be provided\n\n")
		flag.PrintDefaults()
		os.Exit(1)
	}
	grpcConn, err := grpc.Dial(*deploidHost, grpc.WithInsecure())
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer grpcConn.Close()
	deploiClient := protobuf.NewDeploiServerClient(grpcConn)
	a := newAgent(deploiClient)
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
		time.Sleep(10 * time.Second)
	}
}
