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

func RunSessionInBackground(phone string, apiID int32, apiHash string, globalCtx context.Context) (*Session, error) {
	sess, err := NewSession(phone, apiID, apiHash)
	if err != nil {
		log.Fatalf("Не удалось создать сессию для %s: %v", phone, err)
	}

	done := make(chan error, 1)

	go func() {
		defer close(done)
		if err := sess.Start(globalCtx); err != nil {
			done <- fmt.Errorf("сессия %s завершилась с ошибкой: %w", phone, err)
			return
		}
		// Если Start вышел без ошибки → сессия закрыта нормально (по ctx.Done())
		done <- nil
	}()

	// Теперь ждём авторизации (или завершения/ошибки горутины)
	const maxWait = 6 * time.Minute // обычно хватает даже на ввод кода + 2FA
	const checkInterval = 400 * time.Millisecond

	ticker := time.NewTicker(checkInterval)
	defer ticker.Stop()

	deadline := time.Now().Add(maxWait)

	time.Sleep(5 * time.Second)

	for {
		select {
		case err := <-done:
			if err != nil {
				return nil, err
			}
			return nil, fmt.Errorf("сессия %s неожиданно завершилась без авторизации", phone)

		case <-ticker.C:

			authorized, aErr := sess.client.IsAuthorized()
			if aErr != nil {
				return nil, fmt.Errorf("ошибка проверки авторизации %s: %w", phone, aErr)
			}
			if authorized {
				log.Printf("[messenger %s] Успешно авторизован (в RunSessionInBackground)", phone)
				return sess, nil // ← возвращаем только теперь
			}

			if time.Now().After(deadline) {
				return nil, fmt.Errorf("таймаут ожидания авторизации для %s (> %v)", phone, maxWait)
			}

		case <-globalCtx.Done():
			return nil, fmt.Errorf("контекст отменён во время ожидания авторизации %s: %w", phone, globalCtx.Err())
		}
	}
}
