package feed

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_stringCutAndEclipse(t *testing.T) {
	type args struct {
		str string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "超长中文",
			args: args{
				"一个👌🏻水电费连锁的是的是的是的是的",
			},
			want: "一个👌🏻水电费连锁的是的...",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, stringCutAndEclipse(tt.args.str), "stringCutAndEclipse(%v)", tt.args.str)
		})
	}
}
