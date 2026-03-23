package statistics

import (
	"context"

	"gitlab.com/slon/shad-go/gitfame/internal/git"
)

type AuthorStats struct {
	Name    string
	Lines   int
	Commits map[string]struct{}
	Files   map[string]struct{}
}

func AggregateStats(ctx context.Context, results <-chan git.FileStats, totalFiles int, progress ProgressVisualizer) (map[string]*AuthorStats, error) {
	statsMap := make(map[string]*AuthorStats)
	processedFiles := 0

	for {
		select {
		case <-ctx.Done():
			return nil, context.Cause(ctx)
		case res, ok := <-results:
			if !ok {
				if progress != nil {
					progress.Finish()
				}
				return statsMap, nil
			}
			processedFiles++

			if progress != nil {
				progress.Update(processedFiles, totalFiles)
			}

			for _, chunk := range res.Chunks {
				commitInfo, ok := res.Commits[chunk.CommitHash]
				if !ok {
					continue
				}
				author := commitInfo.Author

				if statsMap[author] == nil {
					statsMap[author] = &AuthorStats{
						Name:    author,
						Commits: make(map[string]struct{}),
						Files:   make(map[string]struct{}),
					}
				}

				statsMap[author].Lines += chunk.NumLines
				statsMap[author].Commits[chunk.CommitHash] = struct{}{}
				statsMap[author].Files[res.FilePath] = struct{}{}
			}
		}
	}
}
