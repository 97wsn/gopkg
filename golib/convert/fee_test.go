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
		{name: "0.015", args: args{fee: 0.015}, want: 2},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ToFen(tt.args.fee); got != tt.want {
				t.Errorf("ToFen() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestYuanToFen(t *testing.T) {
	type args struct {
		fee string
	}
	tests := []struct {
		name    string
		args    args
		want    int
		wantErr bool
	}{
		{name: "-3.12", args: args{fee: "-3.12"}, want: -312, wantErr: false},
		{name: "-0.01", args: args{fee: "-0.01"}, want: -1, wantErr: false},
		{name: "-1.59", args: args{fee: "-1.59"}, want: -159, wantErr: false},
		{name: "error", args: args{fee: "error"}, want: 0, wantErr: false},
		{name: "156.7", args: args{fee: "156.7"}, want: 15670, wantErr: false},   // 15669.999999999998
		{name: "0.07", args: args{fee: "0.07"}, want: 7, wantErr: false},         // 7.000000000000001
		{name: "632.17", args: args{fee: "632.17"}, want: 63217, wantErr: false}, // 63216.99999999999
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := YuanToFen(tt.args.fee)
			if (err != nil) != tt.wantErr {
				t.Errorf("YuanToFen() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("YuanToFen() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestToYuanString(t *testing.T) {
	type args struct {
		fee int
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{name: "1999999", args: args{fee: 1999999}, want: "19999.99"},
		{name: "-100", args: args{fee: -100}, want: "-1.00"},
		{name: "399", args: args{fee: 399}, want: "3.99"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ToYuanString(tt.args.fee); got != tt.want {
				t.Errorf("ToYuanString() = %v, want %v", got, tt.want)
			}
		})
	}
}
