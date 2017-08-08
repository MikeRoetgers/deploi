package tests

import (
	"testing"

	"github.com/MikeRoetgers/deploi"
	"github.com/MikeRoetgers/deploi/protobuf"
	"github.com/stretchr/testify/assert"
)

func TestRegisterBuild(t *testing.T) {
	if !*runSystemTests {
		t.SkipNow()
	}
	var err error
	build := &protobuf.Build{
		ProjectName:    "registerTest",
		BuildId:        "123",
		BuildURL:       "buildUrl",
		BuildSystemURL: "buildSystemUrl",
		BranchName:     "thebranch",
		Files: map[string]string{
			deploi.ManifestFile: "This is a test",
		},
	}
	if err = registerBuild(build); err != nil {
		t.Fatal(err)
	}
	builds, err := getBuilds(build.ProjectName)
	if err != nil {
		t.Fatal(err)
	}
	if len(builds) != 1 {
		t.Fatalf("Expected 1 build, got %d", len(builds))
	}
	b := builds[0]
	assert.Equal(t, build.BranchName, b.BranchName)
	assert.Equal(t, build.BuildId, b.BuildId)
	assert.Equal(t, build.BuildURL, b.BuildURL)
	assert.Equal(t, build.BuildSystemURL, b.BuildSystemURL)
	assert.Equal(t, build.Files, b.Files)
}

func TestDeploymentProcess(t *testing.T) {
	if !*runSystemTests {
		t.SkipNow()
	}
	build := &protobuf.Build{
		ProjectName: "deploymentTest",
		BuildId:     "1",
	}
	env := &protobuf.Environment{
		Name:       "deploymentTestEnv",
		Namespaces: []string{"default"},
	}

	// register environment we want to deploy to
	if err := registerEnvironment(env); err != nil {
		t.Fatal(err)
	}

	// register a new build that can be deployed
	if err := registerBuild(build); err != nil {
		t.Fatal(err)
	}

	// mark build to be deployed to the freshly registered environment
	if err := deployBuild(build.ProjectName, build.BuildId, env.Name, env.Namespaces[0]); err != nil {
		t.Fatal(err)
	}

	// agent of the created environment asks for jobs, should get our newly created job
	jobs, err := getNextJobs(env.Name)
	if err != nil {
		t.Fatal(err)
	}

	assert.Len(t, jobs, 1)
	assert.Equal(t, build.BuildId, jobs[0].Build.BuildId)
	assert.Equal(t, build.ProjectName, jobs[0].Build.ProjectName)

	// mark job as done, should not be handed out again
	if err := markJobAsDone(jobs[0]); err != nil {
		t.Fatal(err)
	}

	jobs, err = getNextJobs(env.Name)
	if err != nil {
		t.Fatal(err)
	}

	assert.Empty(t, jobs)
}
