package maskingutil

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMaskPassword(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name string
		v    string
		want string
	}{
		{
			name: "dsn",
			v:    "root:mypassword123@tcp(127.0.0.1:3306)/mydb?charset=utf8",
			want: "root:****@tcp(127.0.0.1:3306)/mydb?charset=utf8",
		}, {
			name: "http",
			v:    "https://user1:mypassword123@example.com",
			want: "https://user1:****@example.com",
		}, {
			name: "not url",
			v:    "plain text message",
			want: "plain text message",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := MaskPassword(tc.v)
			assert.Equal(t, tc.want, got)
		})
	}
}
