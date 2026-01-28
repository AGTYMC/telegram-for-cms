package messenger

import (
	"fmt"

	"github.com/amarnathcjd/gogram/telegram"
)

type ContactsAddCmd struct {
	phone      string
	username   string
	firstName  string
	lastName   string
	resultChan chan Result
}

func NewContactsAddCmd(phone, username, firstName, lastName string) *ContactsAddCmd {
	return &ContactsAddCmd{
		phone:      phone,
		username:   username,
		firstName:  firstName,
		lastName:   lastName,
		resultChan: make(chan Result, 1),
	}
}

func (c *ContactsAddCmd) Result() <-chan Result {
	return c.resultChan
}

func (c *ContactsAddCmd) Execute(client *telegram.Client) error {
	if user, err := checkContact(c.phone, client); err == nil && user != nil {
		c.resultChan <- Result{Success: true, Data: user}
		return nil
	}

	imported, err := client.ContactsImportContacts([]*telegram.InputPhoneContact{
		{
			Phone:     c.phone,
			FirstName: c.firstName,
			LastName:  c.lastName,
			ClientID:  0,
		},
	})

	if err != nil {
		c.resultChan <- Result{Success: false, Err: fmt.Errorf("resolve %q: %w", c.phone, err)}
		return err
	}

	if len(imported.Imported) == 0 {
		c.resultChan <- Result{Success: false, Err: fmt.Errorf("imported contact list is empty")}
		return fmt.Errorf("imported contact list is empty")
	}

	user, err := getUser(imported.Imported[0].UserID, client)
	if err != nil {
		c.resultChan <- Result{Success: false, Err: fmt.Errorf("resolve %q: %w", c.phone, err)}
		return err
	}

	c.resultChan <- Result{Success: true, Data: user}

	return nil
}
