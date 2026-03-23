package statistics

import (
	"sort"
)

type SortableAuthorStats []AuthorStats

type LessFunc func(p1, p2 *AuthorStats) bool

func SortStats(statsMap map[string]*AuthorStats, less LessFunc) []AuthorStats {
	if len(statsMap) == 0 {
		return []AuthorStats{}
	}

	result := make([]AuthorStats, 0, len(statsMap))
	for _, stats := range statsMap {
		result = append(result, *stats)
	}

	sort.Slice(result, func(i, j int) bool {
		return less(&result[i], &result[j])
	})

	return result
}

func LessByLines(p1, p2 *AuthorStats) bool {
	lines1, lines2 := p1.Lines, p2.Lines
	commits1, commits2 := len(p1.Commits), len(p2.Commits)
	files1, files2 := len(p1.Files), len(p2.Files)

	if lines1 != lines2 {
		return lines1 > lines2
	}
	if commits1 != commits2 {
		return commits1 > commits2
	}
	if files1 != files2 {
		return files1 > files2
	}
	return p1.Name < p2.Name
}

func LessByCommits(p1, p2 *AuthorStats) bool {
	lines1, lines2 := p1.Lines, p2.Lines
	commits1, commits2 := len(p1.Commits), len(p2.Commits)
	files1, files2 := len(p1.Files), len(p2.Files)

	if commits1 != commits2 {
		return commits1 > commits2
	}
	if lines1 != lines2 {
		return lines1 > lines2
	}
	if files1 != files2 {
		return files1 > files2
	}
	return p1.Name < p2.Name
}

func LessByFiles(p1, p2 *AuthorStats) bool {
	lines1, lines2 := p1.Lines, p2.Lines
	commits1, commits2 := len(p1.Commits), len(p2.Commits)
	files1, files2 := len(p1.Files), len(p2.Files)

	if files1 != files2 {
		return files1 > files2
	}
	if lines1 != lines2 {
		return lines1 > lines2
	}
	if commits1 != commits2 {
		return commits1 > commits2
	}
	return p1.Name < p2.Name
}
