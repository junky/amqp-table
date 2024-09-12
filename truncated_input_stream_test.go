package amqptable

import (
	"bytes"
	"io"
	"testing"
)

func TestTruncatedInputStream_Read_Basic(t *testing.T) {
	reader := bytes.NewReader([]byte{1, 2, 3, 4, 5})
	truncatedReader := NewTruncatedInputStream(reader, 3)

	buf := make([]byte, 2)
	n, err := truncatedReader.Read(buf)
	if err != nil {
		t.Fatalf("Read error: %v", err)
	}
	if n != 2 {
		t.Fatalf("Read returned %d bytes, expected 2", n)
	}
}

func TestNewTruncatedInputStream(t *testing.T) {
	reader := bytes.NewReader([]byte{1, 2, 3, 4, 5})
	tis := NewTruncatedInputStream(reader, 3)

	if tis.limit != 3 {
		t.Errorf("Expected limit to be 3, got %d", tis.limit)
	}
	if tis.counter != 0 {
		t.Errorf("Expected counter to be 0, got %d", tis.counter)
	}
}

func TestAvailable(t *testing.T) {
	reader := bytes.NewReader([]byte{1, 2, 3, 4, 5})
	tis := NewTruncatedInputStream(reader, 3)

	if available := tis.Available(); available != 3 {
		t.Errorf("Expected available to be 3, got %d", available)
	}

	_, _ = tis.ReadByte()
	if available := tis.Available(); available != 2 {
		t.Errorf("Expected available to be 2, got %d", available)
	}
}

func TestReadByte(t *testing.T) {
	reader := bytes.NewReader([]byte{1, 2, 3})
	tis := NewTruncatedInputStream(reader, 2)

	b, err := tis.ReadByte()
	if b != 1 || err != nil {
		t.Errorf("Expected (1, nil), got (%d, %v)", b, err)
	}

	b, err = tis.ReadByte()
	if b != 2 || err != nil {
		t.Errorf("Expected (2, nil), got (%d, %v)", b, err)
	}

	// Try to read more, should return EOF
	b, err = tis.ReadByte()
	if b != 0 || err != io.EOF {
		t.Errorf("Expected (0, EOF), got (%d, %v)", b, err)
	}
}

func TestReadBoolean(t *testing.T) {
	reader := bytes.NewReader([]byte{0, 1, 2})
	tis := NewTruncatedInputStream(reader, 3)

	b, err := tis.ReadBoolean()
	if b != false || err != nil {
		t.Errorf("Expected (false, nil), got (%t, %v)", b, err)
	}

	b, err = tis.ReadBoolean()
	if b != true || err != nil {
		t.Errorf("Expected (true, nil), got (%t, %v)", b, err)
	}
}

func TestReadInt32(t *testing.T) {
	reader := bytes.NewReader([]byte{0, 0, 0, 5, 0, 0, 0, 10})
	tis := NewTruncatedInputStream(reader, 8)

	i, err := tis.ReadInt32()
	if i != 5 || err != nil {
		t.Errorf("Expected (5, nil), got (%d, %v)", i, err)
	}

	i, err = tis.ReadInt32()
	if i != 10 || err != nil {
		t.Errorf("Expected (10, nil), got (%d, %v)", i, err)
	}

	// Try to read more, should return EOF
	i, err = tis.ReadInt32()
	if i != 0 || err != io.EOF {
		t.Errorf("Expected (0, EOF), got (%d, %v)", i, err)
	}
}

func TestTruncatedInputStream_Read(t *testing.T) {
	tests := []struct {
		name        string
		input       []byte
		limit       int32
		readSize    int
		expectedN   int
		expectedErr error
	}{
		{
			name:        "Read within limit",
			input:       []byte("hello world"),
			limit:       11,
			readSize:    5,
			expectedN:   5,
			expectedErr: nil,
		},
		{
			name:        "Read exactly at limit",
			input:       []byte("hello world"),
			limit:       5,
			readSize:    5,
			expectedN:   5,
			expectedErr: nil,
		},
		{
			name:        "Read beyond limit",
			input:       []byte("hello world"),
			limit:       5,
			readSize:    10,
			expectedN:   0,
			expectedErr: io.EOF,
		},
		{
			name:        "Read from empty input",
			input:       []byte{},
			limit:       5,
			readSize:    5,
			expectedN:   0,
			expectedErr: io.EOF,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader := bytes.NewReader(tt.input)
			tis := NewTruncatedInputStream(reader, tt.limit)

			p := make([]byte, tt.readSize)
			n, err := tis.Read(p)

			if n != tt.expectedN {
				t.Errorf("%s: Expected n = %d, got %d", tt.name, tt.expectedN, n)
			}

			if err != tt.expectedErr {
				t.Errorf("%s: Expected error = %v, got %v", tt.name, tt.expectedErr, err)
			}

			if err == nil && !bytes.Equal(p[:n], tt.input[:n]) {
				t.Errorf("%s: Expected read data = %v, got %v", tt.name, tt.input[:n], p[:n])
			}
		})
	}
}
