package agent

import (
	"context"

	"github.com/anthropics/anthropic-sdk-go"
)

func NewAgent(client *anthropic.Client) *Agent {
	return &Agent{
		client: client,
	}
}

type Agent struct {
	client *anthropic.Client
}

func (a *Agent) Explore(ctx context.Context, path string) error {
	return nil
}
