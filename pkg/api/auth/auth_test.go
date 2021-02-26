package auth

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_hash(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Should match",
			args: args{
				s: "hello-world",
			},
			want: "afa27b44d43b02a9fea41d13cedc2e4016cfcf87c5dbf990e593669aa8ce286d",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := hash(tt.args.s)
			assert.Equal(t, tt.want, res)
		})
	}
}
