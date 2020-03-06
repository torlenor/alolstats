package utils

import (
	"sort"
	"testing"
)

func TestHashInt(t *testing.T) {
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
			want: "1_123_555_-234",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := make([]int, len(tt.args.i))
			copy(b, tt.args.i)
			if got := HashInt(tt.args.i); got != tt.want {
				t.Errorf("HashInt() = %v, want %v", got, tt.want)
			}
			for j, val := range b {
				if val != tt.args.i[j] {
					t.Error("Slice is not the same as before (unsorted)")
					return
				}
			}
		})
	}
}

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
			b := make([]int, len(tt.args.i))
			copy(b, tt.args.i)
			if got := HashSortedInt(tt.args.i); got != tt.want {
				t.Errorf("HashSortedInt() = %v, want %v", got, tt.want)
			}
			sort.Ints(b)
			for j, val := range b {
				if val != tt.args.i[j] {
					t.Error("Slice is not sorted")
					return
				}
			}

		})
	}
}
