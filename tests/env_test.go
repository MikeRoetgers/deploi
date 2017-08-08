package tests

import (
	"testing"

	"github.com/MikeRoetgers/deploi/protobuf"
	"github.com/stretchr/testify/assert"
)

func TestEnvironmentCRUDOperations(t *testing.T) {
	if !*runSystemTests {
		t.SkipNow()
	}
	env := &protobuf.Environment{
		Name:       "envTest",
		Namespaces: []string{"default", "default2", "default3"},
	}
	if err := registerEnvironment(env); err != nil {
		t.Fatal(err)
	}

	myEnv, err := getEnvironment(env.Name)
	if err != nil {
		t.Fatalf("Failed to find environment: %s", err)
	}
	assert.Equal(t, env, myEnv)

	if err := deleteEnvironment(env.Name, env.Namespaces[2]); err != nil {
		t.Error(err)
	}

	myEnv, err = getEnvironment(env.Name)
	if err != nil {
		t.Fatalf("Failed to find environment: %s", err)
	}
	assert.Len(t, myEnv.Namespaces, 2)
	assert.NotContains(t, myEnv.Namespaces, env.Namespaces[2])

	if err := deleteEnvironment(env.Name); err != nil {
		t.Error(err)
	}

	myEnv, _ = getEnvironment(env.Name)
	assert.Nil(t, myEnv)
}
