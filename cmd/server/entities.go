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
		return nil
	}
	if err := proto.Unmarshal(val, proj); err != nil {
		log.Errorf("Failed to unmarshal project %s. Error: %s", name, err)
		return nil
	}
	return proj
}

func getOrCreateProject(bucket *bolt.Bucket, name string) *protobuf.Project {
	res := getProject(bucket, name)
	if res == nil {
		return &protobuf.Project{
			ProjectName: name,
		}
	}
	return res
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
		return nil
	}
	if err := proto.Unmarshal(val, env); err != nil {
		log.Errorf("Failed to unmarshal environment %s. Error: %s", name, err)
		return nil
	}
	return env
}

func getOrCreateEnvironment(bucket *bolt.Bucket, name string) *protobuf.Environment {
	env := getEnvironment(bucket, name)
	if env == nil {
		return &protobuf.Environment{
			Name: name,
		}
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

func environmentHasNamespace(env *protobuf.Environment, namespace string) bool {
	for _, v := range env.Namespaces {
		if v == namespace {
			return true
		}
	}
	return false
}

// Uses a different key pattern (env_jobid) to make the bucket easily searchable
func storePendingJob(bucket *bolt.Bucket, job *protobuf.Job) error {
	val, err := proto.Marshal(job)
	if err != nil {
		return fmt.Errorf("Failed to marshal job: %s", err)
	}
	if err = bucket.Put([]byte(job.Environment.Name+"_"+job.Id), val); err != nil {
		return fmt.Errorf("Failed to store job: %s", err)
	}
	return nil
}

func storeJob(bucket *bolt.Bucket, job *protobuf.Job) error {
	val, err := proto.Marshal(job)
	if err != nil {
		return fmt.Errorf("Failed to marshal job: %s", err)
	}
	if err = bucket.Put([]byte(job.Id), val); err != nil {
		return fmt.Errorf("Failed to store job: %s", err)
	}
	return nil
}

func getPendingJob(bucket *bolt.Bucket, id, environment string) *protobuf.Job {
	val := bucket.Get([]byte(fmt.Sprintf("%s_%s", environment, id)))
	job := &protobuf.Job{}
	if len(val) == 0 {
		return nil
	}
	if err := proto.Unmarshal(val, job); err != nil {
		log.Errorf("Failed to unmarshal job %s. Error: %s", id, err)
		return nil
	}
	return job
}

func getJob(bucket *bolt.Bucket, id string) *protobuf.Job {
	val := bucket.Get([]byte(id))
	job := &protobuf.Job{}
	if len(val) == 0 {
		return nil
	}
	if err := proto.Unmarshal(val, job); err != nil {
		log.Errorf("Failed to unmarshal job %s. Error: %s", id, err)
		return nil
	}
	return job
}
