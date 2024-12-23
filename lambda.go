package main

import (
	"context"
	"fmt"
	"log/slog"
)

type LambdaHandler struct {
	messagesStore  *MessagesStore
	telegramClient *TelegramClient
	telegramChatID string
}

func NewLambdaHandler(messagesStore *MessagesStore, telegramClient *TelegramClient, telegramChatID string) *LambdaHandler {
	return &LambdaHandler{
		messagesStore:  messagesStore,
		telegramClient: telegramClient,
		telegramChatID: telegramChatID,
	}
}

func (h *LambdaHandler) HandleRequest(ctx context.Context) {
	slog.InfoContext(ctx, "handle request")

	msg, err := h.messagesStore.GetRandomMessage()
	if err != nil {
		slog.ErrorContext(ctx, "failed to get random message", "error", err)
	}

	if err = h.telegramClient.SendMessage(ctx, h.telegramChatID, fmt.Sprintf("===============\n%s", msg.Text)); err != nil {
		slog.ErrorContext(ctx, "failed to send message", "path", msg.Path, "error", err)
		_ = h.telegramClient.SendMessage(ctx, h.telegramChatID, fmt.Sprintf("Failed to send message for path %1: %s", msg.Path, err))
	}

	slog.InfoContext(ctx, "request handled")
}
