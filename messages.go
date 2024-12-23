package main

import (
	"crypto/rand"
	"embed"
	"fmt"
	"io/fs"
	"math/big"
	"strings"
)

type (
	Message struct {
		Path string
		Text string
	}

	MessagesStore struct {
		messages []Message
	}
)

//go:embed messages
var messagesFS embed.FS

func NewMessagesStore() (*MessagesStore, error) {
	messages, err := recursiveRead()
	if err != nil {
		return nil, fmt.Errorf("read embedded messages: %w", err)
	}
	return &MessagesStore{
		messages: messages,
	}, nil
}

func recursiveRead() ([]Message, error) {
	var messages []Message

	err := fs.WalkDir(messagesFS, "messages", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return fmt.Errorf("walk embedded path %q: %w", path, err)
		}

		if !d.IsDir() && strings.HasSuffix(d.Name(), ".md") {
			content, readErr := messagesFS.ReadFile(path)
			if readErr != nil {
				return fmt.Errorf("read embedded file %q: %w", path, readErr)
			}

			// Trim the "messages/" prefix to get the relative path
			relPath := strings.TrimPrefix(path, "messages/")

			messages = append(messages, Message{
				Path: relPath,
				Text: string(content),
			})
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return messages, nil
}

func (ms *MessagesStore) GetRandomMessage() (Message, error) {
	index, err := rand.Int(rand.Reader, big.NewInt(int64(len(ms.messages))))
	if err != nil {
		return Message{}, fmt.Errorf("generate random index: %w", err)
	}

	return ms.messages[index.Int64()], nil
}
