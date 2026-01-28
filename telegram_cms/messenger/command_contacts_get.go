package messenger

import (
	"fmt"

	"github.com/amarnathcjd/gogram/telegram"
)

type ContactsGetCmd struct {
	userId     int64
	resultChan chan Result
}

func NewContactsGetCmd(userId int64) *ContactsGetCmd {
	return &ContactsGetCmd{
		userId:     userId,
		resultChan: make(chan Result, 1),
	}
}

func (c *ContactsGetCmd) Result() <-chan Result {
	return c.resultChan
}

func (c *ContactsGetCmd) Execute(client *telegram.Client) error {
	user, err := client.GetUser(c.userId)
	if err != nil {
		c.resultChan <- Result{Success: false, Err: fmt.Errorf("resolve %q: %w", c.userId, err)}
		return err
	}
	c.resultChan <- Result{Success: true, Data: user}
	return nil
}
