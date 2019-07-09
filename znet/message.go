package znet

import "zframe/zinterface"

type Message struct {
	//ID (消息类型)
	ID uint32

	//
	Data []byte

	//
	DataLength uint32
}

func NewMessage(id uint32, data []byte) zinterface.IMessage {
	return &Message{
		ID:         id,
		Data:       data,
		DataLength: uint32(len(data)),
	}
}

//消息id (类型）
func (m *Message) GetID() uint32 {
	return m.ID
}

//数据
func (m *Message) GetData() []byte {
	return m.Data
}

//长度
func (m *Message) GetDataLength() uint32 {
	return m.DataLength
}
func (m *Message) SetID(id uint32) {
	m.ID = id
}
func (m *Message) SetData(data []byte) {
	m.Data = data
}
func (m *Message) SetDataLength(len uint32) {
	m.DataLength = len
}
