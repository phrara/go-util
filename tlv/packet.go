package tlv

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"
)

const HEADER = 8

// Packet
// "TLV" style binary packets: Tag | Len | Value
type Packet struct {
	Tag   uint32
	Len   uint32
	Value []byte
}

func (p *Packet) Load(reader io.Reader) error {
	header := make([]byte, HEADER)
	_, err := io.ReadFull(reader, header)
	if err != nil {
		return err
	}
	err = p.ParseHeader(header)
	if err != nil {
		return err
	}
	val := make([]byte, p.Len)
	_, err = io.ReadFull(reader, val)
	if err != nil {
		return err
	}
	p.Value = val
	return nil
}

func (p *Packet) Wrap() ([]byte, error) {
	buf := bytes.Buffer{}
	err := binary.Write(&buf, binary.LittleEndian, p.Tag)
	if err != nil {
		return nil, err
	}
	err = binary.Write(&buf, binary.LittleEndian, p.Len)
	if err != nil {
		return nil, err
	}
	err = binary.Write(&buf, binary.LittleEndian, p.Value)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (p *Packet) Parse(data []byte) error {
	if len(data) < HEADER {
		return errors.New("too short")
	}
	err := p.ParseHeader(data[:HEADER])
	if err != nil {
		return err
	}
	if int(p.Len) != len(data[HEADER:]) {
		return errors.New("illegal data packet")
	} else {
		p.Value = data[HEADER:]
		return nil
	}
}

func (p *Packet) ParseHeader(header []byte) error {
	if len(header) != HEADER {
		return errors.New("illegal length of header")
	}
	buf := bytes.NewBuffer(header)
	err := binary.Read(buf, binary.LittleEndian, &(p.Tag))
	if err != nil {
		return err
	}
	err = binary.Read(buf, binary.LittleEndian, &(p.Len))
	if err != nil {
		return err
	}
	return nil
}
