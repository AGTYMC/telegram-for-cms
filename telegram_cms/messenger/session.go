package messenger

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/amarnathcjd/gogram/telegram"
)

type Session struct {
	Phone  string
	client *telegram.Client
	cmdCh  chan Command
	done   chan struct{}
}

func (s *Session) Start(ctx context.Context) error {
	// Подключаемся
	if _, err := s.client.Conn(); err != nil {
		return fmt.Errorf("connect: %w", err)
	}

	// Авторизация (интерактивная при необходимости)
	_, err := s.client.Login(s.Phone)
	if err != nil {
		if strings.Contains(err.Error(), "PHONE_NUMBER_INVALID") {
			return fmt.Errorf("номер телефона %s некорректен", s.Phone)
		}
		if strings.Contains(err.Error(), "FLOOD_WAIT") {
			return fmt.Errorf("flood wait при авторизации %s: %w", s.Phone, err)
		}
		return fmt.Errorf("login failed for %s: %w", s.Phone, err)
	}

	fmt.Printf("[messenger %s] Успешно авторизован\n", s.Phone)

	// Запускаем цикл обработки команд
	go s.worker()

	// Ожидаем внешней остановки
	<-ctx.Done()

	err = s.client.Disconnect()
	if err != nil {
		return err
	}

	close(s.done)
	time.Sleep(300 * time.Millisecond)
	return nil
}

func (s *Session) SendCommand(cmd Command) <-chan Result {
	// небольшая задержка оставлена для совместимости с твоим стилем,
	// но в большинстве случаев её можно убрать
	time.Sleep(400 * time.Millisecond)

	select {
	case s.cmdCh <- cmd:
		return cmd.Result()
	default:
		ch := make(chan Result, 1)
		ch <- Result{
			Success: false,
			Err:     fmt.Errorf("канал команд переполнен"),
		}
		close(ch)
		log.Printf("[messenger %s] канал команд переполнен — команда отброшена", s.Phone)
		return ch
	}
}

func (s *Session) worker() {
	for {
		select {
		case <-s.done:
			return
		case cmd := <-s.cmdCh:
			if err := cmd.Execute(s.client); err != nil {
				log.Printf("[messenger %s] ошибка выполнения команды: %v", s.Phone, err)
			}
		}
	}
}

func RunSessionInBackground(phone string, apiID int32, apiHash string, globalCtx context.Context) *Session {
	sess, err := NewSession(phone, apiID, apiHash)
	if err != nil {
		log.Fatalf("Не удалось создать сессию для %s: %v", phone, err)
	}

	go func() {
		if err := sess.Start(globalCtx); err != nil {
			log.Printf("Сессия %s завершилась с ошибкой: %v", phone, err)
		}
	}()

	return sess
}
