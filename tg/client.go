package tg

import (
	"fmt"
	"log/slog"
	"strconv"

	"time"

	tele "gopkg.in/telebot.v3"
)

type Client struct {
	bot      *tele.Bot
	chatID   tele.ChatID
	chat     *tele.Chat
	handlers []BotUpdateHandler
}

func NewClient(token string, chatID tele.ChatID, pollerTimeoutSeconds int) (*Client, error) {
	bot, err := tele.NewBot(tele.Settings{
		Token:  token,
		Poller: &tele.LongPoller{Timeout: time.Duration(pollerTimeoutSeconds) * time.Second},
	})
	if err != nil {
		return nil, err
	}

	return &Client{
		bot:    bot,
		chatID: chatID,
		chat: &tele.Chat{
			ID: int64(chatID),
		},
	}, nil
}

type BotUpdateHandler interface {
	Match(c tele.Context) bool
	Handle(c tele.Context) error
}

func (c *Client) RegisterHandler(handler BotUpdateHandler) {
	c.handlers = append(c.handlers, handler)
}
func (c *Client) Start() {

	handler := func(ctx tele.Context) error {
		for _, handler := range c.handlers {

			if !handler.Match(ctx) {
				continue
			}

			err := handler.Handle(ctx)
			if err != nil {
				slog.Error("handle message", slog.Any("error", err))
			}

		}
		return nil
	}
	c.bot.Handle(tele.OnText, handler)
	go c.bot.Start()
}

func (c *Client) Stop() {
	c.bot.Stop()
}

func (c *Client) SendSpoilerLink(threadID int, header, link string) (int, error) {
	payload := fmt.Sprintf("%s\n%s", header, link)

	message, err := c.bot.Send(c.chatID, payload, &tele.SendOptions{
		ThreadID:              threadID,
		DisableWebPagePreview: true,
		Entities: []tele.MessageEntity{
			{
				Type:   tele.EntitySpoiler,
				Offset: len(header) + 1,
				Length: len(link),
			},
		},
	})
	if err != nil {
		return 0, err
	}

	return message.ID, nil
}

func (c *Client) SendSticker(threadID int, stickerID string) (int, error) {
	sticker := tele.Sticker{
		File: tele.File{
			FileID: stickerID,
		},
	}
	message, err := c.bot.Send(c.chatID, &sticker, &tele.SendOptions{
		ThreadID: threadID,
	})
	if err != nil {
		return 0, fmt.Errorf("send: %w", err)
	}

	return message.ID, nil
}

func (c *Client) ReplyWithSticker(stickerID string, msg *tele.Message) (int, error) {
	sticker := tele.Sticker{
		File: tele.File{
			FileID: stickerID,
		},
	}
	message, err := c.bot.Send(c.chatID, &sticker, &tele.SendOptions{
		ReplyTo: msg,
	})
	if err != nil {
		return 0, fmt.Errorf("send: %w", err)
	}

	return message.ID, nil
}

func (c *Client) Pin(id int) error {
	msg := tele.StoredMessage{
		MessageID: strconv.Itoa(id),
		ChatID:    c.chat.ID,
	}

	return c.bot.Pin(msg, tele.Silent)
}

func (c *Client) Unpin(id int) error {
	return c.bot.Unpin(c.chat, id)
}
