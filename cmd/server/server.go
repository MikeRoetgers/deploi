package main

import (
	"bytes"
	"fmt"
	"time"

	"golang.org/x/net/context"

	"github.com/MikeRoetgers/deploi"
	"github.com/MikeRoetgers/deploi/protobuf"
	"github.com/boltdb/bolt"
	"github.com/golang/protobuf/proto"
	uuid "github.com/satori/go.uuid"
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
		proj := getOrCreateProject(b, req.GetBuild().GetProjectName())
		proj.Builds = append(proj.Builds, req.Build)
		if err := storeProject(b, proj); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		log.Errorf("Failed to register new build: %s", err)
		addInternalError(res.Header)
		return res, nil
	}
	return res, nil
}

func (s *server) GetNextJobs(ctx context.Context, req *protobuf.NextJobRequest) (*protobuf.NextJobResponse, error) {
	res := &protobuf.NextJobResponse{
		Header: &protobuf.ResponseHeader{
			Success: true,
		},
		Jobs: []*protobuf.Job{},
	}
	err := s.db.View(func(tx *bolt.Tx) error {
		c := tx.Bucket(JobBucket).Cursor()
		prefix := []byte(req.Environment)
		for k, v := c.Seek(prefix); k != nil && bytes.HasPrefix(k, prefix); k, v = c.Next() {
			job := &protobuf.Job{}
			if err := proto.Unmarshal(v, job); err != nil {
				return err
			}
			res.Jobs = append(res.Jobs, job)
		}
		return nil
	})
	if err != nil {
		log.Errorf("Failed to load pending jobs: %s", err)
		addInternalError(res.Header)
		return res, nil
	}
	return res, nil
}

func (s *server) MarkJobDone(ctx context.Context, req *protobuf.JobDoneRequest) (*protobuf.StandardResponse, error) {
	log.Debugf("Marking job %s as done", req.Job.Id)
	res := &protobuf.StandardResponse{
		Header: &protobuf.ResponseHeader{
			Success: true,
		},
	}
	err := s.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(JobBucket)
		djb := tx.Bucket(DoneJobBucket)
		job := getPendingJob(b, req.Job.Id, req.Job.Environment.Name)
		if job == nil {
			addError(res.Header, "JOB_MISSING", "The provided job could not be found in the database")
			return &AlreadyHandledError{}
		}
		job.FinishedAt = time.Now().Unix()
		job.Output = req.Job.Output
		if err := storeJob(djb, job); err != nil {
			log.Errorf("Failed to store job %s in done bucket: %s", job.Id, err)
			addInternalError(res.Header)
			return &AlreadyHandledError{}
		}
		if err := b.Delete([]byte(fmt.Sprintf("%s_%s", job.Environment.Name, job.Id))); err != nil {
			log.Errorf("Failed to delete job %s from job queue: %s", job.Id, err)
			addInternalError(res.Header)
			return &AlreadyHandledError{}
		}
		return nil
	})
	if _, ok := err.(*AlreadyHandledError); !ok && err != nil {
		log.Errorf("Failed to mark job as done: %s", err)
		addInternalError(res.Header)
		return res, nil
	}
	log.Debugf("Marking job %s worked", req.Job.Id)
	return res, nil
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
		p := getOrCreateProject(b, req.ProjectName)
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

func (s *server) DeployBuild(ctx context.Context, req *protobuf.DeployRequest) (*protobuf.DeployResponse, error) {
	res := &protobuf.DeployResponse{
		Header: &protobuf.ResponseHeader{
			Success: true,
		},
	}
	job := &protobuf.Job{
		Id: uuid.NewV4().String(),
	}

	// Load and validate existance of entities needed to compose a job
	err := s.db.View(func(tx *bolt.Tx) error {
		pb := tx.Bucket(ProjectBucket)
		eb := tx.Bucket(EnvironmentBucket)
		proj := getProject(pb, req.Project)
		if proj == nil {
			addError(res.Header, "PROJECT_MISSING", "The supplied project does not exist in the database.")
		} else {
			for _, build := range proj.Builds {
				if build.BuildId == req.BuildId {
					job.Build = build
					break
				}
			}
			if job.Build == nil {
				addError(res.Header, "BUILD_MISSING", "The supplied build does not exist in the project.")
			}
		}

		env := getEnvironment(eb, req.Environment)
		if env == nil {
			addError(res.Header, "ENVIRONMENT_MISSING", "The supplied environment does not exist in the database.")
		}
		if (env != nil) && (!environmentHasNamespace(env, req.Namespace)) {
			addError(res.Header, "ENVIRONMENT_NAMESPACE_MISSING", "The supplied namespace does not exist in the environment.")
		}

		if !res.Header.Success {
			return &AlreadyHandledError{}
		}
		job.Environment = &protobuf.Environment{
			Name:       env.Name,
			Namespaces: []string{req.Namespace},
		}
		job.CreatedAt = time.Now().Unix()
		return nil
	})
	if err != nil {
		return res, nil
	}

	// overwrite manifest in build in case one was provided in deploy request
	if manifest, ok := req.Files[deploi.ManifestFile]; ok {
		job.Build.Files[deploi.ManifestFile] = manifest
	}

	// write job to database
	err = s.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(JobBucket)
		if err := storePendingJob(b, job); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		log.Errorf("Failed to store new job: %s", err)
		addInternalError(res.Header)
		return res, nil
	}
	return res, nil
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
		env := getOrCreateEnvironment(b, req.Environment.Name)
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

func (s *server) GetJobs(ctx context.Context, req *protobuf.GetJobsRequest) (*protobuf.GetJobsResponse, error) {
	res := &protobuf.GetJobsResponse{
		Header: &protobuf.ResponseHeader{
			Success: true,
		},
	}
	err := s.db.View(func(tx *bolt.Tx) error {
		var b *bolt.Bucket
		if req.Pending && req.Id == "" {
			b = tx.Bucket(JobBucket)
		} else {
			b = tx.Bucket(DoneJobBucket)
		}
		if req.Id != "" {
			job := getJob(b, req.Id)
			if job == nil {
				addError(res.Header, "JOB_NOT_FOUND", "Job could not be found in the database")
				return nil
			}
			res.Jobs = []*protobuf.Job{job}
			return nil
		}
		jobs, err := getJobs(b)
		if err != nil {
			return err
		}
		res.Jobs = jobs
		return nil
	})
	if err != nil {
		log.Errorf("Failed to load jobs: %s", err)
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

type AlreadyHandledError struct {
}

func (e *AlreadyHandledError) Error() string { return "" }
