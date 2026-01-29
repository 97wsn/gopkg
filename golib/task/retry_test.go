package task

import (
	"errors"
	"testing"
	"time"
)

func TestRetry(t *testing.T) {
	type args struct {
		num int
		fn  func() error
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{name: "error", args: args{3, func() error { return errors.New("mock error") }}, wantErr: true},
		{name: "nil", args: args{3, func() error { return nil }}, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := Retry(tt.args.num, tt.args.fn); (err != nil) != tt.wantErr {
				t.Errorf("Retry() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestRetryWithDuration(t *testing.T) {
	t.Run("", func(t *testing.T) {
		fn := func() error {
			return nil
		}
		err := RetryWithDuration(2, time.Second, func() error {
			return fn()
		})
		t.Log(err)
	})
}
