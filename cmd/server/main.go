package main

import (
	"fmt"
	"net"
	"os"
	"path/filepath"

	"google.golang.org/grpc"

	"github.com/MikeRoetgers/deploi/protobuf"
	"github.com/boltdb/bolt"
	logging "github.com/op/go-logging"
)

var (
	ProjectBucket     = []byte("Projects")
	EnvironmentBucket = []byte("Environments")
	JobBucket         = []byte("Jobs")
	DoneJobBucket     = []byte("DoneJobs")
	AutomationBucket  = []byte("Automations")
	UserBucket        = []byte("Users")
	log               = logging.MustGetLogger("app")
	config            *Config
	JWTKey            = []byte("NbZMecQ3UcgzsnzmNaHTTUTfH6w3XB")
)

func main() {
	var err error
	configFile := os.Getenv("DEPLOID_CONFIG")
	if configFile == "" {
		config = newConfig()
	} else {
		if config, err = newConfigFromFile(configFile); err != nil {
			log.Fatalf("Failed to start deamon. Reason: %s", err)
		}
	}
	if err = os.MkdirAll(filepath.Base(config.Database.Path), os.FileMode(int(0755))); err != nil {
		log.Fatalf("Failed to start daemon. Reason: Could not create database path. %s", err)
	}
	db, err := bolt.Open(config.Database.Path, 0600, nil)
	if err != nil {
		log.Fatalf("Failed to start daemon. Reason: %s", err)
	}
	defer db.Close()

	err = db.Update(func(tx *bolt.Tx) error {
		buckets := [][]byte{ProjectBucket, EnvironmentBucket, JobBucket, DoneJobBucket, AutomationBucket, UserBucket}
		for _, b := range buckets {
			_, txErr := tx.CreateBucketIfNotExists(b)
			if txErr != nil {
				return fmt.Errorf("creating bucket failed: %s", txErr)
			}
		}
		return nil
	})
	if err != nil {
		log.Fatalf("Failed to start daemon. Reason: %s", err)
	}

	s := newServer(db)
	lis, err := net.Listen("tcp", config.ListenAddr)
	if err != nil {
		log.Fatalf("Failed to start daemon. Reason: %v", err)
	}
	go enforceRetention(db)
	grpcServer := grpc.NewServer()
	protobuf.RegisterDeploiServerServer(grpcServer, s)
	grpcServer.Serve(lis)
}
