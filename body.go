package httpclient

import "io"

type body struct {
	raw io.ReadCloser
}

// Raw returns the payload without reading it.
// The `Close` should be called, otherwise,
// it's never will be released from the `net/http` pool.
func (b *body) Raw() io.Reader {
	return b.raw
}

// Close release the resource from `net/http`
// and close the payload.
//
// It should be use only with `Raw`.
func (b *body) Close() error {
	return b.raw.Close()
}

// Bytes reads the payload and returns the content.
//
// It could return an error while reading or closing the payload
func (b *body) Bytes() ([]byte, error) {
	body, err := io.ReadAll(b.raw)
	if err != nil {
		return nil, err
	}

	if err := b.raw.Close(); err != nil {
		return nil, err
	}

	return body, nil
}
