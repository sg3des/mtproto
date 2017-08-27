package mtproto

import (
	"fmt"
	"reflect"
)

const (
	CHAT_TYPE_EMPTY             = "EMPTY"
	CHAT_TYPE_CHAT              = "CHAT"
	CHAT_TYPE_CHAT_FORBIDDEN    = "CHAT_FORBIDDEN"
	CHAT_TYPE_CHANNEL           = "CHANNEL"
	CHAT_TYPE_CHANNEL_FORBIDDEN = "CHANNEL_FORBIDDEN"
)

type ChatProfilePhoto struct {
	PhotoSmall FileLocation
	PhotoBig   FileLocation
}
type Chat struct {
	flags        int32
	ID           int32
	Username     string
	Type         string
	Title        string
	Photo        *ChatProfilePhoto
	Participants int32
	Members 	  []ChatMember
	Date         int32
	Left         bool
	Version      int32
	AccessHash   int64
	Address      string
	Venue        string
	CheckedIn    bool
}
type ChatMember struct {
	UserID 		int32
	InviterID	int32
	Date 		int32
}
type ChannelParticipantFilter struct {}

func (ch *Chat) GetPeer() TL {
	switch ch.Type {
	case CHAT_TYPE_CHAT, CHAT_TYPE_CHAT_FORBIDDEN:
		return TL_peerChat{
			Chat_id: ch.ID,
		}
	case CHAT_TYPE_CHANNEL, CHAT_TYPE_CHANNEL_FORBIDDEN:
		return TL_peerChannel{
			Channel_id: ch.ID,

		}
	default:
		return nil
	}
}
func (ch *Chat) GetInputPeer() TL {
	switch ch.Type {
	case CHAT_TYPE_CHAT, CHAT_TYPE_CHAT_FORBIDDEN:
		return TL_inputPeerChat{
			Chat_id: ch.ID,
		}
	case CHAT_TYPE_CHANNEL, CHAT_TYPE_CHANNEL_FORBIDDEN:
		return TL_inputPeerChannel{
			Channel_id:  ch.ID,
			Access_hash: ch.AccessHash,
		}
	default:
		return nil
	}
}

func NewChatProfilePhoto(input TL) (photo *ChatProfilePhoto) {
	photo = new(ChatProfilePhoto)
	switch p := input.(type) {
	case TL_chatPhotoEmpty:
		return nil
	case TL_chatPhoto:
		switch big := p.Photo_big.(type) {
		case TL_fileLocationUnavailable:
		case TL_fileLocation:
			photo.PhotoBig.DC = big.Dc_id
			photo.PhotoBig.LocalID = big.Local_id
			photo.PhotoBig.Secret = big.Secret
			photo.PhotoBig.VolumeID = big.Volume_id
		}
		switch small := p.Photo_small.(type) {
		case TL_fileLocationUnavailable:
		case TL_fileLocation:
			photo.PhotoSmall.DC = small.Dc_id
			photo.PhotoSmall.LocalID = small.Local_id
			photo.PhotoSmall.Secret = small.Secret
			photo.PhotoSmall.VolumeID = small.Volume_id
		}
	}
	return photo
}
func NewChat(input TL) (chat *Chat) {
	chat = new(Chat)
	chat.Members = []ChatMember{}
	switch ch := input.(type) {
	case TL_chatEmpty:
		chat.Type = CHAT_TYPE_EMPTY
		chat.ID = ch.Id
	case TL_chatForbidden:
		chat.Type = CHAT_TYPE_CHAT_FORBIDDEN
		chat.ID = ch.Id
		chat.Title = ch.Title
	case TL_chat:
		chat.flags = ch.Flags
		chat.Type = CHAT_TYPE_CHAT
		chat.ID = ch.Id
		chat.Title = ch.Title
		chat.Date = ch.Date
		chat.Photo = NewChatProfilePhoto(ch.Photo)
		chat.Version = ch.Version
		chat.Participants = ch.Participants_count
	case TL_chatFull:
		chat.ID = ch.Id
		participants := ch.Participants.(TL_chatParticipants)
		chat.Version = participants.Version
		for _, tl := range ch.Participants.(TL_chatParticipants).Participants {
			m := tl.(TL_chatParticipant)
			chat.Members = append(chat.Members, ChatMember{m.User_id, m.Inviter_id, m.Date})
		}
	case TL_channelFull:

	case TL_channelForbidden:
		chat.flags = ch.Flags
		chat.Type = CHAT_TYPE_CHANNEL_FORBIDDEN
		chat.ID = ch.Id
		chat.Title = ch.Title
		chat.AccessHash = ch.Access_hash
	case TL_channel:
		chat.flags = ch.Flags
		chat.Type = CHAT_TYPE_CHANNEL
		chat.ID = ch.Id
		chat.Username = ch.Username
		chat.Title = ch.Title
		chat.Date = ch.Date
		chat.Photo = NewChatProfilePhoto(ch.Photo)
		chat.Version = ch.Version
		chat.AccessHash = ch.Access_hash
	default:
		fmt.Println(reflect.TypeOf(ch).String())
		return nil
	}
	return chat

}
func NewInputPeerUser(userID int32, accessHash int64) TL {
	return TL_inputPeerUser{
		User_id: userID,
		Access_hash: accessHash,
	}
}
func NewInputPeerChat(chatID int32) TL {
	return TL_inputPeerChat{
		Chat_id: chatID,
	}
}
func NewInputPeerChannel(channelID int32, accessHash int64) TL {
	return TL_inputPeerChannel{
		Channel_id:  channelID,
		Access_hash: accessHash,
	}
}
func NewPeerChat(chatID int32) TL {
	return TL_peerChat{
		Chat_id: chatID,
	}
}
func NewPeerChannel(channelID int32) TL {
	return TL_peerChannel{
		Channel_id: channelID,
	}
}