package pkg

import "sync"

type ringLog struct {
	mu    sync.Mutex
	lines []string
	max   int
	pos   int
	full  bool
}

var recentLogs = &ringLog{max: 1000, lines: make([]string, 1000)}

func (r *ringLog) Write(p []byte) (int, error) {
	r.mu.Lock()
	r.lines[r.pos] = string(p)
	r.pos = (r.pos + 1) % r.max
	if r.pos == 0 {
		r.full = true
	}
	r.mu.Unlock()
	return len(p), nil
}

func GetRecentLogs(limit int) []string {
	recentLogs.mu.Lock()
	defer recentLogs.mu.Unlock()

	var ordered []string
	if recentLogs.full {
		ordered = append(ordered, recentLogs.lines[recentLogs.pos:]...)
	}
	ordered = append(ordered, recentLogs.lines[:recentLogs.pos]...)

	if limit > 0 && limit < len(ordered) {
		ordered = ordered[len(ordered)-limit:]
	}
	return ordered
}