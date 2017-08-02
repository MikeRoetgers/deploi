package main

import (
	"sort"
	"time"

	"github.com/MikeRoetgers/deploi/protobuf"
	"github.com/boltdb/bolt"
	"github.com/golang/protobuf/proto"
)

func enforceRetention(db *bolt.DB) {
	for {
		cleanupJobs(db)
		cleanupBuilds(db)
		time.Sleep(10 * time.Minute)
	}
}

type byDate []*protobuf.Job

func (a byDate) Len() int           { return len(a) }
func (a byDate) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a byDate) Less(i, j int) bool { return a[i].FinishedAt < a[j].FinishedAt }

func cleanupJobs(db *bolt.DB) {
	jobs := []*protobuf.Job{}
	db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(DoneJobBucket)
		b.ForEach(func(_ []byte, v []byte) error {
			job := &protobuf.Job{}
			if err := proto.Unmarshal(v, job); err != nil {
				return nil
			}
			jobs = append(jobs, job)
			return nil
		})
		if len(jobs) <= config.Retention.Jobs {
			return nil
		}
		sort.Sort(sort.Reverse(byDate(jobs)))
		for _, j := range jobs[config.Retention.Jobs:] {
			b.Delete([]byte(j.Id))
		}
		return nil
	})
}

func cleanupBuilds(db *bolt.DB) {
	db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(ProjectBucket)
		b.ForEach(func(k []byte, _ []byte) error {
			p := getProject(b, string(k))
			if p == nil {
				return nil
			}
			var buildsToDelete int
			if num, ok := config.Retention.BuildsPerProject[string(k)]; ok {
				buildsToDelete = len(p.Builds) - num
			} else {
				buildsToDelete = len(p.Builds) - config.Retention.Builds
			}
			if buildsToDelete > 0 {
				p.Builds = p.Builds[buildsToDelete:]
				if err := storeProject(b, p); err != nil {
					log.Errorf("Failed to store cleaned version of project: %s", err)
					return nil
				}
			}
			return nil
		})
		return nil
	})
}
