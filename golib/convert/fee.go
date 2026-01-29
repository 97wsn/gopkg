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

// MustYuanToFen 字符串元转分，如果出错则panic
func MustYuanToFen(fee string) int {
	f, err := YuanToFen(fee)
	if err != nil {
		panic(err)
	}
	return f
}

// GetDividedAmount
// 公平算法，在保证总金额完整的前提下，将金额尽量等分，且任意两个元素之间最大差距不超过1分钱
// 如果不能均分，则金额小的在数组前面
// 注意，该方法只支持int
// 如：amount=3,totalFee = 100,则返回[33,33,34]
// 如：amount=3,totalFee = 98,则返回[32,33,33]
// 如：amount=3,totalFee = 91,则返回[30,30,31]
func GetDividedAmount[T integer](amount, count T) []T {
	if count == 1 {
		return []T{amount}
	}
	ans := make([]T, 0)
	for i := count; i > 0; i-- {
		aPart := amount / i
		amount = amount - aPart
		ans = append(ans, aPart)
	}

	return ans
}
