package git

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"

	"gitlab.com/slon/shad-go/gitfame/internal/errors"
)

func LsTree(repoPath, revision string) ([]string, error) {
	cmd := exec.Command("git", "ls-tree", "-r", "--name-only", revision)
	cmd.Dir = repoPath

	var stdout bytes.Buffer
	cmd.Stdout = &stdout
	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf(errors.MsgGitLsTree, err)
	}

	output := strings.TrimSpace(stdout.String())
	if output == "" {
		return []string{}, nil
	}

	files := strings.Split(output, "\n")
	return files, nil
}
