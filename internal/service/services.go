package service

import (
	"context"

	"github.com/dylanewe/moni/internal/store"
	"github.com/openai/openai-go/v3"
)

type Service struct {
	LLMParser interface {
		ParseStatement(ctx context.Context, categories []string, filepath string) ([]store.Transaction, error)
	}
}

func NewService(client *openai.Client) Service {
	return Service{
		LLMParser: &LLMParserService{client},
	}
}
