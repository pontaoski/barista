package barista

import (
	"context"
	"math/rand"

	"github.com/appadeia/barista/barista-go/commandlib"
	"github.com/appadeia/barista/barista-go/config"
	"github.com/sashabaranov/go-openai"
)

func init() {
	commandlib.RegisterCommand(commandlib.Command{
		Name:  "Define?",
		Usage: "Define? a toki pona word",
		ID:    "lmaodefine",
		Matches: []string{
			"ilo o pana sona e nimi",
			"ilo o define e nimi",
		},
		Action: define,
	})
}

func funnyPrompt() string {
	prompts := []string{
		"using marketing speak, explain the following toki pona word:",
		"using many cliches and buzzwords, explain the following toki pona word:",
		"describe the following toki pona word in an impatient and irritable manner:",
		"explain the following toki pona word shyly:",
		"pretend to be billy mays and explain the following toki pona word:",
		"using iambic pentameter, explain the following toki pona word:",
		"give a jeopardy prompt for the following toki pona word:",
	}
	return prompts[rand.Intn(len(prompts))]
}

func define(c commandlib.Context) {
	word := c.Arg(0)
	if word == "" {
		c.SendMessage("response", commandlib.ErrorEmbed("sina pana ala e nimi a!"))
		return
	}
	client := openai.NewClient(config.BotConfig.Tokens.OpenAI)
	resp, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model:       openai.GPT3Dot5Turbo,
			Temperature: 1.0,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: funnyPrompt(),
				},
				{
					Role:    openai.ChatMessageRoleUser,
					Content: word,
				},
			},
		},
	)
	if err != nil {
		c.SendMessage("response", commandlib.ErrorEmbed(`ilo li pakala a! `+err.Error()))
		return
	}
	c.SendMessage("response", commandlib.Embed{
		Body:   resp.Choices[0].Message.Content,
		Colour: 0x3daee9,
	})
}
