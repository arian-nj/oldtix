package telegrambot

import (
	"context"
	"strconv"

	"github.com/arian-nj/master-card/back/internal/babel"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func (app *ApplicationTelegram) botRoutes(b *bot.Bot) {
	b.RegisterHandler(bot.HandlerTypeMessageText, "/start", bot.MatchTypeExact, app.StartHandler)
}

func (app *ApplicationTelegram) StartHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	go func() {
		if update.Message != nil {
			_, err := app.Queries.InsertBotUser(context.Background(), strconv.Itoa(int(update.Message.From.ID)))
			if err != nil {
				app.Logger.Error(err.Error())
				return
			}
		}
	}()
	_, err := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   babel.Persian.Welcome,
		ReplyMarkup: models.InlineKeyboardMarkup{
			InlineKeyboard: [][]models.InlineKeyboardButton{
				{
					{Text: "بازی", WebApp: &models.WebAppInfo{URL: "https://tix.filelord.ir/telegram/"}},
				},
			},
		},
	})
	if err != nil {
		app.Logger.Error(err.Error())
	}
}

func (app *ApplicationTelegram) defaultHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	_, err := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   babel.Persian.Default,
	})
	if err != nil {
		app.Logger.Error(err.Error())
	}
}
