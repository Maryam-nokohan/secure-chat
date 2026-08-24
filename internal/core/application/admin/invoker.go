package admin

import (
	"context"
	"errors"
	"sync"
)

type HistoryEntry struct {
	Description string `json:"description"`
}

type Invoker struct {
	mu      sync.Mutex
	history []Command
	max     int
}

func NewInvoker(max int) *Invoker {
	if max <= 0 {
		max = 50
	}
	return &Invoker{max: max}
}

func (i *Invoker) Do(ctx context.Context, cmd Command) error {
	if err := cmd.Execute(ctx); err != nil {
		return err
	}
	i.mu.Lock()
	i.history = append(i.history, cmd)
	if len(i.history) > i.max {
		i.history = i.history[len(i.history)-i.max:]
	}
	i.mu.Unlock()
	return nil
}

func (i *Invoker) UndoLast(ctx context.Context) (string, error) {
	i.mu.Lock()
	if len(i.history) == 0 {
		i.mu.Unlock()
		return "", errors.New("nothing to undo")
	}
	cmd := i.history[len(i.history)-1]
	i.history = i.history[:len(i.history)-1]
	i.mu.Unlock()

	if err := cmd.Undo(ctx); err != nil {
		return "", err
	}
	return cmd.Description(), nil
}

func (i *Invoker) History() []HistoryEntry {
	i.mu.Lock()
	defer i.mu.Unlock()
	out := make([]HistoryEntry, 0, len(i.history))
	for _, c := range i.history {
		out = append(out, HistoryEntry{Description: c.Description()})
	}
	return out
}