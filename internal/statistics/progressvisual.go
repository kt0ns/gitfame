package statistics

import (
	"fmt"
	"os"
)

type ProgressVisualizer interface {
	Update(current, total int)
	Finish()
}

type stdoutProgress struct{}

func NewStdoutProgress() ProgressVisualizer {
	return &stdoutProgress{}
}

func (p *stdoutProgress) Update(current, total int) {
	if total == 0 {
		return
	}
	_, _ = fmt.Fprintf(os.Stderr, "Processing files: [%d/%d] %.1f%%",
		current,
		total,
		float64(current)/float64(total)*100,
	)
}

func (p *stdoutProgress) Finish() {
	// fmt.Println()
}
