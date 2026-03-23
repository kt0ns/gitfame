package git

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os/exec"
	"strconv"
	"strings"

	"gitlab.com/slon/shad-go/gitfame/internal/errors"
)

const (
	HashLen        = 40
	TargetAuthor   = "author "
	TargetCommiter = "committer "
	SystemLineSep  = '\n'
)

type BlamePortion struct {
	CommitHash string
	NumLines   int
}

type FileStats struct {
	FilePath string
	Commits  map[string]CommitInfo
	Chunks   []BlamePortion
}

type CommitInfo struct {
	Author string
}

func BlameFile(repoPath, revision, filePath string, useCommitter bool) (FileStats, error) {
	cmd := exec.Command("git", "blame", "--incremental", revision, "--", filePath)
	cmd.Dir = repoPath

	var stdout bytes.Buffer
	cmd.Stdout = &stdout

	if err := cmd.Run(); err != nil {
		return FileStats{}, fmt.Errorf(errors.MsgGitBlame, filePath, err)
	}

	if stdout.Len() == 0 {
		return logEmptyFile(repoPath, revision, filePath, useCommitter)
	}

	return parseIncremental(&stdout, filePath, useCommitter), nil
}

func logEmptyFile(repoPath, revision, filePath string, useCommitter bool) (FileStats, error) {
	format := "%H %an"
	if useCommitter {
		format = "%H %cn"
	}

	cmd := exec.Command("git", "log", "-1", "--format="+format, revision, "--", filePath)
	cmd.Dir = repoPath

	var stdout bytes.Buffer
	cmd.Stdout = &stdout

	if err := cmd.Run(); err != nil {
		return FileStats{}, fmt.Errorf(errors.MsgGitLog, filePath, err)
	}

	if stdout.Len() == 0 {
		return FileStats{FilePath: filePath}, nil
	}

	output, _ := stdout.ReadString(SystemLineSep)
	output = strings.TrimSpace(output)

	parts := strings.SplitN(output, " ", 2)
	hash := parts[0]
	author := parts[1]

	return FileStats{
		FilePath: filePath,
		Commits: map[string]CommitInfo{
			hash: {Author: author},
		},
		Chunks: []BlamePortion{{
			CommitHash: hash,
			NumLines:   0,
		}},
	}, nil
}

func parseIncremental(r io.Reader, filePath string, useCommitter bool) FileStats {
	scanner := bufio.NewScanner(r)

	commits := make(map[string]CommitInfo)
	var chunks []BlamePortion

	var currentCommit string

	targetField := TargetAuthor
	if useCommitter {
		targetField = TargetCommiter
	}

	for scanner.Scan() {
		line := scanner.Text()

		parts := strings.Split(line, " ")

		if len(parts[0]) == HashLen {
			currentCommit = parts[0]

			if n, err := strconv.Atoi(parts[3]); err == nil {
				chunks = append(chunks, BlamePortion{
					CommitHash: currentCommit,
					NumLines:   n,
				})
			}
			continue
		}

		if strings.HasPrefix(line, targetField) {
			author := strings.TrimPrefix(line, targetField)
			if _, exists := commits[currentCommit]; !exists {
				commits[currentCommit] = CommitInfo{Author: author}
			}
		}
	}

	return FileStats{
		FilePath: filePath,
		Commits:  commits,
		Chunks:   chunks,
	}
}
