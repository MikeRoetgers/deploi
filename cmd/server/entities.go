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

func getEnvironment(bucket *bolt.Bucket, name string) *protobuf.Environment {
	val := bucket.Get([]byte(name))
	env := &protobuf.Environment{}
	if len(val) == 0 {
		env.Name = name
		return env
	}
	if err := proto.Unmarshal(val, env); err != nil {
		log.Errorf("Failed to unmarshal environment %s. Entity was reset. Error: %s", name, err)
		env.Name = name
	}
	return env
}

func storeEnvironment(bucket *bolt.Bucket, env *protobuf.Environment) error {
	val, err := proto.Marshal(env)
	if err != nil {
		return fmt.Errorf("Failed to marshal environment: %s", err)
	}
	if err = bucket.Put([]byte(env.Name), val); err != nil {
		return fmt.Errorf("Failed to store environment: %s", err)
	}
	return nil
}
