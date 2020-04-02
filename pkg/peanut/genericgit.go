package peanut

import (
	"errors"
	"fmt"
	"os/exec"
)

type genericGitSourceFilesystem struct {
	repositoryUrl string
}

func NewGenericGitSourceFS(repositoryUrl string) SourceFilesystem {
	sourceFs := genericGitSourceFilesystem{repositoryUrl}
	return &sourceFs
}

func (sourceFs *genericGitSourceFilesystem) FetchAll(destination string) error {
	cloneCommand := exec.Command("git", "clone", sourceFs.repositoryUrl, destination)
	err := cloneCommand.Run()

	if err != nil {
		errMessage := fmt.Sprintf("Error cloning repository: %s. %s", sourceFs.repositoryUrl, err)
		return errors.New(errMessage)
	}

	return nil
}
