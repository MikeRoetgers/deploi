package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"time"

	"github.com/MikeRoetgers/deploi"
	"github.com/MikeRoetgers/deploi/protobuf"
)

type kubectlExecutor struct {
	volumePath string
}

func (k *kubectlExecutor) ProcessJob(job *protobuf.Job) (output []byte, err error) {
	if k.volumePath == "" {
		k.volumePath = "/tmp"
	}
	dir := fmt.Sprintf("%s/%s", k.volumePath, job.Id)
	if err = os.MkdirAll(dir, os.FileMode(int(0755))); err != nil {
		return output, fmt.Errorf("Failed to create directory: %s", err)
	}
	if err = os.Chdir(dir); err != nil {
		return output, fmt.Errorf("Failed to change to directory %s: %s", dir, err)
	}
	if err = ioutil.WriteFile(deploi.ManifestFile, []byte(job.Build.Files[deploi.ManifestFile]), os.FileMode(int(0644))); err != nil {
		return output, fmt.Errorf("Failed to write manifest file: %s", err)
	}
	cmd := exec.Command("/usr/local/bin/kubectl", "describe", "-f", fmt.Sprintf("%s/%s", dir, "manifest"))
	cmd.Env = os.Environ()
	if err = cmd.Run(); err != nil {
		return output, fmt.Errorf("Failed to apply manifest file via kubectl: %s", err)
	}
	time.Sleep(time.Second)
	cmd = exec.Command("/usr/local/bin/kubectl", "describe", "-f", fmt.Sprintf("%s/%s", dir, "manifest"))
	cmd.Env = os.Environ()
	if output, err = cmd.Output(); err != nil {
		return output, fmt.Errorf("Failed to fetch kubectl describe output: %s", err)
	}
	return
}
