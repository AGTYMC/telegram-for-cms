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

func (s *Session) Start(ctx context.Context, result chan<- Result) Result {
	// Подключаемся
	if _, err := s.client.Conn(); err != nil {
		return Result{Success: false, Err: err}
	}

	// Авторизация (интерактивная при необходимости)
	_, err := s.client.Login(s.Phone)
	if err != nil {
		if strings.Contains(err.Error(), "PHONE_NUMBER_INVALID") {
			return Result{Success: false, Err: err, Message: "PHONE_NUMBER_INVALID"}
		}
		if strings.Contains(err.Error(), "FLOOD_WAIT") {
			return Result{Success: false, Err: err, Message: "FLOOD_WAIT"}
		}
		return Result{Success: false, Err: err, Message: "ERROR"}
	}

	// Успешная авторизация
	result <- Result{Success: true, Message: fmt.Sprintf("[messenger %s] Успешно авторизован\n", s.Phone)}

	// Цикл обработки команд
	go s.worker()

	// Ожидаем внешней остановки
	<-ctx.Done()

	// Завершаем соединение
	if err = s.client.Disconnect(); err != nil {
		return Result{Success: false, Err: err, Message: "DISCONNECT"}
	}

	close(s.done)
	close(result)

	time.Sleep(300 * time.Millisecond)
	return Result{Success: true, Err: nil}
}

func (s *Session) SendCommand(cmd Command) <-chan Result {
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

func RunSessionInBackground(phone string, apiID int32, apiHash string, globalCtx context.Context) (*Session, error) {
	sess, err := NewSession(phone, apiID, apiHash)
	if err != nil {
		return nil, fmt.Errorf("не удалось создать сессию для %s: %w", phone, err)
	}

	//Канал с результатом
	result := make(chan Result, 1)

	//Запускаем в отдельной горутине, где получаем результат авторизации
	go func() {
		if res := sess.Start(globalCtx, result); res.Success {
			return
		}
	}()

	//Ожидаем результат авторизации
	if res := <-result; !res.Success {
		return nil, res.Err
	}

	return sess, nil
}
