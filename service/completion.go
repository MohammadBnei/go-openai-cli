package service

import (
	"context"
	"io"
	"strings"
	"time"

	"github.com/briandowns/spinner"
	"github.com/sashabaranov/go-openai"
	"github.com/thoas/go-funk"
)

func CreateChatCompletion(ctx context.Context) (string, error) {
	s := spinner.New(spinner.CharSets[26], 100*time.Millisecond)
	s.Start()
	defer s.Stop()

	resp, err := OpenAiService.client.CreateChatCompletionStream(
		ctx,
		openai.ChatCompletionRequest{
			Model:    OpenAiService.model,
			Messages: OpenAiService.Messages,
			Stream:   true,
		},
	)
	if err != nil {
		return "", err
	}
	defer resp.Close()

	fullMsg := ""
	role := ""

	for {
		msg, err := resp.Recv()
		s.Stop()
		if err == io.EOF {
			break
		}
		if err != nil {
			return "", err
		}

		OpenAiService.output.Write([]byte(msg.Choices[0].Delta.Content))
		fullMsg = strings.Join([]string{fullMsg, msg.Choices[0].Delta.Content}, "")
		if role == "" {
			role = msg.Choices[0].Delta.Role
		}
	}

	OpenAiService.AddMessage(openai.ChatCompletionMessage{
		Content: fullMsg,
		Role:    role,
	})

	OpenAiService.output.Write([]byte("\n"))

	return fullMsg, nil
}

func CreateCompletion(ctx context.Context) (string, error) {
	s := spinner.New(spinner.CharSets[26], 100*time.Millisecond)
	s.Start()
	defer s.Stop()

	resp, err := OpenAiService.client.CreateCompletionStream(
		ctx,
		openai.CompletionRequest{
			Model: OpenAiService.model,
			Prompt: funk.Map(OpenAiService.Messages, func(msg openai.ChatCompletionMessage) string {
				return msg.Content
			}).([]string),
			Stream: true,
		},
	)

	if err != nil {
		return "", err
	}
	defer resp.Close()

	fullMsg := ""
	role := ""

	for {
		msg, err := resp.Recv()
		s.Stop()
		if err == io.EOF {
			break
		}
		if err != nil {
			return "", err
		}

		OpenAiService.output.Write([]byte(msg.Choices[0].Text))
		fullMsg = strings.Join([]string{fullMsg, msg.Choices[0].Text}, "")
	}

	OpenAiService.AddMessage(openai.ChatCompletionMessage{
		Content: fullMsg,
		Role:    role,
	})

	OpenAiService.output.Write([]byte("\n"))

	return fullMsg, nil
}
