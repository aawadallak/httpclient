package httpclient

type (
	ResponseOption func(*response)

	response struct {
		contentLenght int64
		status        int
		body          Body
		header        Header
	}
)

var _ Response = (*response)(nil)

func (r *response) ContentLength() int64 {
	return r.contentLenght
}

func (r *response) Body() Body {
	return r.body
}

func (r *response) StatusCode() int {
	return r.status
}

func (r *response) Header() Header {
	return r.header
}
