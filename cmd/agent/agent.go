package main

import (
	"context"

	"github.com/MikeRoetgers/deploi/protobuf"
	"github.com/spf13/viper"
)

type agent struct {
	deploiClient protobuf.DeploiServerClient
	executor     JobExecutor
}

func newAgent(c protobuf.DeploiServerClient, e JobExecutor) *agent {
	return &agent{
		deploiClient: c,
		executor:     e,
	}
}

type JobExecutor interface {
	ProcessJob(*protobuf.Job) ([]byte, error)
}

func (a *agent) processJob(job *protobuf.Job) error {
	output, err := a.executor.ProcessJob(job)
	if err != nil {
		log.Errorf("Job: %s | Error: %s", job.Id, err)
	}
	job.Output = output
	req := &protobuf.JobDoneRequest{
		Header: &protobuf.RequestHeader{},
		Job:    job,
	}
	res, err := a.deploiClient.MarkJobDone(context.Background(), req)
	if err != nil {
		return err
	}
	if !res.Header.Success {
		return newResponseError(res.Header)
	}
	return nil
}

func (a *agent) fetchJobs() ([]*protobuf.Job, error) {
	req := &protobuf.NextJobRequest{
		Header:      &protobuf.RequestHeader{},
		Environment: viper.GetString("environment"),
	}
	res, err := a.deploiClient.GetNextJobs(context.Background(), req)
	if err != nil {
		return nil, err
	}
	if !res.Header.Success {
		return nil, newResponseError(res.Header)
	}
	return res.Jobs, nil
}
