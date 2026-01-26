package convert

import "testing"

func TestToFen(t *testing.T) {
	type args struct {
		fee float64
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{name: "0.01", args: args{fee: 0.01}, want: 1},
		{name: "1.4", args: args{fee: 1.4}, want: 140},
		{name: "-0.01", args: args{fee: -0.01}, want: -1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ToFen(tt.args.fee); got != tt.want {
				t.Errorf("ToFen() = %v, want %v", got, tt.want)
			}
		})
	}
}
