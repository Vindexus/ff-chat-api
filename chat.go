package userchat

import (
	"fmt"
	"time"

	"github.com/Vindexus/userchat-api/pkg/fixtures"
)

type Message struct {
	Date     time.Time
	Message  string
	SenderId int
	Username string
}

// TODO: Don't use memory to store messages, obviously this should be in a database
var DataMessages = map[string][]*Message{
	GetChatId(fixtures.User1Id, fixtures.User2Id): []*Message{
		{
			SenderId: fixtures.User1Id,
			Date:     time.Now().Add(time.Minute * -1),
			Message:  "Hello, how are you?",
			Username: fixtures.User1Username, // This would normally come from a JOIN of messages table to users table
		},
		{
			SenderId: fixtures.User2Id,
			Date:     time.Now(),
			Message:  "I am fine",
			Username: fixtures.User2Username,
		},
	},
}

func GetUsersChatId(user1, user2 *User) string {
	return GetChatId(user1.Id, user2.Id)
}

func GetChatId(id1, id2 int) string {
	if id1 < id2 {
		return fmt.Sprintf("%d_%d", id1, id2)
	}
	return fmt.Sprintf("%d_%d", id2, id1)
}

// TODO: Hook this up to an actual database to fetch the chat messages
func GetChatMessages(currentUser *User, otherUser *User) ([]*Message, error) {
	id := GetUsersChatId(currentUser, otherUser)

	_, ok := DataMessages[id]
	if !ok {
		DataMessages[id] = make([]*Message, 0)
	}

	return DataMessages[id], nil
}
