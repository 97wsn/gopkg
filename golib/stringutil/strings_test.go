package stringutil

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func BenchmarkGetRandomString(b *testing.B) {
	for i := 0; i < b.N; i++ {
		RandomString(8)
	}
}

func TestGetRandomString(t *testing.T) {
	for i := 0; i < 10; i++ {
		s := RandomString(8)
		if len(s) != 8 {
			t.Error("string length error:" + s)
		}
	}
}

func TestTitleCasedName(t *testing.T) {
	type args struct {
		name string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "test1",
			args: args{"userID"},
			want: "Userid",
		}, {
			name: "test2",
			args: args{"UserName"},
			want: "Username",
		}, {
			name: "test3",
			args: args{"user_id"},
			want: "UserId",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := TitleCasedName(tt.args.name); got != tt.want {
				t.Errorf("TitleCasedName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSnakeCasedName(t *testing.T) {
	type args struct {
		name string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "test1",
			args: args{"UserId"},
			want: "user_id",
		}, {
			name: "test2",
			args: args{"userID"},
			want: "user_i_d",
		}, {
			name: "test3",
			args: args{"user_id"},
			want: "user_id",
		}, {
			name: "test4",
			args: args{"user_ID"},
			want: "user_i_d",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SnakeCasedName(tt.args.name); got != tt.want {
				t.Errorf("SnakeCasedName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMustJsonEncode(t *testing.T) {
	m := map[string]any{
		"id":     1,
		"name":   "test",
		"status": true,
	}
	s := MustJsonEncode(m)
	assert.Equal(t, `{"id":1,"name":"test","status":true}`, s)
}

func TestLeftpad(t *testing.T) {
	s := Leftpad("test", 10)
	assert.Equal(t, "      test", s)
	var a rune = 42
	s2 := Leftpad("test", 7, a)
	assert.Equal(t, "***test", s2)

}

func TestTrimBom(t *testing.T) {
	s := TrimBom(string([]byte{239, 187, 191}) + "abcd")
	assert.Equal(t, "abcd", s)
	s = TrimBom("E")
	assert.Equal(t, "E", s)
}

func TestRightPad(t *testing.T) {
	type args struct {
		str     string
		length  int
		padChar string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "test1",
			args: args{
				str:     "test",
				length:  10,
				padChar: "0",
			},
			want: "test000000",
		},
		{
			name: "test2",
			args: args{
				str:     "testtesttest",
				length:  10,
				padChar: "0",
			},
			want: "testtesttest",
		},
		{
			name: "test3",
			args: args{
				str:     "testtestte",
				length:  10,
				padChar: "0",
			},
			want: "testtestte",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, RightPad(tt.args.str, tt.args.length, tt.args.padChar), "RightPad(%v, %v, %v)", tt.args.str, tt.args.length, tt.args.padChar)
		})
	}
}

func TestMark(t *testing.T) {
	type args struct {
		str   string
		start int
		end   int
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "保留第一位和最后两位",
			args: args{
				str:   "1234",
				start: 1,
				end:   2,
			},
			want: "1****34",
		},
		{
			name: "保留第一位和最后两位",
			args: args{
				str:   "123456789",
				start: 1,
				end:   2,
			},
			want: "1****89",
		},
		{
			name: "保留前三位和最后两位",
			args: args{
				str:   "123456789",
				start: 4,
				end:   6,
			},
			want: "123****89",
		},
		{
			name: "保留最后四位",
			args: args{
				str:   "123456789",
				start: 0,
				end:   4,
			},
			want: "****6789",
		},
		{
			name: "保留前四位",
			args: args{
				str:   "123456789",
				start: 4,
				end:   0,
			},
			want: "1234****",
		},
		{
			name: "保留前四位和最后一位",
			args: args{
				str:   "123456789",
				start: 4,
				end:   1,
			},
			want: "1234****9",
		},
		{
			name: "保留前三位和最后两位",
			args: args{
				str:   "123456789",
				start: 4,
				end:   100,
			},
			want: "123****89",
		},
		{
			name: "保留前三位和最后两位",
			args: args{
				str:   "123456789",
				start: 14,
				end:   100,
			},
			want: "123****89",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, Mark(tt.args.str, tt.args.start, tt.args.end), "Mark(%v, %v, %v)", tt.args.str, tt.args.start, tt.args.end)
		})
	}
}

func TestIsNumeric(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "纯数字",
			args: args{
				s: "123",
			},
			want: true,
		},
		{
			name: "非全数字",
			args: args{
				s: "123.123",
			},
			want: false,
		},
		{
			name: "空字符串",
			args: args{
				s: "",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, IsNumeric(tt.args.s), "IsNumeric(%v)", tt.args.s)
		})
	}
}

func TestIsChinaMobile(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "手机号+86",
			args: args{
				s: "+8613812345678",
			},
			want: true,
		},
		{
			name: "手机号非+86开头",
			args: args{
				s: "13812345678",
			},
			want: true,
		},
		{
			name: "+86非正常手机号",
			args: args{
				s: "+861234567",
			},
			want: false,
		},
		{
			name: "非正常手机号",
			args: args{
				s: "1234567",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, IsChinaMobile(tt.args.s), "IsChinaMobile(%v)", tt.args.s)
		})
	}
}
