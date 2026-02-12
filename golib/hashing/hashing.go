package hashing

import (
	"bufio"
	"bytes"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"hash"
	"io"
	"net/url"
	"sort"
	"strings"
)

func Encode(v url.Values) string {
	if v == nil {
		return ""
	}
	var buf bytes.Buffer
	keys := make([]string, 0, len(v))
	for k := range v {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		vs := v[k]
		prefix := k + "="
		for _, v := range vs {
			if buf.Len() > 0 {
				buf.WriteByte('&')
			}
			buf.WriteString(prefix)
			buf.WriteString(v)
		}
	}
	return buf.String()
}

func Sign(values url.Values, appSecret string, hash ...string) string {
	hashName := "md5"
	if len(hash) > 0 {
		hashName = strings.ToLower(hash[0])
	}

	for key := range values {
		if values.Get(key) == "" || key == "sign" {
			values.Del(key)
		}
	}

	str := Encode(values) + appSecret

	switch hashName {
	case "sha1":
		return Sha1(str)
	case "sha256":
		return Sha256(str)
	case "sha512":
		return Sha512(str)
	default:
		return Md5(str)
	}
}

func Md5(text string) string {
	algorithm := md5.New()
	return stringHasher(algorithm, text)
}

func Md5File(file io.Reader) string {
	r := bufio.NewReader(file)
	md5h := md5.New()
	_, err := io.Copy(md5h, r)
	if err != nil {
		return ""
	}
	return hex.EncodeToString(md5h.Sum(nil))
}

// Sha1 hashes using sha1 algorithm
func Sha1(text string) string {
	algorithm := sha1.New()
	return stringHasher(algorithm, text)
}

// Sha256 hashes using sha256 algorithm
func Sha256(text string) string {
	algorithm := sha256.New()
	return stringHasher(algorithm, text)
}

// Sha512 hashes using sha512 algorithm
func Sha512(text string) string {
	algorithm := sha512.New()
	return stringHasher(algorithm, text)
}

func stringHasher(algorithm hash.Hash, text string) string {
	algorithm.Write([]byte(text))
	return hex.EncodeToString(algorithm.Sum(nil))
}

func GetSign(values url.Values, appSecret string) string {
	return Sign(values, appSecret)
}

func MakeSign(data any, apiKey string) string {
	values := make(url.Values)
	if params, ok := data.(map[string]any); ok {
		for k, v := range params {
			switch v := v.(type) {
			case float32:
				values.Add(k, fmt.Sprint(int64(v)))
			case float64:
				values.Add(k, fmt.Sprint(int64(v)))
			default:
				values.Add(k, fmt.Sprint(v))
			}
		}
	}

	if params, ok := data.(map[string]string); ok {
		for k, v := range params {
			values.Add(k, v)
		}
	}

	return Sign(values, apiKey)
}
