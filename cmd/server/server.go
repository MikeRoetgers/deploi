package main

import (
	"fmt"

	"golang.org/x/net/context"

	"github.com/MikeRoetgers/deploi/protobuf"
	"github.com/boltdb/bolt"
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
