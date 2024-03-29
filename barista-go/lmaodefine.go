package barista

import (
	"context"
	_ "embed"
	"encoding/json"
	"fmt"
	"math/rand"

	"github.com/appadeia/barista/barista-go/commandlib"
	"github.com/appadeia/barista/barista-go/config"
	"github.com/sashabaranov/go-openai"
)

//go:embed "linku.json"
var linkuString []byte

type linkuWord struct {
	Definition string
	Commentary string
}

var linkuDict = map[string]linkuWord{}

func init() {
	var s []struct {
		Word       string `json:"prompt"`
		Definition string `json:"completion"`
		Commentary string `json:"commentary"`
	}
	err := json.Unmarshal(linkuString, &s)
	if err != nil {
		panic(err)
	}
	for _, word := range s {
		linkuDict[word.Word] = linkuWord{word.Definition, word.Commentary}
	}
	linkuDict["eliki"] = linkuWord{
		Definition: "distress, affliction, adversity, pain",
		Commentary: "a nimisin made by kala suno pi kulupu telo mun with the intent to help people describe a range of experiences that may not fit in the space of alternative words—particularly those of ike or pakala—either because they are personally uncomfortable with it or feel it does not properly express the meaning they wish to convey.",
	}
	commandlib.RegisterCommand(commandlib.Command{
		Name:  "Define?",
		Usage: "Define? a toki pona word",
		ID:    "lmaodefine",
		Matches: []string{
			"ilo o pana sona e nimi",
			"ilo o define e nimi",
		},
		Flags: commandlib.FlagList{
			commandlib.StringFlag{
				LongFlag:  "prompt",
				ShortFlag: "p",
				FlagUsage: "which prompt to use",
				Value:     "",
			},
			commandlib.BoolFlag{
				LongFlag:  "musi",
				ShortFlag: "m",
				FlagUsage: "make the bot know less about toki pona (often generates funnier responses)",
				Value:     false,
			},
		},
		Action: define,
	})
}

type prompt struct {
	Contextless string
	Contextual  string
}

var prompts = map[string]prompt{
	"marketing": {
		Contextless: "using marketing speak, explain the following toki pona word:",
		Contextual:  "rephrase the definition of the following toki pona word using many cliches and buzzwords:",
	},
	"buzzword": {
		Contextless: "using many cliches and buzzwords, explain the following toki pona word:",
		Contextual:  "rephrase the definition of the following toki pona word using many cliches and buzzwords:",
	},
	"angry": {
		Contextless: "describe the following toki pona word in an impatient and irritable manner:",
		Contextual:  "rephrase the definition of the following toki pona word in an impatient and irritable manner:",
	},
	"shy": {
		Contextless: "explain the following toki pona word shyly:",
		Contextual:  "rephrase the definition of the following toki pona word in a shy manner:",
	},
	"billy": {
		Contextless: "pretend to be billy mays and explain the following toki pona word:",
		Contextual:  "rephrase the definition of the following toki pona word in the style of pretending to be billy mays:",
	},
	"poem": {
		Contextless: "using iambic pentameter, explain the following toki pona word:",
		Contextual:  "rephrase the definition of the following toki pona word in iambic pentameter:",
	},
	"jeopardy": {
		Contextless: "give a jeopardy prompt for the following toki pona word:",
		Contextual:  "rephrase the definition of the following toki pona word in the style of a jeopardy prompt: ",
	},
	"horse": {
		Contextless: "talk about the following toki pona word like you are commentating on a horse race:",
		Contextual:  "rephrase the definition of the following toki pona word in the style of commentary on a horse race:",
	},
	"flight": {
		Contextless: "explain the following toki pona word as if you were a flight attendant on a plane going down:",
		Contextual:  "pretend to be a flight attendant on a plane going down and explain the following toki pona word:",
	},
	"wwe": {
		Contextless: "give a wwe announcement for the following toki pona word (make sure to explain what their name means):",
		Contextual:  "give a wwe announcement for the following toki pona word (make sure to explain what their name means):",
	},
	"auction": {
		Contextless: "pretend to be an auctioneer and make a sales pitch for the following toki pona word:",
		Contextual:  "pretend to be an auctioneer and make a sales pitch for the following toki pona word:",
	},
	"walter": {
		Contextless: "Pretend you are Walter White from breaking bad, but instead of cooking methamphetamine you are telling Jesse how you need to cook the following toki pona word:",
		Contextual:  "Pretend you are Walter White from breaking bad, but instead of cooking methamphetamine you are telling Jesse how you need to cook the following toki pona word:",
	},
	"romance": {
		Contextless: "describe the following toki pona word as if you fell in love with it:",
		Contextual:  "describe the following toki pona word as if you fell in love with it:",
	},
}

func funnyPrompt(p string) prompt {
	if v, ok := prompts[p]; ok {
		return v
	}
	ps := []string{
		"marketing",
		"buzzword",
		"angry",
		"shy",
		"billy",
		"poem",
		"jeopardy",
	}

	return prompts[ps[rand.Intn(len(ps))]]
}

func define(c commandlib.Context) {
	word := c.Arg(0)
	if word == "" {
		c.SendMessage("response", commandlib.ErrorEmbed("sina pana ala e nimi a!"))
		return
	}
	c.CacheHint(5)
	client := openai.NewClient(config.BotConfig.Tokens.OpenAI)
	prompt := funnyPrompt(c.FlagValue("prompt"))
	var messages []openai.ChatCompletionMessage
	if def, ok := linkuDict[word]; ok && !c.IsFlagSet("musi") {
		messages = []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleUser,
				Content: fmt.Sprintf(`the toki pona dictionary says that "%s" means "%s"`, word, def.Definition),
			},
		}
		if def.Commentary != "" {
			messages = append(messages, openai.ChatCompletionMessage{
				Role:    openai.ChatMessageRoleUser,
				Content: fmt.Sprintf("it elaborates that %s", def.Commentary),
			})
		}
		messages = append(messages,
			openai.ChatCompletionMessage{
				Role:    openai.ChatMessageRoleUser,
				Content: "avoid repeating the dictionary verbatim",
			},
			openai.ChatCompletionMessage{
				Role:    openai.ChatMessageRoleUser,
				Content: prompt.Contextual,
			},
			openai.ChatCompletionMessage{
				Role:    openai.ChatMessageRoleUser,
				Content: word,
			},
		)
	} else {
		messages = []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleUser,
				Content: prompt.Contextless,
			},
			{
				Role:    openai.ChatMessageRoleUser,
				Content: word,
			},
		}
	}
	resp, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model:       openai.GPT3Dot5Turbo,
			Temperature: 1.0,
			Messages:    messages,
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
