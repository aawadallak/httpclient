package httpclient

import "io"

type body struct {
	raw io.ReadCloser
}

func (b *body) Raw() io.Reader {
	return b.raw
}

func (b *body) Close() error {
	return b.raw.Close()
}

func (b *body) Bytes() ([]byte, error) {
	body, err := io.ReadAll(b.raw)
	if err != nil {
		return nil, err
	}

	if err := b.raw.Close(); err != nil {
		return nil, err
	}
	
	b.raw = io.NopCloser(bytes.NewBuffer(body))

	return body, nil
}
