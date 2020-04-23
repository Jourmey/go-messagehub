package messagehub


type Message interface {
	GetMessageID() MessageID
}