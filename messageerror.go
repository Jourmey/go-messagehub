package messagehub

import "fmt"

type MessageError struct {
	msg interface{}
}

func (m *MessageError) Error() string {
	return fmt.Sprint("MessageError panic,value =", m.msg)
}

func NewMessageError(msg interface{}) error {
	m := new(MessageError)
	m.msg = msg
	return m
}
