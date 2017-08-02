package main

import (
	"fmt"
	"os"
	"testing"

	"github.com/MikeRoetgers/deploi/protobuf"
	"github.com/boltdb/bolt"
	uuid "github.com/satori/go.uuid"
)

func TestMain(m *testing.M) {
	config = newConfig()
	os.Exit(m.Run())
}

func initDb(t *testing.T) *bolt.DB {
	dbFile := fmt.Sprintf("%s/%s", os.TempDir(), uuid.NewV4().String())
	db, err := bolt.Open(dbFile, 0600, nil)
	if err != nil {
		t.Fatal(err)
	}
	err = db.Update(func(tx *bolt.Tx) error {
		buckets := [][]byte{ProjectBucket, EnvironmentBucket, JobBucket, DoneJobBucket, AutomationBucket}
		for _, b := range buckets {
			_, txErr := tx.CreateBucketIfNotExists(b)
			if txErr != nil {
				return txErr
			}
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	return db
}

func destroyDb(db *bolt.DB) {
	p := db.Path()
	db.Close()
	os.Remove(p)
}

func TestJobCleanup(t *testing.T) {
	db := initDb(t)
	config.Retention.Jobs = 5
	db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(DoneJobBucket)
		for i := 0; i < 10; i++ {
			if err := storeJob(b, &protobuf.Job{
				Id:         fmt.Sprintf("%d", i),
				FinishedAt: int64(i),
			}); err != nil {
				t.Fatalf("Failed to write test data: %s", err)
			}
		}
		return nil
	})
	cleanupJobs(db)
	db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(DoneJobBucket)
		keys := map[string]struct{}{
			"5": struct{}{},
			"6": struct{}{},
			"7": struct{}{},
			"8": struct{}{},
			"9": struct{}{},
		}
		b.ForEach(func(k []byte, _ []byte) error {
			if _, ok := keys[string(k)]; !ok {
				t.Errorf("Unexpected key %s", k)
			}
			return nil
		})
		return nil
	})
	destroyDb(db)
}

func TestBuildCleanup(t *testing.T) {
	db := initDb(t)
	config.Retention.Builds = 5
	builds := []*protobuf.Build{}
	for i := 0; i < 10; i++ {
		builds = append(builds, &protobuf.Build{
			BuildId:   fmt.Sprintf("%d", i),
			CreatedAt: int64(i),
		})
	}
	p := &protobuf.Project{
		ProjectName: "test",
		Builds:      builds,
	}
	db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(ProjectBucket)
		if err := storeProject(b, p); err != nil {
			t.Fatal(err)
		}
		return nil
	})
	cleanupBuilds(db)
	db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(ProjectBucket)
		p := getProject(b, "test")
		if p == nil {
			t.Error("test project is missing")
			return nil
		}
		if len(p.Builds) != 5 {
			t.Errorf("Expected 5 builds, got %d", len(p.Builds))
			return nil
		}
		expected := map[string]struct{}{
			"5": struct{}{},
			"6": struct{}{},
			"7": struct{}{},
			"8": struct{}{},
			"9": struct{}{},
		}
		for _, b := range p.Builds {
			if _, ok := expected[b.BuildId]; !ok {
				t.Errorf("Unexpected build %s", b.BuildId)
			}
		}
		return nil
	})
	destroyDb(db)
}
