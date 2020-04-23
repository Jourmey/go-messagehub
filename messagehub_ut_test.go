package messagehub

import (
	"fmt"
	"testing"
)

type TestMessageA struct {
	context string
}

func (t TestMessageA) GetMessageID() MessageID {

	b := [16]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16,}

	return MessageID(b)
}

func TestNewMessageHub(t *testing.T) {

	mh := NewMessageHub(func(id MessageID, message Message) {
		fmt.Println("Global, message =", message)
	}, func(id MessageID, err error) {
		fmt.Println("Global, error =", err)
	})

	a := TestMessageA{context: "TestMessageA???"}

	aid := a.GetMessageID()

	funcId := 0
	for i := 0; i < 100; i++ {
		handID, err := mh.Subscribe(aid, func(mes Message) {
			fmt.Println("No.", funcId, "Subscribe, value =", mes)
			funcId++
			if funcId%3 == 0 {
				panic(funcId)
			}
		})
		if err != nil {
			t.Fail()
		}
		if mh.IsSubscribed(aid, handID) == false {
			t.Fail()
		}

	}

	mh.Publish(a)

}
