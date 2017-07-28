package main

import (
	"context"
	"fmt"

	"github.com/MikeRoetgers/deploi"
	"github.com/MikeRoetgers/deploi/protobuf"
)

type agent struct {
	deploiClient protobuf.DeploiServerClient
}

func newAgent(c protobuf.DeploiServerClient) *agent {
	return &agent{
		deploiClient: c,
	}
}

func (a *agent) processJob(job *protobuf.Job) error {
	fmt.Printf("ID: %s\nProject: %s\nBuild: %s\nEnv: %s\n\n", job.Id, job.Build.ProjectName, job.Build.BuildId, job.Environment.Name)
	req := &protobuf.JobDoneRequest{
		Header: &protobuf.RequestHeader{},
		Job:    job,
	}
	res, err := a.deploiClient.MarkJobDone(context.Background(), req)
	if err != nil {
		return err
	}
	if !res.Header.Success {
		return deploi.NewResponseError(res.Header)
	}
	return nil
}

func (a *agent) fetchJobs() ([]*protobuf.Job, error) {
	req := &protobuf.NextJobRequest{
		Header:      &protobuf.RequestHeader{},
		Environment: *environment,
	}
	res, err := a.deploiClient.GetNextJobs(context.Background(), req)
	if err != nil {
		return nil, err
	}
	if !res.Header.Success {
		return nil, deploi.NewResponseError(res.Header)
	}
	return res.Jobs, nil
}
