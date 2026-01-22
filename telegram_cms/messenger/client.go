package messenger

import (
	"context"
)

type Client struct {
	Cancel  context.CancelFunc
	Session *Session
}

func (c Client) Close() {
	CloseTelegramClient(c.Cancel)
}

func (c Client) SendCommandAsync(cmd Command) {
	c.Session.SendCommand(cmd)
}

func (c Client) SendMessageAsync(target string, message string) {
	c.Session.SendCommand(NewSendMessageCmd(target, message))
}

func (c Client) ListContactsAsync() {
	c.Session.SendCommand(NewListContactsCmd())
}

func (c Client) SendMessage(target, message string) (bool, error) {
	ch := c.Session.SendCommand(NewSendMessageCmd(target, message))
	res := <-ch
	return res.Success, res.Err
}

func (c Client) ListContacts() ([]map[string]any, error) {
	ch := c.Session.SendCommand(NewListContactsCmd())
	res := <-ch
	if !res.Success {
		return nil, res.Err
	}
	if data, ok := res.Data.([]map[string]any); ok {
		return data, nil
	}
	return nil, nil
}

func (c Client) Execute(cmd Command) Result {
	ch := c.Session.SendCommand(cmd)
	return <-ch
}
