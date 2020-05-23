package peanut

import (
	"fmt"
	"os/exec"
)

type genericGitSourceFilesystem struct {
	repositoryURL string
	revision      string
}

// NewGenericGitSourceFS creates a generic-git source filesystem which interacts with Git using the shell command.
func NewGenericGitSourceFS(repositoryURL string, revision string) SourceFilesystem {
	sourceFs := genericGitSourceFilesystem{repositoryURL, revision}
	return &sourceFs
}

func (sourceFs *genericGitSourceFilesystem) FetchAll(destination string) error {
	cloneCommand := exec.Command("git", "clone", sourceFs.repositoryURL, destination)
	err := cloneCommand.Run()
	if err != nil {
		return fmt.Errorf("error cloning repository \"%s\": %s", sourceFs.repositoryURL, err)
	}

	checkoutCommand := exec.Command("git", "-C", destination, "checkout", sourceFs.revision)
	err = checkoutCommand.Run()
	if err != nil {
		return fmt.Errorf("failed to checkout revision \"%s\": %s", sourceFs.revision, err)
	}

	return nil
}
