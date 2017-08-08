package tests

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/MikeRoetgers/deploi/protobuf"
)

func createUser(email, password string) error {
	res, err := deploiClient.CreateUser(context.Background(), &protobuf.CreateUserRequest{
		Header:   &protobuf.RequestHeader{},
		Email:    email,
		Password: password,
	})
	if gErr := handleGRPCResponse(res, err); gErr != nil {
		return gErr
	}
	return nil
}

func login(email, password string) (string, error) {
	res, err := deploiClient.Login(context.Background(), &protobuf.LoginRequest{
		Header:   &protobuf.RequestHeader{},
		Username: email,
		Password: password,
	})
	if gErr := handleGRPCResponse(res, err); gErr != nil {
		return "", gErr
	}
	return res.Token, nil
}

func registerBuild(build *protobuf.Build) error {
	req := &protobuf.NewBuildRequest{
		Header: getReqHeader(),
		Build:  build,
	}
	res, err := deploiClient.RegisterNewBuild(context.Background(), req)
	if gErr := handleGRPCResponse(res, err); gErr != nil {
		return fmt.Errorf("Failed to register new build: %s", gErr)
	}
	return nil
}

func getBuilds(projectName string) ([]*protobuf.Build, error) {
	res, err := deploiClient.GetBuilds(context.Background(), &protobuf.GetBuildsRequest{
		Header:      getReqHeader(),
		ProjectName: projectName,
	})
	if gErr := handleGRPCResponse(res, err); gErr != nil {
		return nil, fmt.Errorf("Failed to load projects: %s", gErr)
	}
	return res.Builds, nil
}

func registerEnvironment(env *protobuf.Environment) error {
	req := &protobuf.RegisterEnvironmentRequest{
		Header:      getReqHeader(),
		Environment: env,
	}
	res, err := deploiClient.RegisterEnvironment(context.Background(), req)
	if gErr := handleGRPCResponse(res, err); gErr != nil {
		return fmt.Errorf("Failed to register env: %s", gErr)
	}
	return nil
}

func getEnvironments() ([]*protobuf.Environment, error) {
	res, err := deploiClient.GetEnvironments(context.Background(), &protobuf.StandardRequest{
		Header: getReqHeader(),
	})
	if gErr := handleGRPCResponse(res, err); gErr != nil {
		return nil, fmt.Errorf("Failed to get envs: %s", gErr)
	}
	return res.Environments, nil
}

func getEnvironment(name string) (*protobuf.Environment, error) {
	envs, err := getEnvironments()
	if err != nil {
		return nil, err
	}
	for _, e := range envs {
		if e.Name == name {
			return e, nil
		}
	}
	return nil, errors.New("Could not find requested env in database")
}

func deleteEnvironment(name string, namespaces ...string) error {
	req := &protobuf.DeleteEnvironmentRequest{
		Header: getReqHeader(),
		Environment: &protobuf.Environment{
			Name:       name,
			Namespaces: namespaces,
		},
	}
	res, err := deploiClient.DeleteEnvironment(context.Background(), req)
	if gErr := handleGRPCResponse(res, err); gErr != nil {
		return fmt.Errorf("Failed to register env: %s", gErr)
	}
	return nil
}

func deployBuild(project, buildId, env, namespace string) error {
	req := &protobuf.DeployRequest{
		Header:      getReqHeader(),
		Project:     project,
		BuildId:     buildId,
		Environment: env,
		Namespace:   namespace,
	}
	res, err := deploiClient.DeployBuild(context.Background(), req)
	if err := handleGRPCResponse(res, err); err != nil {
		return fmt.Errorf("Failed to deploy: %s", err)
	}
	return nil
}

func getNextJobs(env string) ([]*protobuf.Job, error) {
	req := &protobuf.NextJobRequest{
		Header:      getReqHeader(),
		Environment: env,
	}
	res, err := deploiClient.GetNextJobs(context.Background(), req)
	if err := handleGRPCResponse(res, err); err != nil {
		return nil, fmt.Errorf("Failed to get next jobs: %s", err)
	}
	return res.Jobs, nil
}

func markJobAsDone(job *protobuf.Job) error {
	req := &protobuf.JobDoneRequest{
		Header: getReqHeader(),
		Job:    job,
	}
	res, err := deploiClient.MarkJobDone(context.Background(), req)
	if err := handleGRPCResponse(res, err); err != nil {
		return fmt.Errorf("Failed to mark job as done: %s", err)
	}
	return nil
}

func getReqHeader() *protobuf.RequestHeader {
	return &protobuf.RequestHeader{
		Token: defaultToken,
	}
}

func handleGRPCResponse(res interface{}, err error) error {
	if err != nil {
		return fmt.Errorf("RPC request failed: %s", err)
	}
	if resp, ok := res.(protobuf.Response); ok {
		header := resp.GetHeader()
		if !header.Success {
			strs := []string{}
			for _, er := range header.Errors {
				strs = append(strs, fmt.Sprintf("Code: %s | Message: %s\n", er.Code, er.Message))
			}
			return errors.New(strings.Join(strs, " -- "))
		}
	}
	return nil
}
