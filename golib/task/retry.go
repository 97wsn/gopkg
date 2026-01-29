package task

import "time"

// Retry 重试函数
func Retry(num int, fn func() error) error {
	var err error
	for i := 0; i < num; i++ {
		err = fn()
		if err == nil {
			break
		}
	}
	return err
}

// RetryWithDuration Retry 重试函数
func RetryWithDuration(num int, duration time.Duration, fn func() error) error {
	var err error
	for i := 0; i < num; i++ {
		err = fn()
		if err == nil {
			break
		}

		// 休息一下再重试
		time.Sleep(duration)
	}
	return err
}
