package amqptable

import (
	"encoding/binary"
	"io"
)

type TruncatedInputStream struct {
	Reader  io.Reader
	limit   int32
	counter int32
}

func NewTruncatedInputStream(reader io.Reader, limit int32) *TruncatedInputStream {
	return &TruncatedInputStream{Reader: reader, limit: limit, counter: 0}
}

func (t *TruncatedInputStream) Available() int32 {
	return t.limit - t.counter
}

func (t *TruncatedInputStream) Read(p []byte) (n int, err error) {

	if t.counter+int32(len(p)) > t.limit {
		return 0, io.EOF
	}

	n, err = t.Reader.Read(p)
	if err != nil {
		return 0, err
	}

	t.counter += int32(n)

	return n, nil
}

func (t *TruncatedInputStream) ReadByte() (byte, error) {
	if t.counter >= t.limit {
		return 0, io.EOF
	}

	var b byte
	err := binary.Read(t.Reader, binary.BigEndian, &b)
	if err != nil {
		return 0, err
	}

	t.counter++
	return b, nil
}

func (t *TruncatedInputStream) ReadBoolean() (bool, error) {
	if t.limit <= 0 {
		return false, io.EOF
	}

	var b byte
	err := binary.Read(t.Reader, binary.BigEndian, &b)
	if err != nil {
		return false, err
	}

	t.counter++
	return b != 0, nil
}

func (t *TruncatedInputStream) ReadInt8() (int8, error) {
	var i int8
	err := binary.Read(t.Reader, binary.BigEndian, &i)
	if err != nil {
		return 0, err
	}
	t.counter += 1
	return i, nil
}

func (t *TruncatedInputStream) ReadUnsignedInt8() (uint8, error) {
	var i uint8
	err := binary.Read(t.Reader, binary.BigEndian, &i)
	if err != nil {
		return 0, err
	}
	t.counter += 1
	return i, nil
}

func (t *TruncatedInputStream) ReadInt16() (int16, error) {
	var i int16
	err := binary.Read(t.Reader, binary.BigEndian, &i)
	if err != nil {
		return 0, err
	}
	t.counter += 2
	return i, nil
}

func (t *TruncatedInputStream) ReadUnsignedInt16() (uint16, error) {
	var i uint16
	err := binary.Read(t.Reader, binary.BigEndian, &i)
	if err != nil {
		return 0, err
	}
	t.counter += 2
	return i, nil
}

func (t *TruncatedInputStream) ReadInt32() (int32, error) {
	var i int32
	err := binary.Read(t.Reader, binary.BigEndian, &i)
	if err != nil {
		return 0, err
	}
	t.counter += 4
	return i, nil
}

func (t *TruncatedInputStream) ReadUnsignedInt32() (uint32, error) {
	var i uint32
	err := binary.Read(t.Reader, binary.BigEndian, &i)
	if err != nil {
		return 0, err
	}
	t.counter += 4
	return i, nil
}

func (t *TruncatedInputStream) ReadInt64() (int64, error) {
	var i int64
	err := binary.Read(t.Reader, binary.BigEndian, &i)
	if err != nil {
		return 0, err
	}
	t.counter += 8
	return i, nil
}

func (t *TruncatedInputStream) ReadUnsignedInt64() (uint64, error) {
	var i uint64
	err := binary.Read(t.Reader, binary.BigEndian, &i)
	if err != nil {
		return 0, err
	}
	t.counter += 8
	return i, nil
}

func (t *TruncatedInputStream) ReadFloat32() (float32, error) {
	var f float32
	err := binary.Read(t.Reader, binary.BigEndian, &f)
	if err != nil {
		return 0, err
	}
	t.counter += 4
	return f, nil
}

func (t *TruncatedInputStream) ReadFloat64() (float64, error) {
	var f float64
	err := binary.Read(t.Reader, binary.BigEndian, &f)
	if err != nil {
		return 0, err
	}
	t.counter += 8
	return f, nil
}
