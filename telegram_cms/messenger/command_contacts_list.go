package messenger

import (
	"fmt"

	"github.com/amarnathcjd/gogram/telegram"
)

type ContactsListCmd struct {
	resultChan chan Result
}

func NewContactsListCmd() *ContactsListCmd {
	return &ContactsListCmd{
		resultChan: make(chan Result, 1),
	}
}

func (c *ContactsListCmd) Result() <-chan Result {
	return c.resultChan
}

func (c *ContactsListCmd) Execute(client *telegram.Client) error {
	contactsObj, err := client.ContactsGetContacts(0)
	if err != nil {
		c.resultChan <- Result{Success: false, Err: err}
		return err
	}

	var data = make(map[int64]*telegram.UserObj)

	if obj, ok := contactsObj.(*telegram.ContactsContactsObj); ok {
		for _, user := range obj.Users {
			if u, ok := user.(*telegram.UserObj); ok {
				data[u.ID] = u
			}
		}

		c.resultChan <- Result{Success: true, Data: data}
	} else {
		err = fmt.Errorf("неожиданный тип ответа: %T", contactsObj)
		c.resultChan <- Result{Success: false, Err: err}
		return err
	}

	return nil
}
