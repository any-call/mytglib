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

func (self api) SendMessage(client *Client, chatId, replyToMessage int64, message string) (*Message, error) {
	if client == nil {
		return nil, fmt.Errorf("invalid client")
	}

	return client.SendMessage(chatId, 0, replyToMessage, nil, nil,
		NewInputMessageText(NewFormattedText(message, nil), true, true))
}

func (self api) SendDice(client *Client, chatId, replyToMessage int64) (*Message, error) {
	if client == nil {
		return nil, fmt.Errorf("invalid client")
	}

	return client.SendMessage(chatId, 0, replyToMessage, nil, nil,
		NewInputMessageDice("🎲", true))
}

func (self api) DelMessage(client *Client, chatID, messageID int64) error {
	if client == nil {
		return fmt.Errorf("invalid client")
	}

	_, err := client.DeleteMessages(chatID, []int64{messageID}, true)
	return err
}
