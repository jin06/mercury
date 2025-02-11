package mqtt

import (
	"io"
	"testing"
)

func TestReader_ReadBytePtr(t *testing.T) {
	r, w := io.Pipe()
	conn := &connMocker{r: r, w: w}
	reader := NewReader(conn)

	go func() {
		w.Write([]byte{42})
		w.Close()
	}()

	b, err := reader.ReadBytePtr()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if *b != 42 {
		t.Fatalf("expected 42, got %v", *b)
	}
}

func TestReader_ReadBool(t *testing.T) {
	r, w := io.Pipe()
	conn := &connMocker{r: r, w: w}
	reader := NewReader(conn)

	go func() {
		w.Write([]byte{1})
		w.Close()
	}()

	b, err := reader.ReadBool()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if !b {
		t.Fatalf("expected true, got %v", b)
	}
}

func TestReader_ReadUint8(t *testing.T) {
	r, w := io.Pipe()
	conn := &connMocker{r: r, w: w}
	reader := NewReader(conn)

	go func() {
		w.Write([]byte{255})
		w.Close()
	}()

	u8, err := reader.ReadUint8()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if u8 != 255 {
		t.Fatalf("expected 255, got %v", u8)
	}
}

func TestReader_ReadUint16(t *testing.T) {
	r, w := io.Pipe()
	conn := &connMocker{r: r, w: w}
	reader := NewReader(conn)

	go func() {
		w.Write([]byte{0x01, 0x02})
		w.Close()
	}()

	u16, err := reader.ReadUint16()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if u16 != 258 {
		t.Fatalf("expected 258, got %v", u16)
	}
}

func TestReader_ReadUTF8Str(t *testing.T) {
	r, w := io.Pipe()
	conn := &connMocker{r: r, w: w}
	reader := NewReader(conn)

	go func() {
		w.Write([]byte{0x00, 0x05, 0x68, 0x65, 0x6c, 0x6c, 0x6f})
		w.Close()
	}()

	str, n, err := reader.ReadUTF8Str()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if str != "hello" {
		t.Fatalf("expected 'hello', got %s", str)
	}
	if n != 7 {
		t.Fatalf("expected 7 bytes read, got %d", n)
	}
}

func TestReader_ReadVariableByteInteger(t *testing.T) {
	r, w := io.Pipe()
	conn := &connMocker{r: r, w: w}
	reader := NewReader(conn)

	go func() {
		w.Write([]byte{0x81, 0x01})
		w.Close()
	}()

	length, n, err := reader.ReadVariableByteInteger()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if length != 129 {
		t.Fatalf("expected 129, got %d", length)
	}
	if n != 2 {
		t.Fatalf("expected 2 bytes read, got %d", n)
	}
}
