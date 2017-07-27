package main

import (
	"fmt"

	"golang.org/x/net/context"

	"github.com/MikeRoetgers/deploi/protobuf"
	"github.com/boltdb/bolt"
	"github.com/golang/protobuf/proto"
)

type server struct {
	db *bolt.DB
}

func newServer(db *bolt.DB) *server {
	return &server{
		db: db,
	}
}

func (s *server) RegisterNewBuild(ctx context.Context, req *protobuf.NewBuildRequest) (*protobuf.StandardResponse, error) {
	res := &protobuf.StandardResponse{
		Header: &protobuf.ResponseHeader{
			Success: true,
		},
	}
	err := s.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(ProjectBucket)
		proj := getProject(b, req.GetBuild().GetProjectName())
		proj.Builds = append(proj.Builds, req.Build)
		storeProject(b, proj)
		return nil
	})
	if err != nil {
		log.Errorf("Failed to register new build: %s", err)
		addInternalError(res.Header)
		return res, nil
	}
	return res, nil
}

func (s *server) GetNextJob(context.Context, *protobuf.NextJobRequest) (*protobuf.NextJobResponse, error) {
	return nil, nil
}

func (s *server) MarkJobDone(context.Context, *protobuf.JobDoneRequest) (*protobuf.StandardResponse, error) {
	return nil, nil
}

func (s *server) GetProjects(ctx context.Context, req *protobuf.StandardRequest) (*protobuf.GetProjectsResponse, error) {
	res := &protobuf.GetProjectsResponse{
		Header: &protobuf.ResponseHeader{
			Success: true,
		},
		Projects: []string{},
	}
	err := s.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(ProjectBucket)
		err := b.ForEach(func(k []byte, _ []byte) error {
			res.Projects = append(res.Projects, string(k))
			return nil
		})
		if err != nil {
			return fmt.Errorf("Failed to load projects: %s", err)
		}
		return nil
	})
	if err != nil {
		addInternalError(res.Header)
		log.Errorf("Failed to load list of projects: %s", err)
		return res, nil
	}
	return res, nil
}

func (s *server) GetBuilds(ctx context.Context, req *protobuf.GetBuildsRequest) (*protobuf.GetBuildsResponse, error) {
	res := &protobuf.GetBuildsResponse{
		Header: &protobuf.ResponseHeader{
			Success: true,
		},
		Builds: []*protobuf.Build{},
	}
	err := s.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(ProjectBucket)
		p := getProject(b, req.ProjectName)
		res.Builds = p.Builds
		return nil
	})
	if err != nil {
		addInternalError(res.Header)
		log.Errorf("Failed to load list of builds: %s", err)
		return res, nil
	}
	return res, nil
}

func (s *server) DeployBuild(context.Context, *protobuf.DeployRequest) (*protobuf.DeployResponse, error) {

	return nil, nil
}

func (s *server) AutomateDeployment(context.Context, *protobuf.AutomationRequest) (*protobuf.AutomationResponse, error) {
	return nil, nil
}

func (s *server) RegisterEnvironment(ctx context.Context, req *protobuf.RegisterEnvironmentRequest) (*protobuf.StandardResponse, error) {
	res := &protobuf.StandardResponse{
		Header: &protobuf.ResponseHeader{
			Success: true,
		},
	}
	err := s.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(EnvironmentBucket)
		env := getEnvironment(b, req.Environment.Name)
		env.Namespaces = append(env.Namespaces, req.Environment.Namespaces...)
		if err := storeEnvironment(b, env); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		log.Errorf("Failed to register environment %s: %s", req.Environment.Name, err)
		addInternalError(res.Header)
		return res, nil
	}
	return res, nil
}

func (s *server) GetEnvironments(context.Context, *protobuf.StandardRequest) (*protobuf.GetEnvironmentResponse, error) {
	res := &protobuf.GetEnvironmentResponse{
		Header: &protobuf.ResponseHeader{
			Success: true,
		},
		Environments: []*protobuf.Environment{},
	}
	err := s.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(EnvironmentBucket)
		b.ForEach(func(_ []byte, v []byte) error {
			env := &protobuf.Environment{}
			if err := proto.Unmarshal(v, env); err != nil {
				return fmt.Errorf("Failed to unmarshal environment: %s", err)
			}
			res.Environments = append(res.Environments, env)
			return nil
		})
		return nil
	})
	if err != nil {
		log.Errorf("Failed to load environments: %s", err)
		addInternalError(res.Header)
		return res, nil
	}
	return res, nil
}

func (s *server) DeleteEnvironment(ctx context.Context, req *protobuf.DeleteEnvironmentRequest) (*protobuf.StandardResponse, error) {
	res := &protobuf.StandardResponse{
		Header: &protobuf.ResponseHeader{
			Success: true,
		},
	}
	err := s.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(EnvironmentBucket)
		if len(req.Environment.Namespaces) == 0 {
			if err := b.Delete([]byte(req.Environment.Name)); err != nil {
				return err
			}
			return nil
		}
		env := getEnvironment(b, req.Environment.Name)
		toDelete := map[string]struct{}{}
		for _, ns := range req.Environment.Namespaces {
			toDelete[ns] = struct{}{}
		}
		for k, v := range env.Namespaces {
			if _, ok := toDelete[v]; ok {
				env.Namespaces = append(env.Namespaces[:k], env.Namespaces[k+1:]...)
			}
		}
		if err := storeEnvironment(b, env); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		log.Errorf("Failed to delete environment or namespace: %s", err)
		addInternalError(res.Header)
		return res, nil
	}
	return res, nil
}

func addError(header *protobuf.ResponseHeader, code, message string) {
	header.Success = false
	header.Errors = append(header.Errors, &protobuf.Error{
		Code:    code,
		Message: message,
	})
}

func addInternalError(header *protobuf.ResponseHeader) {
	header.Success = false
	header.Errors = append(header.Errors, &protobuf.Error{
		Code:    "INTERNAL_ERROR",
		Message: "An internal server error occured",
	})
}
