package i18n

import (
	"errors"
	"fmt"
	"strings"

	"github.com/leonelquinteros/gotext"
)

type Message int
type MessageMap map[Message]string

func gettext(s string) string {
	return gotext.Get(s)
}

func Msg(id Message, params ...interface{}) string {
	msg := Messages[id]

	if strings.Contains(msg, "%") {
		msg = fmt.Sprintf(msg, params...)
	}

	return msg
}

func Error(id Message, params ...interface{}) error {
	return errors.New(Msg(id, params...))
}
