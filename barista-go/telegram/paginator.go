package telegram

import (
	"sync"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type Paginator struct {
	Pages []tgbotapi.MessageConfig
	Index int

	message  *tgbotapi.Message
	bot      *tgbotapi.BotAPI
	lastused time.Time
}

var paginatorsMutex = &sync.Mutex{}
var paginators map[int]*Paginator = make(map[int]*Paginator)

func NewPaginator(bot *tgbotapi.BotAPI) Paginator {
	return Paginator{
		bot: bot,
	}
}

func paginatorHandler(messageID int, direction string) {
	if val, ok := paginators[messageID]; ok {
		if direction == "previous" {
			val.Prev()
		} else {
			val.Next()
		}
	}
}

func cleaner() {
	for {
		time.Sleep(5 * time.Minute)
		var rmkeys []int
		for key, cmd := range paginators {
			if time.Now().Sub(cmd.lastused) >= 10*time.Minute {
				rmkeys = append(rmkeys, key)
			}
		}
		for _, key := range rmkeys {
			paginatorsMutex.Lock()
			delete(paginators, key)
			paginatorsMutex.Unlock()
		}
	}
}

func (p *Paginator) AddPage(msg tgbotapi.MessageConfig) {
	p.Pages = append(p.Pages, msg)
}

var keyboard = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Previous", "previous"),
		tgbotapi.NewInlineKeyboardButtonData("Next", "next"),
	),
)

func (p *Paginator) Send() {
	p.Index = 0
	send := p.Pages[p.Index]
	send.ReplyMarkup = keyboard
	msg, err := p.bot.Send(send)
	if err == nil {
		p.message = &msg
		p.lastused = time.Now()
		paginators[msg.MessageID] = p
	}
}

func (p *Paginator) Prev() {
	p.Index--
	if p.Index < 0 {
		p.Index = len(p.Pages) - 1
	}
	send := p.Pages[p.Index]
	send.ReplyMarkup = keyboard
	edit := tgbotapi.NewEditMessageText(p.message.Chat.ID, p.message.MessageID, "")
	edit.Text = p.Pages[p.Index].Text
	edit.ParseMode = p.Pages[p.Index].ParseMode
	edit.ReplyMarkup = &keyboard
	msg, err := p.bot.Send(edit)
	if err != nil {
		p.message = &msg
		p.lastused = time.Now()
	}
}

func (p *Paginator) Next() {
	p.Index++
	if p.Index+1 > len(p.Pages) {
		p.Index = 0
	}
	send := p.Pages[p.Index]
	send.ReplyMarkup = keyboard
	edit := tgbotapi.NewEditMessageText(p.message.Chat.ID, p.message.MessageID, "")
	edit.Text = p.Pages[p.Index].Text
	edit.ParseMode = p.Pages[p.Index].ParseMode
	edit.ReplyMarkup = &keyboard
	msg, err := p.bot.Send(edit)
	if err != nil {
		p.message = &msg
		p.lastused = time.Now()
	}
}
