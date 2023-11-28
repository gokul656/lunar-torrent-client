package transport

import (
	"encoding/binary"
	"fmt"
	"io"
)

type MessageID uint8

func (m MessageID) Name() string {
	switch m {
	case 0:
		return "choke"
	case 1:
		return "unchoke"
	case 2:
		return "intrested"
	case 3:
		return "not-intrested"
	case 4:
		return "have"
	case 5:
		return "bitfield"
	case 6:
		return "request"
	case 7:
		return "piece"
	case 8:
		return "cancel"
	default:
		return fmt.Sprintf("invalid message ID %d", m)
	}
}

const (
	MsgChoke         MessageID = 0
	MsgUnchoke       MessageID = 1
	MsgInterested    MessageID = 2
	MsgNotInterested MessageID = 3
	MsgHave          MessageID = 4
	MsgBitfield      MessageID = 5
	MsgRequest       MessageID = 6
	MsgPiece         MessageID = 7
	MsgCancel        MessageID = 8
)

type Message struct {
	ID      MessageID
	Payload Bitfield
}

func ReadMessage(conn io.Reader) (*Message, error) {
	lengthBuf := make([]byte, 4) // made of 32 bit int, so the size will be 4 byte
	_, err := io.ReadFull(conn, lengthBuf)
	if err != nil {
		return nil, err
	}

	length := binary.BigEndian.Uint32(lengthBuf)

	msgBuf := make([]byte, length)
	_, err = io.ReadFull(conn, msgBuf)
	if err != nil {
		return nil, err
	}

	msg := &Message{
		ID:      MessageID(msgBuf[0]),
		Payload: msgBuf[1:],
	}

	return msg, nil
}

func (m *Message) Serialize() []byte {
	length := uint32(len(m.Payload) + 1)
	buf := make([]byte, 4+length)
	binary.BigEndian.PutUint32(buf[0:4], length)
	buf[5] = byte(m.ID)
	copy(buf[6:], m.Payload)

	return buf
}
