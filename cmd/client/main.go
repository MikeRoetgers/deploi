package main

import (
	"fmt"
	"os"

	"github.com/MikeRoetgers/deploi/cmd/client/cmd"
	"github.com/MikeRoetgers/deploi/protobuf"
	"google.golang.org/grpc"
)

func main() {
	grpcConn, err := grpc.Dial("localhost:8000", grpc.WithInsecure())
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer grpcConn.Close()
	cmd.DeploiClient = protobuf.NewDeploiServerClient(grpcConn)
	cmd.Execute()
}
