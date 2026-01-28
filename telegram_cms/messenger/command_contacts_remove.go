package messenger

import (
	"fmt"

	"github.com/amarnathcjd/gogram/telegram"
)

type ContactsRemoveCmd struct {
	userId     int64
	resultChan chan Result
}

func NewContactsRemoveCmd(userId int64) *ContactsRemoveCmd {
	return &ContactsRemoveCmd{
		userId:     userId,
		resultChan: make(chan Result, 1),
	}
}

func (c *ContactsRemoveCmd) Result() <-chan Result {
	return c.resultChan
}

func (c *ContactsRemoveCmd) Execute(client *telegram.Client) error {
	user, err := getUser(c.userId, client)
	if err != nil {
		c.resultChan <- Result{Success: false, Err: fmt.Errorf("GetContact() resolve %q: %w", c.userId, err)}
		return err
	}

	if user.AccessHash == 0 {
		c.resultChan <- Result{Success: false, Err: fmt.Errorf("у пользователя нет AccessHash — невозможно сформировать InputUser")}
		return err
	}

	input := &telegram.InputUserObj{
		UserID:     user.ID,
		AccessHash: user.AccessHash,
	}

	_, err = client.ContactsDeleteContacts([]telegram.InputUser{input})

	if err != nil {
		c.resultChan <- Result{Success: false, Err: fmt.Errorf("ContactsDeleteContacts() resolve %q: %w", c.userId, err)}
		return err
	}

	c.resultChan <- Result{Success: true, Err: nil}

	return nil
}
