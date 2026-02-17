package AI

import (
	"context"
	"errors"
	"github.com/saleh-ghazimoradi/Projectopher/config"
	"github.com/tmc/langchaingo/llms/openai"
	"strings"
)

type OpenAI interface {
	GetSentiment(ctx context.Context, review string, sentiments []string) (string, error)
}

type openAI struct {
	config *config.Config
	llm    *openai.LLM
}

func (o *openAI) GetSentiment(ctx context.Context, review string, sentiments []string) (string, error) {
	if len(sentiments) == 0 {
		return "", errors.New("no sentiments provided")
	}

	delimited := strings.Join(sentiments, ",")

	basePrompt := strings.Replace(o.config.OpenAI.BasePromptTemplate, "{rankings}", delimited, 1)

	fullPrompt := basePrompt + review

	response, err := o.llm.Call(ctx, fullPrompt)
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(response), nil
}

func NewOpenAI(config *config.Config, llm *openai.LLM) OpenAI {
	return &openAI{
		config: config,
		llm:    llm,
	}
}
