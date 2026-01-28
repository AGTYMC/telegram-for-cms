package messenger

import (
	"fmt"
	"strings"

	"github.com/amarnathcjd/gogram/telegram"
)

type ContactsCheckCmd struct {
	contact    string
	resultChan chan Result
}

func NewContactsCheckCmd(contact string) *ContactsCheckCmd {
	return &ContactsCheckCmd{
		contact:    strings.ReplaceAll(contact, "+", ""),
		resultChan: make(chan Result, 1),
	}
}

func (c *ContactsCheckCmd) Result() <-chan Result {
	return c.resultChan
}

func (c *ContactsCheckCmd) Execute(client *telegram.Client) error {
	contacts, err := client.ContactsGetContacts(0)
	if err != nil {
		c.resultChan <- Result{Success: false, Err: err}
		return err
	}

	if obj, ok := contacts.(*telegram.ContactsContactsObj); ok {
		for _, user := range obj.Users {
			if u, ok := user.(*telegram.UserObj); ok {
				if u.Phone == c.contact || u.Username == c.contact {
					c.resultChan <- Result{Success: true, Data: user}
					return nil
				}
			}
		}
	} else {
		err = fmt.Errorf("неожиданный тип ответа: %T", contacts)
		c.resultChan <- Result{Success: false, Err: err}
		return err
	}

	c.resultChan <- Result{Success: false, Data: nil, Err: nil}
	return nil
}
