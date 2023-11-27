package main

import (
	"encoding/binary"
	"io"
)

type messageID uint8

const (
	MsgChoke         messageID = 0
	MsgUnchoke       messageID = 1
	MsgInterested    messageID = 2
	MsgNotInterested messageID = 3
	MsgHave          messageID = 4
	MsgBitfield      messageID = 5
	MsgRequest       messageID = 6
	MsgPiece         messageID = 7
	MsgCancel        messageID = 8
)

type Message struct {
	ID      messageID
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
		ID:      messageID(msgBuf[0]),
		Payload: msgBuf[1:],
	}

	return msg, nil
}

func (m *Message) Serialize() []byte {
	buf := make([]byte, 4)
	return buf
}
