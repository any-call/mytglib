package mytglib

import (
	"fmt"
	"math"
)

type api struct {
}

func ImpApi() api {
	return api{}
}

func (self api) GetChatList(client *Client, limit int) ([]*Chat, error) {
	if client == nil {
		return nil, fmt.Errorf("invalid client")
	}

	allChats := make([]*Chat, 0)
	var haveFullChatList bool = false

	for !haveFullChatList && limit > len(allChats) {
		offsetOrder := int64(math.MaxInt64)
		offsetChatID := int64(0)

		var chatList = NewChatListMain()
		var lastChat *Chat

		if len(allChats) > 0 {
			lastChat = allChats[len(allChats)-1]
			for i := 0; i < len(lastChat.Positions); i++ {
				//Find the main chat list
				if lastChat.Positions[i].List.GetChatListEnum() == ChatListMainType {
					offsetOrder = int64(lastChat.Positions[i].Order)
				}
			}
			offsetChatID = lastChat.ID
		}

		// get chats (ids) from tdlib
		chats, err := client.GetChats(chatList, JSONInt64(offsetOrder),
			offsetChatID, int32(limit-len(allChats)))
		if err != nil {
			return nil, err
		}

		if len(chats.ChatIDs) == 0 {
			haveFullChatList = true
			return allChats, nil
		}

		for _, chatID := range chats.ChatIDs {
			// get chat info from tdlib
			chat, err := client.GetChat(chatID)
			if err == nil {
				allChats = append(allChats, chat)
			} else {
				return nil, err
			}
		}
	}

	return allChats, nil
}
