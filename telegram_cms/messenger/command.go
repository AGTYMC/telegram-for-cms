package messenger

import (
	"github.com/amarnathcjd/gogram/telegram"
)

type Command interface {
	Execute(*telegram.Client) error
	Result() <-chan Result
}
