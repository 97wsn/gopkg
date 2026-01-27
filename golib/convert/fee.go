package convert

import (
	"fmt"
	"math"
	"strconv"
)

type signed interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64
}

type unsigned interface {
	~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr
}

type integer interface {
	signed | unsigned
}

func ToFen(fee float64) int {
	return int(math.Round(fee * 100))
}

// ToYuan 分转换为元float格式，如int(199)->float(1.99)
func ToYuan[T integer](fee T) float64 {
	return float64(fee) / 100
}

// ToYuanString "分"转换成显示为"元"的字符串，如int(199)->string(1.99)
func ToYuanString[T integer](fee T) string {
	return fmt.Sprintf("%.2f", ToYuan(fee))
}

// YuanToFen 将字符串单位元转成整型单位分
func YuanToFen(fee string) (int, error) {
	f, err := strconv.ParseFloat(fee, 64)
	if err != nil {
		return 0, nil
	}

	return ToFen(f), nil
}
