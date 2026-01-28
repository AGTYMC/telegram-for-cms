package messenger

import "github.com/amarnathcjd/gogram/telegram"

// TODO: дублируется код получения и проверки контактов
func checkContact(phone string, client *telegram.Client) (*telegram.UserObj, error) {
	contacts, err := client.ContactsGetContacts(0)
	if err != nil {
		return nil, err
	}

	if obj, ok := contacts.(*telegram.ContactsContactsObj); ok {
		for _, user := range obj.Users {
			if u, ok := user.(*telegram.UserObj); ok {
				if u.Phone == phone || u.Username == phone {
					return u, nil
				}
			}
		}
	} else {
		return nil, err
	}

	return nil, nil
}

func getContact(userId int64, client *telegram.Client) (*telegram.UserObj, error) {
	user, err := client.GetUser(userId)
	if err != nil {
		return nil, err
	}
	return user, nil
}
