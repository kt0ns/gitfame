//go:build !solution

package main

import (
	"context"
	"fmt"
	"os"
	"runtime"
	"sync"

	"gitlab.com/slon/shad-go/gitfame/configs"
	"gitlab.com/slon/shad-go/gitfame/internal/config/config"
	"gitlab.com/slon/shad-go/gitfame/internal/config/languagebase"
	"gitlab.com/slon/shad-go/gitfame/internal/config/languagebase/JSON"
	"gitlab.com/slon/shad-go/gitfame/internal/errors"
	"gitlab.com/slon/shad-go/gitfame/internal/filter"
	"gitlab.com/slon/shad-go/gitfame/internal/format"
	"gitlab.com/slon/shad-go/gitfame/internal/git"
	"gitlab.com/slon/shad-go/gitfame/internal/statistics"
)

const (
	Workers = runtime.GOOS
)

func main() {
	loader := &JSON.JSONLoader{Data: configs.LanguageExtensionsJSON}
	langBase, err := languagebase.New(loader)
	if err != nil {
		fmt.Printf(errors.MsgLoadLanguageBase, err)
		os.Exit(1)
	}

	cfg := config.MustLoad(langBase)

	files, err := getTargetFiles(cfg)
	if err != nil {
		fmt.Printf(errors.MsgGetTargetFiles, err)
		os.Exit(1)
	}

	allFiles := len(files)

	jobs := make(chan string, allFiles)
	results := make(chan git.FileStats)

	for _, f := range files {
		jobs <- f
	}
	close(jobs)

	ctx, cancel := context.WithCancelCause(context.Background())

	processFile := func() {
		for {
			select {
			case <-ctx.Done():
				return
			case file, ok := <-jobs:
				if !ok {
					return
				}

				stats, err := git.BlameFile(cfg.Repository, cfg.Revision, file, cfg.UseCommitter)
				if err != nil {
					cancel(fmt.Errorf(errors.MsgWorkerBlame, file, err))
					return
				}

				select {
				case results <- stats:
				case <-ctx.Done():
					return
				}
			}
		}
	}

	wg := sync.WaitGroup{}

	for range Workers {
		wg.Add(1)
		go func() {
			defer wg.Done()
			processFile()
		}()
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	progressVisualizer := statistics.NewStdoutProgress()
	statsMap, err := statistics.AggregateStats(ctx, results, allFiles, progressVisualizer)
	if err != nil {
		fmt.Printf(errors.MsgAggregateBlame, err)
		cancel(nil)
		os.Exit(1)
	}

	orderToLessFunc := map[string]statistics.LessFunc{
		"lines":   statistics.LessByLines,
		"commits": statistics.LessByCommits,
		"files":   statistics.LessByFiles,
	}

	lessFunc, ok := orderToLessFunc[cfg.OrderBy]
	if !ok {
		fmt.Printf(errors.MsgUnknownOrderBy, cfg.OrderBy)
		cancel(nil)
		os.Exit(1)
	}

	sortedStats := statistics.SortStats(statsMap, lessFunc)

	if err := format.PrintStats(os.Stdout, cfg.Format, sortedStats); err != nil {
		fmt.Printf(errors.MsgPrintStats, err)
		cancel(nil)
		os.Exit(1)
	}
}

func getTargetFiles(cfg *config.Config) ([]string, error) {
	allFiles, err := git.LsTree(cfg.Repository, cfg.Revision)
	if err != nil {
		return nil, err
	}

	conditions := []filter.Condition{
		filter.ByExtensions(cfg.Extensions),
		filter.ExcludePatterns(cfg.Exclude),
		filter.RestrictToPatterns(cfg.RestrictTo),
	}

	return filter.ChainFilters(allFiles, conditions...), nil
}
