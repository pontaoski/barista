package telegram

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"

func (p *telegramPaginator) Send() {
	p.Index = 0
	send := p.Pages[p.Index]
	send.ReplyMarkup = p.context.keyboard()
	msg, err := p.bot.Send(send)
	if err == nil {
		p.message = &msg
		paginatorCache.Add(msg.MessageID, p)
	}
}

func (p *telegramPaginator) Prev() {
	p.Index--
	if p.Index < 0 {
		p.Index = len(p.Pages) - 1
	}
	send := p.Pages[p.Index]
	send.ReplyMarkup = p.context.keyboard()
	edit := tgbotapi.NewEditMessageText(p.message.Chat.ID, p.message.MessageID, "")
	edit.Text = p.Pages[p.Index].Text
	edit.ParseMode = p.Pages[p.Index].ParseMode
	kb := p.context.keyboard()
	edit.ReplyMarkup = &kb
	msg, err := p.bot.Send(edit)
	if err != nil {
		p.message = &msg
	}
}

func (p *telegramPaginator) Next() {
	p.Index++
	if p.Index+1 > len(p.Pages) {
		p.Index = 0
	}
	send := p.Pages[p.Index]
	send.ReplyMarkup = p.context.keyboard()
	edit := tgbotapi.NewEditMessageText(p.message.Chat.ID, p.message.MessageID, "")
	edit.Text = p.Pages[p.Index].Text
	edit.ParseMode = p.Pages[p.Index].ParseMode
	kb := p.context.keyboard()
	edit.ReplyMarkup = &kb
	msg, err := p.bot.Send(edit)
	if err != nil {
		p.message = &msg
	}
}
