package stringutil

import (
	"bytes"
	"encoding/json"
	"math/rand"
	"regexp"
	"strings"
	"unicode/utf8"
)

// 生成随机字符串
const randomStr = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func RandomString(length int) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = randomStr[rand.Intn(len(randomStr))]
	}
	return string(b)
}

func TitleCasedName(name string) string {
	newstr := make([]rune, 0)
	upNextChar := true

	name = strings.ToLower(name)

	for _, chr := range name {
		switch {
		case upNextChar:
			upNextChar = false
			if 'a' <= chr && chr <= 'z' {
				chr -= 'a' - 'A'
			}
		case chr == '_':
			upNextChar = true
			continue
		}

		newstr = append(newstr, chr)
	}

	return string(newstr)
}

func SnakeCasedName(name string) string {
	newstr := make([]rune, 0)
	for idx, chr := range name {
		if isUpper := 'A' <= chr && chr <= 'Z'; isUpper {
			if idx > 0 && name[idx-1] != '_' {
				newstr = append(newstr, '_')
			}
			chr -= 'A' - 'a'
		}
		newstr = append(newstr, chr)
	}

	return string(newstr)
}

func TrimBom(s string) string {
	buf := TrimBomBytes([]byte(s))
	return string(buf)
}

func TrimBomBytes(s []byte) []byte {
	buf := s
	if len(buf) > 3 {
		// 0xef, 0xbb, 0xbf 239 187 191
		if buf[0] == 239 && buf[1] == 187 && buf[2] == 191 {
			return buf[3:]
		}
	}
	return s
}

// Leftpad returns a new string of given length, filled with
func Leftpad(s string, length int, ch ...rune) string {
	c := ' '
	if len(ch) > 0 {
		c = ch[0]
	}
	l := length - utf8.RuneCountInString(s)
	if l > 0 {
		s = strings.Repeat(string(c), l) + s
	}
	return s
}

// RightPad 函数用于在字符串右侧填充指定字符至指定长度
func RightPad(str string, length int, padChar string) string {
	if len(str) >= length {
		return str
	}
	padding := strings.Repeat(padChar, length-len(str))
	return str + padding
}

func MustJsonEncode(v any) string {
	b, _ := json.Marshal(v)
	return string(b)
}

func MustJsonEncodeUnescape(v any) string {
	var buf bytes.Buffer
	jsonEncoder := json.NewEncoder(&buf)
	jsonEncoder.SetEscapeHTML(false)
	_ = jsonEncoder.Encode(v)
	return buf.String()
}

// Mark 保留指定位置字符串
func Mark(str string, start, end int) string {
	length := len(str)

	// 边界条件提前检查
	if length == 0 || start < 0 || end < 0 {
		return str
	}

	// 计算实际保留的左段和右段
	leftLen := start
	rightLen := end

	// 若保留段之和超过总长，进行合理缩减
	if leftLen+rightLen >= length {
		// 三等分保留长度，避免奇数长度无法打码
		half := length / 3
		leftLen = half
		rightLen = (length - leftLen) / 3
	}

	left := str[:leftLen]
	right := str[length-rightLen:]

	return left + strings.Repeat("*", 4) + right
}

var reIsNumber = regexp.MustCompile(`^[0-9]+$`) // 正则表达式只允许数字字符
func IsNumeric(s string) bool {
	return reIsNumber.MatchString(s)
}

// IsChinaMobile 中国大陆手机号码
var regChinaMobile = regexp.MustCompile(`^(?:\+?86)?1[3-9]\d{9}$`)

func IsChinaMobile(s string) bool {
	return regChinaMobile.MatchString(s)
}
