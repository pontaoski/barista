package barista

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/appadeia/barista/barista-go/commandlib"
)

func init() {
	commandlib.RegisterCommand(commandlib.Command{
		Name:    "toki pona vocabulary lookup",
		Usage:   "Look up words in pu",
		ID:      "pu-lookup",
		Matches: []string{"ilo o pu"},
		Flags: []commandlib.Flag{
			commandlib.BoolFlag{
				LongFlag: "browse",
			},
		},
		Action: PuSearch,
	})
	commandlib.RegisterCommand(commandlib.Command{
		Name:    "toki pona quiz",
		Usage:   "Quiz yourself on toki pona (pu)",
		ID:      "tpo-quiz",
		Matches: []string{"o quiz pu"},
		Action:  QuizPu,
	})
	commandlib.RegisterCommand(commandlib.Command{
		Name:    "toki pona quiz",
		Usage:   "Quiz yourself on toki pona (nimi ale pona)",
		ID:      "tpo-quiz",
		Matches: []string{"o quiz nap"},
		Action:  QuizNap,
	})
}

func lookupPu(word string) (Word, bool) {
	for _, val := range words {
		for _, name := range val.Names {
			if name == word {
				return val, true
			}
		}
	}
	return Word{}, false
}

func toEmbed(word Word) commandlib.Embed {
	embed := commandlib.Embed{}
	embed.Title.Text = strings.Join(word.Names, ", ")
	embed.Fields = []commandlib.EmbedField{
		{
			Title: "Meaning",
			Body:  word.Definition,
		},
		{
			Title:  "Kind",
			Body:   word.Category,
			Inline: true,
		},
		{
			Title:  "Source Language",
			Body:   word.SourceLanguage,
			Inline: true,
		},
	}
	if word.Etymology != "" {
		embed.Fields = append(embed.Fields, commandlib.EmbedField{
			Title:  "Etymology",
			Body:   word.Etymology,
			Inline: true,
		})
	}
	if word.Tag != "" {
		embed.Fields = append(embed.Fields, commandlib.EmbedField{
			Title:  "Notable Thing",
			Body:   word.Tag,
			Inline: true,
		})
	}

	return embed
}

func PuSearch(c commandlib.Context) {
	if c.IsFlagSet("browse") {
		embedlist := commandlib.EmbedList{
			ItemTypeName: c.I18n("Word"),
		}
		for _, word := range words {
			embedlist.Embeds = append(embedlist.Embeds, toEmbed(word))
		}
		c.SendMessage("primary", embedlist)
		return
	}
	if c.NArgs() < 1 {
		c.SendMessage(
			"primary",
			commandlib.ErrorEmbed(c.I18n("Please provide a word from pu to look up the meaning of.")),
		)
		return
	}
	if word, ok := lookupPu(c.Arg(0)); ok {
		c.SendMessage("primary", toEmbed(word))
		return
	}

	c.SendMessage(
		"primary",
		commandlib.ErrorEmbed(c.I18n("Sorry, I don't think that's a word.")),
	)
	return
}

const nextRespDuration = 3 * time.Second

type meaningWordList []string

func (m meaningWordList) contains(s string) bool {
	lower := strings.ToLower(s)
	for _, word := range strings.Fields(lower) {
		for _, meaning := range m {
			if strings.TrimSpace(word) == strings.TrimSpace(meaning) {
				return true
			}
		}
	}
	return false
}

func (p Word) toMeaningWordList() (ret meaningWordList) {
	for _, split := range strings.Fields(strings.ReplaceAll(strings.ReplaceAll(p.Definition, ",", " "), "!", " ")) {
		ret = append(ret, strings.ToLower(split))
	}
	return
}

func Quiz(c commandlib.Context, words []Word) {
	score := 0
	rand.Seed(time.Now().Unix())
	c.SendMessage("starting", fmt.Sprintf(c.I18n("Starting quiz... Type '%s' to cancel"), c.I18n("cancel")))

	i := -1
Quizzer:
	for _, word := range words {
		i++
		c.SendMessage(fmt.Sprintf("primary-%d", i), commandlib.Embed{
			Title: commandlib.EmbedHeader{
				Text: strings.Join(word.Names, ", "),
			},
			Body: c.I18n("What does this word mean?"),
			Footer: commandlib.EmbedHeader{
				Text: fmt.Sprintf("Word %d out of %d", i+1, 10),
			},
		})
		timeoutChan := make(chan struct{})
		go func() {
			time.Sleep(7 * time.Second)
			timeoutChan <- struct{}{}
		}()
		for {
			select {
			case msg := <-c.NextResponse():
				if strings.Contains(msg, c.I18n("cancel")) {
					c.SendMessage(fmt.Sprintf("primary-%d", i), commandlib.ErrorEmbed(c.I18n("Quiz cancelled.")))
					return
				}
				if word.toMeaningWordList().contains(msg) {
					grats := toEmbed(word)
					grats.Colour = 0x00ff00
					c.SendMessage(fmt.Sprintf("primary-%d", i), grats)
					c.SendMessage(fmt.Sprintf("congrats-%d", i), "Correct!")
					score++
					time.Sleep(nextRespDuration)
					continue Quizzer
				}
			case <-timeoutChan:
				wrong := toEmbed(word)
				wrong.Colour = 0xff0000
				c.SendMessage(fmt.Sprintf("primary-%d", i), wrong)
				time.Sleep(nextRespDuration)
				continue Quizzer
			}
		}
	}

	c.SendMessage("primary", commandlib.Embed{
		Title: commandlib.EmbedHeader{
			Text: c.I18n("Quiz Results"),
		},
		Fields: []commandlib.EmbedField{
			{
				Title: c.I18n("Correct Results"),
				Body:  strconv.Itoa(score),
			},
			{
				Title: c.I18n("Incorrect Results"),
				Body:  strconv.Itoa(10 - score),
			},
		},
	})
}

func QuizPu(c commandlib.Context) {
	rand.Seed(time.Now().Unix())

	slice := make([]Word, 10)
	for i := 0; i < 10; i++ {
		word := Word{}
		for word.Category != "pu" {
			word = words[rand.Intn(len(words))]
		}
		slice[i] = word
	}
	Quiz(c, slice)
}

func QuizNap(c commandlib.Context) {
	rand.Seed(time.Now().Unix())

	slice := make([]Word, 10)
	for i := 0; i < 10; i++ {
		slice[i] = words[rand.Intn(len(words))]
	}
	Quiz(c, slice)
}
