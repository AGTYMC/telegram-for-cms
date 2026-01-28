package messenger

import (
	"context"
	"fmt"

	"github.com/amarnathcjd/gogram/telegram"
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
	c.Session.SendCommand(NewContactsListCmd())
}

func (c Client) SendMessage(target, message string) (bool, error) {
	ch := c.Session.SendCommand(NewSendMessageCmd(target, message))
	res := <-ch
	return res.Success, res.Err
}

func (c Client) ListContacts() (map[int64]*telegram.UserObj, error) {
	ch := c.Session.SendCommand(NewContactsListCmd())
	res := <-ch
	if !res.Success {
		return nil, res.Err
	}
	if data, ok := res.Data.(map[int64]*telegram.UserObj); ok {
		return data, nil
	}
	return nil, nil
}

func (c Client) CheckContacts(contact string) (*telegram.UserObj, error) {
	ch := c.Session.SendCommand(NewContactsCheckCmd(contact))
	res := <-ch
	if !res.Success && res.Err != nil {
		return nil, res.Err
	}
	if data, ok := res.Data.(*telegram.UserObj); ok {
		return data, nil
	}
	return nil, fmt.Errorf("check contacts is failed")
}

func (c Client) AddContact(phone, username, firstName, lastName string) (*telegram.UserObj, error) {
	ch := c.Session.SendCommand(NewContactsAddCmd(phone, username, firstName, lastName))
	res := <-ch
	if !res.Success && res.Err != nil {
		return nil, res.Err
	}
	if data, ok := res.Data.(*telegram.UserObj); ok {
		return data, nil
	}
	return nil, fmt.Errorf("imported contact list is empty")
}

func (c Client) GetContact(userId int64) (*telegram.UserObj, error) {
	ch := c.Session.SendCommand(NewContactsGetCmd(userId))
	res := <-ch
	if !res.Success && res.Err != nil {
		return nil, res.Err
	}
	if data, ok := res.Data.(*telegram.UserObj); ok {
		return data, nil
	}
	return nil, fmt.Errorf("get user is failed")
}

func (c Client) RemoveContact(userId int64) (bool, error) {
	ch := c.Session.SendCommand(NewContactsRemoveCmd(userId))
	res := <-ch
	if !res.Success && res.Err != nil {
		return false, res.Err
	}
	if res.Success {
		return true, nil
	}
	return false, fmt.Errorf("remove user is failed")
}

func (c Client) Execute(cmd Command) Result {
	ch := c.Session.SendCommand(cmd)
	return <-ch
}
