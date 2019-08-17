package utils

import "testing"

func TestHashSortedInt(t *testing.T) {
	type args struct {
		i []int
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Test 1 - feed some ints",
			args: args{
				i: []int{1, 123, 555, -234},
			},
			want: "-234_1_123_555",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := HashSortedInt(tt.args.i); got != tt.want {
				t.Errorf("HashSortedInt() = %v, want %v", got, tt.want)
			}
		})
	}
}
