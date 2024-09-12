package amqptable

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"time"
)

func ReadTable(data []byte) (map[string]any, error) {
	byteReader := bytes.NewReader(data)
	reader := bufio.NewReader(byteReader)

	var tableLength int32
	binary.Read(reader, binary.BigEndian, &tableLength)

	table := make(map[string]any)
	tableIn := NewTruncatedInputStream(reader, tableLength)

	for tableIn.Available() > 0 {
		name, err := readShortstr(tableIn)
		if err != nil {
			return nil, err
		}
		value, err := readFieldValue(tableIn)
		if err != nil {
			return nil, err
		}
		if _, ok := table[name]; !ok {
			table[name] = value
		}
	}
	return table, nil
}

func WriteTable(table map[string]any) ([]byte, error) {
	buf := new(bytes.Buffer)
	writer := bufio.NewWriter(buf)

	for key, value := range table {
		writeShortstr(writer, key)
		writeFieldValue(writer, value)
	}
	writer.Flush()

	msgBody := buf.Bytes()
	length := int32(len(msgBody))
	lengthBuf := new(bytes.Buffer)
	lengthWriter := bufio.NewWriter(lengthBuf)

	binary.Write(lengthWriter, binary.BigEndian, length)
	lengthWriter.Flush()
	lengthBytes := lengthBuf.Bytes()
	return append(lengthBytes, msgBody...), nil
}

func writeShortstr(writer *bufio.Writer, str string) error {
	length := int8(len(str))
	binary.Write(writer, binary.BigEndian, length)
	writer.WriteString(str)
	return nil
}

func writeLongstr(writer *bufio.Writer, str string) error {
	length := int32(len(str))
	binary.Write(writer, binary.BigEndian, length)
	writer.WriteString(str)
	return nil
}

func writeTimestamp(writer *bufio.Writer, timestamp time.Time) error {
	seconds := int32(timestamp.Unix())
	nanoseconds := int32(timestamp.Nanosecond())
	binary.Write(writer, binary.BigEndian, seconds)
	binary.Write(writer, binary.BigEndian, nanoseconds)
	return nil
}

func writeFieldValue(writer *bufio.Writer, value any) error {
	switch v := value.(type) {
	case string:
		writer.WriteByte('S')
		writeLongstr(writer, v)
	case int8:
		writer.WriteByte('b')
		binary.Write(writer, binary.BigEndian, v)
	case int16:
		writer.WriteByte('s')
		binary.Write(writer, binary.BigEndian, v)
	case int32:
		writer.WriteByte('I')
		binary.Write(writer, binary.BigEndian, v)
	case int64:
		writer.WriteByte('l')
		binary.Write(writer, binary.BigEndian, v)
	case float64:
		writer.WriteByte('d')
		binary.Write(writer, binary.BigEndian, v)
	case float32:
		writer.WriteByte('f')
		binary.Write(writer, binary.BigEndian, v)
	case bool:
		writer.WriteByte('t')
		binary.Write(writer, binary.BigEndian, v)
	case []byte:
		writer.WriteByte('x')
		binary.Write(writer, binary.BigEndian, v)
	case time.Time:
		writer.WriteByte('T')
		writeTimestamp(writer, v)
	default:
		return fmt.Errorf("unsupported type: %T", v)
	}
	return nil
}

func readShortstr(reader *TruncatedInputStream) (string, error) {
	length, err := reader.ReadByte()
	if err != nil {
		return "", err
	}
	bytes := make([]byte, int(length))
	_, err = reader.Read(bytes)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func readFieldValue(reader *TruncatedInputStream) (any, error) {
	value := any(nil)

	b, err := reader.ReadByte()
	if err != nil {
		return nil, err
	}

	switch b {
	case 'S':
		value, err = readLongstr(reader)
	case 'I':
		value, err = reader.ReadInt32()
	case 'i':
		value, err = reader.ReadUnsignedInt32()
	case 'T':
		value, err = readTimestamp(reader)
		/*
			case 'D':
				scale := reader.ReadByte()
				unscaled := make([]byte, 4)
				_, err = reader.Read(unscaled)
				if err != nil {
					return nil, err
				}
				value = big.NewInt(0).SetBytes(unscaled)
			case 'F':
				value, err = readTable(reader)
			case 'A':
				value, err = readArray(reader)
		*/
	case 'b':
		value, err = reader.ReadInt8()
	case 'B':
		value, err = reader.ReadUnsignedInt8()
	case 'd':
		value, err = reader.ReadFloat64()
	case 'f':
		value, err = reader.ReadFloat32()
	case 'l':
		value, err = reader.ReadInt64()
	case 's':
		value, err = reader.ReadInt16()
	case 'u':
		value, err = reader.ReadUnsignedInt16()
	case 't':
		value, err = reader.ReadBoolean()
	case 'x':
		value, err = readBytes(reader)
	case 'V':
		value = nil
	default:
		return nil, fmt.Errorf("unrecognised type in table")
	}
	if err != nil {
		return nil, err
	}

	return value, nil
}

func readLongstr(reader *TruncatedInputStream) (string, error) {
	length, err := reader.ReadInt32()
	if err != nil {
		return "", err
	}
	bytes := make([]byte, length)
	_, err = reader.Read(bytes)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func readBytes(reader *TruncatedInputStream) ([]byte, error) {
	length, err := reader.ReadInt32()
	if err != nil {
		return nil, err
	}
	bytes := make([]byte, length)
	_, err = reader.Read(bytes)
	if err != nil {
		return nil, err
	}
	return bytes, nil
}

func readTimestamp(reader *TruncatedInputStream) (time.Time, error) {
	seconds, err := reader.ReadInt32()
	if err != nil {
		return time.Time{}, err
	}
	nanoseconds, err := reader.ReadInt32()
	if err != nil {
		return time.Time{}, err
	}
	return time.Unix(int64(seconds), int64(nanoseconds)), nil
}
