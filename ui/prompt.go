package ui

import (
	"context"
	"fmt"
	"os"

	"github.com/MohammadBnei/go-openai-cli/service"
	"github.com/ktr0731/go-fuzzyfinder"
	"github.com/manifoldco/promptui"
	"github.com/sashabaranov/go-openai"
	"github.com/thoas/go-funk"
)

func OpenAiPrompt() {

	var label string
	help := `
		q: quit
		h: help
		s: save the response to a file
		f: add files to the messages (won't send to openAi until you send a prompt)
		c: clear messages and files
		
		any other text will be sent to openAI
		`

	fmt.Println("for help type 'h'")

	previousRes := ""
	previousPrompt := ""

	// lastImagePath := ""

	fileNumber := 0
PromptLoop:
	for {
		label = "prompt"
		totalCharacters := funk.Reduce(service.GetMessages(), func(acc int, elem openai.ChatCompletionMessage) int {
			return acc + len(elem.Content)
		}, 0)
		if totalCharacters != 0 {
			label = fmt.Sprintf("%d🔤 follow up", totalCharacters)
		}
		if fileNumber != 0 {
			label = fmt.Sprintf("%d💾 %s ", fileNumber, label)
		}

		prompt := promptui.Prompt{
			Label:     label,
			AllowEdit: false,
			Default:   previousPrompt,
		}

		userPrompt, err := prompt.Run()
		if err != nil {
			fmt.Println(err)
			return
		}

		switch userPrompt {
		case "q":
			break PromptLoop
		case "h":
			fmt.Println(help)

		case "s":
			SaveToFile([]byte(previousRes))

		case "e":
			messages := service.GetMessages()
			messages = funk.Filter(messages, func(msg openai.ChatCompletionMessage) bool {
				return msg.Role != "user"
			}).([]openai.ChatCompletionMessage)

			idx, err := fuzzyfinder.Find(
				messages,
				func(i int) string {
					return fmt.Sprintf("%.35s...", messages[i].Content)
				},
				fuzzyfinder.WithPreviewWindow(func(i, w, h int) string {
					return messages[i].Content
				}))

			if err != nil {
				fmt.Println(err)
				continue PromptLoop
			}
			instructionPrompt := promptui.Prompt(promptui.Prompt{
				Label:     "clarify your edit instruction",
				IsVimMode: true,
			})

			instruction, err := instructionPrompt.Run()
			if err != nil {
				fmt.Println(err)
				continue PromptLoop
			}

			response, err := service.SendEditPrompt(context.Background(), instruction, idx, os.Stdout)
			if err != nil {
				fmt.Println(err)
				return
			}

			service.GetMessages()[idx].Content = response

		case "i":
			// lastImagePath = AskForImage()
			AskForImage()

		// case "ei":
		// 	lastImagePath = AskForEditImage(lastImagePath)

		case "c":
			service.ClearMessages()
			fileNumber = 0
			fmt.Println("cleared messages")

		case "f":
			FileSelectionFzf(&fileNumber)

		default:
			response, err := service.SendPrompt(context.Background(), userPrompt, os.Stdout)
			if err != nil {
				fmt.Println(err)
				return
			}
			previousRes = response
			fileNumber = 0
		}

		previousPrompt = userPrompt
	}
}
