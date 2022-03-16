package helpers

import (
	"github.com/stretchr/testify/assert"
	"ssi-gitlab.teda.th/ssi/core/utils"
	"testing"
)

func TestWriteNestedJson(t *testing.T) {
	type args struct {
		nestFields []string
		value      interface{}
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "string",
			args: args{
				nestFields: []string{"foo", "bar"},
				value:      "message",
			},
			want: `{"foo":{"bar":"message"}}`,
		},
		{
			name: "bool",
			args: args{
				nestFields: []string{"foo", "bar"},
				value:      true,
			},
			want: `{"foo":{"bar":true}}`,
		},
		{
			name: "number",
			args: args{
				nestFields: []string{"foo", "bar"},
				value:      1,
			},
			want: `{"foo":{"bar":1}}`,
		},
		{
			name: "null",
			args: args{
				nestFields: []string{"foo", "bar"},
				value:      nil,
			},
			want: `{"foo":{"bar":null}}`,
		},
		{
			name: "single field",
			args: args{
				nestFields: []string{"foo"},
				value:      "bar",
			},
			want: `{"foo":"bar"}`,
		}, {
			name: "empty",
			args: args{
				nestFields: []string{},
				value:      "bar",
			},
			want: `{}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := WriteNestedJson(tt.args.nestFields, tt.args.value); got != tt.want {
				assert.Equal(t, tt.want, got)
				var m map[string]interface{}
				err := utils.JSONParse(utils.StringToBytes(got), &m)
				assert.NoError(t, err)
			}
		})
	}
}
