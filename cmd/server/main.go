package main

import (
	"fmt"
	"net"

	"google.golang.org/grpc"

	"github.com/MikeRoetgers/deploi/protobuf"
	"github.com/boltdb/bolt"
	logging "github.com/op/go-logging"
)

var (
	ProjectBucket = []byte("Projects")
	log           = logging.MustGetLogger("app")
)

func main() {
	db, err := bolt.Open("my.db", 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	err = db.Update(func(tx *bolt.Tx) error {
		_, txErr := tx.CreateBucketIfNotExists(ProjectBucket)
		if txErr != nil {
			return fmt.Errorf("create bucket: %s", err)
		}
		return nil
	})
	if err != nil {
		log.Fatalf("Failed to initiate database: %s", err)
	}

	s := newServer(db)

	lis, err := net.Listen("tcp", ":8000")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	protobuf.RegisterDeploiServerServer(grpcServer, s)
	grpcServer.Serve(lis)
}
