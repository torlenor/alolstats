package analyzer

import "testing"

func Test_determineRole(t *testing.T) {
	type args struct {
		lane string
		role string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Test 1 - Determine TOP",
			args: args{
				lane: "TOP",
				role: "something",
			},
			want: "TOP",
		},
		{
			name: "Test 2 - Determine TOP - mixed case",
			args: args{
				lane: "ToP",
				role: "something",
			},
			want: "TOP",
		},
		{
			name: "Test 3 - Determine MIDDLE from MID",
			args: args{
				lane: "MID",
				role: "something",
			},
			want: "MIDDLE",
		},
		{
			name: "Test 4 - Determine MIDDLE from MIDDLE",
			args: args{
				lane: "MIDDLE",
				role: "whatever",
			},
			want: "MIDDLE",
		},
		{
			name: "Test 5 - Determine JUNGLE",
			args: args{
				lane: "JUNGLE",
				role: "whatever",
			},
			want: "JUNGLE",
		},
		{
			name: "Test 6 - Determine CARRY from BOT - DUO_CARRY",
			args: args{
				lane: "BOT",
				role: "DUO_CARRY",
			},
			want: "CARRY",
		},
		{
			name: "Test 7 - Determine SUPPORT from BOTTOM - DUO_carry",
			args: args{
				lane: "BOTTOM",
				role: "DUO_carry",
			},
			want: "CARRY",
		},
		{
			name: "Test 8 - Determine SUPPORT from BOT - DUO_SUPPORT",
			args: args{
				lane: "BOT",
				role: "DUO_SUPPORT",
			},
			want: "SUPPORT",
		},
		{
			name: "Test 9 - Determine SUPPORT from BOTTOM - DUO_support",
			args: args{
				lane: "BOTTOM",
				role: "DUO_support",
			},
			want: "SUPPORT",
		},
		{
			name: "Test 10 - Determine BOTTOM_UNKNOWN",
			args: args{
				lane: "BOTTOM",
				role: "whatever",
			},
			want: "BOTTOM_UNKNOWN",
		},
		{
			name: "Test 11 - Determine UNKNOWN - everything else must be unknown",
			args: args{
				lane: "something",
				role: "whatever",
			},
			want: "UNKNOWN",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := determineRole(tt.args.lane, tt.args.role); got != tt.want {
				t.Errorf("determineRole() = %v, want %v", got, tt.want)
			}
		})
	}
}
