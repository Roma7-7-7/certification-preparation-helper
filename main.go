package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/aws/aws-lambda-go/lambda"
)

var httpClient = &http.Client{
	Timeout: 15 * time.Second,
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	slog.InfoContext(context.Background(), "lambda started")
	conf, err := GetConfig()
	if err != nil {
		slog.ErrorContext(ctx, "failed to get config", "error", err)
		os.Exit(1)
	}

	store, err := NewMessagesStore()
	if err != nil {
		slog.ErrorContext(ctx, "failed to create messages store", "error", err)
		os.Exit(1)
	}
	
	handler := NewLambdaHandler(
		store,
		NewTelegramClient(conf.TelegramToken, httpClient),
		conf.TelegramChatID,
	)
	if conf.IsDev {
		handler.HandleRequest(ctx)
		return
	} else {
		lambda.Start(handler.HandleRequest)
	}
}

func init() {
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})))
}
