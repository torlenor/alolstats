package utils

import (
	"reflect"
	"testing"
)

func TestSplitNumericVersion(t *testing.T) {
	type args struct {
		version string
	}
	tests := []struct {
		name    string
		args    args
		want    []uint32
		wantErr bool
	}{
		{
			name: "Test 1 - Valid version",
			args: args{
				version: "9.1.1",
			},
			want:    []uint32{9, 1, 1},
			wantErr: false,
		},
		{
			name: "Test 2 - Valid version",
			args: args{
				version: "4.18.1",
			},
			want:    []uint32{4, 18, 1},
			wantErr: false,
		},
		{
			name: "Test 3 - Invalid version, not enough digits",
			args: args{
				version: "4.18",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Test 4 - Invalid version, contains chars",
			args: args{
				version: "4.18a.2",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Test 5 - Invalid version, too long",
			args: args{
				version: "4.18.2.3",
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := SplitNumericVersion(tt.args.version)
			if (err != nil) != tt.wantErr {
				t.Errorf("SplitNumericVersion() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SplitNumericVersion() = %v, want %v", got, tt.want)
			}
		})
	}
}
