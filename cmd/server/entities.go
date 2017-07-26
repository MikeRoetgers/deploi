package main

import (
	"fmt"

	"github.com/MikeRoetgers/deploi/protobuf"
	"github.com/boltdb/bolt"
	"github.com/golang/protobuf/proto"
)

func getProject(bucket *bolt.Bucket, name string) *protobuf.Project {
	val := bucket.Get([]byte(name))
	proj := &protobuf.Project{}
	if len(val) == 0 {
		proj.ProjectName = name
		return proj
	}
	if err := proto.Unmarshal(val, proj); err != nil {
		log.Errorf("Failed to unmarshal project %s. Entity was reset. Error: %s", name, err)
		proj.ProjectName = name
	}
	return proj
}

func storeProject(bucket *bolt.Bucket, project *protobuf.Project) error {
	val, err := proto.Marshal(project)
	if err != nil {
		return fmt.Errorf("Failed to marshal project: %s", err)
	}
	if err = bucket.Put([]byte(project.ProjectName), val); err != nil {
		return fmt.Errorf("Failed to store project: %s", err)
	}
	return nil
}
