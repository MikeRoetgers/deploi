package tests

import (
	"context"
	"testing"

	"github.com/MikeRoetgers/deploi/protobuf"
	"github.com/stretchr/testify/assert"
)

func TestAutomationCRUDOperations(t *testing.T) {
	if !*runSystemTests {
		t.SkipNow()
	}
	project := ""
	branch := ""
	env := &protobuf.Environment{
		Name:       "",
		Namespaces: []string{"default"},
	}
	a := &protobuf.Automation{
		Automation: &protobuf.Automation_BranchAutomation{
			BranchAutomation: &protobuf.BranchAutomation{
				Project:     project,
				Branch:      branch,
				Environment: env,
			},
		},
	}
	if err := createAutomation(a); err != nil {
		t.Fatal(err)
	}

	as, err := getAutomations()
	if err != nil {
		t.Fatal(err)
	}
	myA := findAutomationByProject(as, project)
	if myA == nil {
		t.Fatal("Could not find automation that was just created")
	}

	if err := deleteAutomation(myA.Id); err != nil {
		t.Fatal(err)
	}

	as, err = getAutomations()
	if err != nil {
		t.Fatal(err)
	}
	myA = findAutomationByProject(as, project)
	assert.Nil(t, myA)
}

func TestAutomatedDeploy(t *testing.T) {
	if !*runSystemTests {
		t.SkipNow()
	}
	env := &protobuf.Environment{
		Name:       "testAutomatedDeploy",
		Namespaces: []string{"default"},
	}
	build := &protobuf.Build{
		ProjectName: "testAutomatedDeployProject",
		BuildId:     "1",
		BranchName:  "testAutomatedDeploy",
	}
	ba := &protobuf.BranchAutomation{
		Branch:      build.BranchName,
		Project:     build.ProjectName,
		Environment: env,
	}
	automation := &protobuf.Automation{
		Automation: &protobuf.Automation_BranchAutomation{
			BranchAutomation: ba,
		},
	}

	if err := registerEnvironment(env); err != nil {
		t.Fatal(err)
	}
	if err := createAutomation(automation); err != nil {
		t.Fatal(err)
	}
	if err := registerBuild(build); err != nil {
		t.Fatal(err)
	}
	req := &protobuf.GetJobsRequest{
		Header:  getReqHeader(),
		Pending: true,
	}
	res, err := deploiClient.GetJobs(context.Background(), req)
	if gErr := handleGRPCResponse(res, err); gErr != nil {
		t.Fatal(gErr)
	}

	jobs, err := getNextJobs(env.Name)
	if err != nil {
		t.Fatal(err)
	}
	assert.Len(t, jobs, 1)
}

func findAutomationByProject(automations []*protobuf.Automation, project string) *protobuf.Automation {
	for _, a := range automations {
		ba := a.GetBranchAutomation()
		if ba != nil && ba.Project == project {
			return a
		}
	}
	return nil
}
