package messagehub

import (
	uuid "github.com/nu7hatch/gouuid"
	"sync"
)

//TODO :最大并发数限制

type MessageHub interface {
	Publish(message Message)
	Subscribe(messageID MessageID, messageHandler MessageHandler) (HandlerID, error)
	Unsubscribe(messageID MessageID, handlerID HandlerID)
	IsSubscribed(messageID MessageID, handlerID HandlerID) bool
}

type (
	MessageID [16]byte
	HandlerID [16]byte

	GlobalMessageHandler      func(MessageID, Message)
	GlobalErrorMessageHandler func(MessageID, error)
	MessageHandler            func(Message)
)

func NewMessageHub(messageHandler GlobalMessageHandler, errorHandler GlobalErrorMessageHandler) MessageHub {

	m := new(messageHub)
	m.globalHandler = messageHandler
	m.errorHandler = errorHandler
	m.messageMap = make(map[MessageID]map[HandlerID]MessageHandler)
	m.lock = new(sync.Mutex)

	return m
}

type messageHub struct {
	globalHandler GlobalMessageHandler
	errorHandler  GlobalErrorMessageHandler

	messageMap map[MessageID]map[HandlerID]MessageHandler
	lock       sync.Locker
}

func (m *messageHub) Publish(message Message) {
	messageID := message.GetMessageID()
	if m.globalHandler != nil {
		m.globalHandler(messageID, message)
	}

	if mes, ok := m.messageMap[messageID]; ok {

		for _, handler := range mes {
			go func() {
				defer func() {
					if p := recover(); p != nil {
						if m.errorHandler != nil {
							m.errorHandler(messageID, NewMessageError(p))
						}
					}
				}()
				handler(message)
			}()
		}
	}

}

func (m *messageHub) Subscribe(messageID MessageID, messageHandler MessageHandler) (HandlerID, error) {
	if id, err := uuid.NewV4(); err != nil {
		return HandlerID{}, err
	} else {
		uuid := HandlerID(*id)

		m.lock.Lock()
		if mes, ok := m.messageMap[messageID]; ok {
			mes[uuid] = messageHandler
		} else {
			m.messageMap[messageID] = map[HandlerID]MessageHandler{uuid: messageHandler}
		}
		m.lock.Unlock()
		return uuid, nil
	}
}

func (m *messageHub) Unsubscribe(messageID MessageID, handlerID HandlerID) {
	if m.IsSubscribed(messageID, handlerID) {
		m.lock.Lock()
		delete(m.messageMap[messageID], handlerID)
		m.lock.Unlock()
	}
}

func (m *messageHub) IsSubscribed(messageID MessageID, handlerID HandlerID) bool {
	if mes, ok := m.messageMap[messageID]; ok {
		_, ok := mes[handlerID]
		return ok
	} else {
		return false
	}
}
