package messenger

import (
	"fmt"

	"github.com/amarnathcjd/gogram/telegram"
)

type SendMessageCmd struct {
	Target     string
	Text       string
	resultChan chan Result
}

func NewSendMessageCmd(target, text string) *SendMessageCmd {
	return &SendMessageCmd{
		Target:     target,
		Text:       text,
		resultChan: make(chan Result, 1),
	}
}

func (c *SendMessageCmd) Result() <-chan Result {
	return c.resultChan
}

func (c *SendMessageCmd) Execute(client *telegram.Client) error {
	peer, err := client.ResolvePeer(c.Target)
	if err != nil {
		c.resultChan <- Result{Success: false, Err: fmt.Errorf("resolve %q: %w", c.Target, err)}
		return err
	}

	_, err = client.SendMessage(peer, c.Text)
	if err != nil {
		c.resultChan <- Result{Success: false, Err: fmt.Errorf("send: %w", err)}
		return err
	}

	c.resultChan <- Result{Success: true}
	return nil
}
