package httpclient

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_header_Add(t *testing.T) {
	type args struct {
		key   string
		value string
	}
	tests := []struct {
		name string
		args args
		want http.Header
	}{
		{
			name: "Test_header_Add_with_valid_key",
			args: args{
				key:   "key-example",
				value: "value-example",
			},
			want: map[string][]string{
				"Key-Example": {"value-example"},
			},
		},
		{
			name: "Test_header_Add_with_invalid_key",
			want: map[string][]string{
				"": {""},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &header{
				metadata: make(http.Header),
			}

			h.Add(tt.args.key, tt.args.value)

			assert.Equal(t, tt.want, h.metadata)
		})
	}
}
