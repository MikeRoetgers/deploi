package tests

import (
	"flag"
	"fmt"
	"os"
	"testing"

	"google.golang.org/grpc"

	"github.com/MikeRoetgers/deploi/protobuf"
)

var deploiClient protobuf.DeploiServerClient
var runSystemTests = flag.Bool("system", false, "Run system tests?")
var defaultToken string

func TestMain(m *testing.M) {
	flag.Parse()
	if *runSystemTests {
		host := os.Getenv("DEPLOID_HOST")
		if host == "" {
			host = "127.0.0.1:8000"
		}
		grpcConn, err := grpc.Dial(host, grpc.WithInsecure())
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		defer grpcConn.Close()
		deploiClient = protobuf.NewDeploiServerClient(grpcConn)
		userCred := []string{"test@example.org", "123456"}
		if err = createUser(userCred[0], userCred[1]); err != nil {
			fmt.Printf("Failed to create default user: %s", err)
			os.Exit(1)
		}
		if defaultToken, err = login(userCred[0], userCred[1]); err != nil {
			fmt.Printf("Failed to login default user: %s", err)
			os.Exit(1)
		}
	}
	os.Exit(m.Run())
}
