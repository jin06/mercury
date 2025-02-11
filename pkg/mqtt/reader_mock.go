package mqtt

import (
	"io"
	"net"
	"time"
)

type connMocker struct {
	r *io.PipeReader
	w *io.PipeWriter
}

func (p *connMocker) Read(b []byte) (n int, err error) {
	return p.r.Read(b)
}

func (p *connMocker) Write(b []byte) (n int, err error) {
	return p.w.Write(b)
}

func (p *connMocker) Close() error {
	p.w.Close()
	p.r.Close()
	return nil
}

func (p *connMocker) LocalAddr() net.Addr {
	return nil
}

func (p *connMocker) RemoteAddr() net.Addr {
	return nil
}

func (p *connMocker) SetDeadline(t time.Time) error {
	return nil
}

func (p *connMocker) SetReadDeadline(t time.Time) error {
	return nil
}

func (p *connMocker) SetWriteDeadline(t time.Time) error {
	return nil
}
