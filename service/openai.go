package service

import (
	"context"
	"io"
	"os"

	"github.com/sashabaranov/go-openai"
	"github.com/spf13/viper"
)

var OpenAiService *openAiService

type openAiService struct {
	client             *openai.Client
	Messages           []openai.ChatCompletionMessage
	chatCompletionList map[string]any
	model              string
	output             io.Writer
}

func InitOpenAiService() {
	chatCompletionList := map[string]any{"code-davinci-002": true, "text-davinci-002": true}
	model := viper.GetString("model")
	if model == "" {
		model = openai.GPT3Dot5Turbo
	}

	OpenAiService = &openAiService{
		client:             openai.NewClient(viper.GetString("OPENAI_KEY")),
		Messages:           []openai.ChatCompletionMessage{},
		chatCompletionList: chatCompletionList,
		model:              model,
		output:             os.Stdout,
	}
}

func (s *openAiService) SendPrompt(ctx context.Context, text string, output io.Writer) (string, error) {
	s.AddMessage(openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleUser,
		Content: text,
	})
	model := viper.GetString("model")

	switch {
	case s.chatCompletionList[model] == nil:
		return CreateChatCompletion(ctx)
	default:
		return CreateCompletion(ctx)
	}

}

func (s *openAiService) AddMessage(msg openai.ChatCompletionMessage) {
	s.Messages = append(s.Messages, msg)

	if len(s.Messages) > viper.GetInt("messages-length") {
		s.Messages = s.Messages[1:]
	}
}

func (s *openAiService) ClearMessages() {
	s.Messages = []openai.ChatCompletionMessage{}
}

func (s *openAiService) GetMessages() []openai.ChatCompletionMessage {
	return s.Messages
}

func (s *openAiService) GetModelList() ([]string, error) {
	models, err := s.client.ListModels(context.Background())
	if err != nil {
		return nil, err
	}

	modelsList := []string{}
	for _, model := range models.Models {
		modelsList = append(modelsList, model.ID)
	}

	return modelsList, nil
}
