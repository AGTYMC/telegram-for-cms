package messenger

import (
	"fmt"

	"github.com/amarnathcjd/gogram/telegram"
)

type ContactsListCmd struct {
	resultChan chan Result
}

func NewListContactsCmd() *ContactsListCmd {
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

	var data []map[string]any

	if obj, ok := contactsObj.(*telegram.ContactsContactsObj); ok {
		fmt.Printf("[contacts] найдено сохранённых: %d\n", obj.SavedCount)

		for _, contact := range obj.Contacts {
			fmt.Printf("  • %d (mutual:%v)\n", contact.UserID, contact.Mutual)
		}

		for _, user := range obj.Users {
			if u, ok := user.(*telegram.UserObj); ok {
				name := u.FirstName
				if u.LastName != "" {
					name += " " + u.LastName
				}
				fmt.Printf("  → %d | %s | @%s | phone:%s\n",
					u.ID, name, u.Username, u.Phone)

				data = append(data, map[string]any{
					"id":       u.ID,
					"name":     name,
					"username": u.Username,
					"phone":    u.Phone,
					"mutual":   false, // можно уточнить, если нужно
				})
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
