package progress

import (
	"io"
	"sync"
)

// ReadCloser counts the bytes read through it.
type ReadCloser struct {
	r io.ReadCloser

	lock sync.RWMutex // protects n and err
	n    int64
	err  error
}

// NewReadCloser makes a new Reader that counts the bytes
// read through it.
func NewReadCloser(r io.ReadCloser) *ReadCloser {
	return &ReadCloser{
		r: r,
	}
}

func (r *ReadCloser) Close() error {
	return r.r.Close()
}

func (r *ReadCloser) Read(p []byte) (n int, err error) {
	n, err = r.r.Read(p)
	r.lock.Lock()
	r.n += int64(n)
	r.err = err
	r.lock.Unlock()
	return
}

// N gets the number of bytes that have been read
// so far.
func (r *ReadCloser) N() int64 {
	var n int64
	r.lock.RLock()
	n = r.n
	r.lock.RUnlock()
	return n
}

// Err gets the last error from the Reader.
func (r *ReadCloser) Err() error {
	var err error
	r.lock.RLock()
	err = r.err
	r.lock.RUnlock()
	return err
}
