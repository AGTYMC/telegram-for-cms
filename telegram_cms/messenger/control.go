package messenger

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/AGTYMC/telegram-for-cms/telegram_cms/storage"

	"github.com/amarnathcjd/gogram/telegram"
)

func NewSession(phone string, apiID int32, apiHash string) (*Session, error) {
	sessionPath := storage.SessionFilePath(phone)

	if err := os.MkdirAll(filepath.Dir(sessionPath), 0700); err != nil {
		return nil, err
	}

	client, err := telegram.NewClient(telegram.ClientConfig{
		AppID:    apiID,
		AppHash:  apiHash,
		Session:  sessionPath,
		LogLevel: telegram.ErrorLevel,
	})
	if err != nil {
		return nil, fmt.Errorf("create client: %w", err)
	}

	sess := &Session{
		Phone:  phone,
		client: client,
		cmdCh:  make(chan Command, 32),
		done:   make(chan struct{}),
	}

	return sess, nil
}

func CreateTelegramClient(phone string, apiID int32, apiHash string) *Client {
	ctx, cancel := context.WithCancel(context.Background())
	sess, err := RunSessionInBackground(phone, apiID, apiHash, ctx)
	if err != nil {
		fmt.Println(err)
		cancel()
	}
	return &Client{Cancel: cancel, Session: sess}
}

func CloseTelegramClient(cancel context.CancelFunc) {
	cancel()
	time.Sleep(2 * time.Second)
}
