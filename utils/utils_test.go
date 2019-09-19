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
		{
			name: "Test 6 - Invalid version, negative number",
			args: args{
				version: "-9.1.1",
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

func TestSplitNumericMatchVersion(t *testing.T) {
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
				version: "9.1.1.123",
			},
			want:    []uint32{9, 1, 1, 123},
			wantErr: false,
		},
		{
			name: "Test 2 - Valid version",
			args: args{
				version: "4.18.1.2",
			},
			want:    []uint32{4, 18, 1, 2},
			wantErr: false,
		},
		{
			name: "Test 3a - Invalid version, not enough digits",
			args: args{
				version: "4.18",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Test 3b - Invalid version, not enough digits",
			args: args{
				version: "4.18.2",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Test 4 - Invalid version, contains chars",
			args: args{
				version: "4.18a.2.123",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Test 5 - Invalid version, too long",
			args: args{
				version: "4.18.2.123.555",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Test 6 - Invalid version, negative values",
			args: args{
				version: "4.-18.1.2",
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := SplitNumericMatchVersion(tt.args.version)
			if (err != nil) != tt.wantErr {
				t.Errorf("SplitNumericMatchVersion() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SplitNumericMatchVersion() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGenerateStatusResponse(t *testing.T) {
	type args struct {
		statusCode uint16
		statusText string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Test 1 - Status code 200, some text",
			args: args{
				statusCode: 200,
				statusText: "some text",
			},
			want: `{"error": { "code": 200, "message": "some text" } }`,
		},
		{
			name: "Test 2 - Status code 404, some other text",
			args: args{
				statusCode: 404,
				statusText: "some other text",
			},
			want: `{"error": { "code": 404, "message": "some other text" } }`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GenerateStatusResponse(tt.args.statusCode, tt.args.statusText); got != tt.want {
				t.Errorf("GenerateStatusResponse() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCleanUpSummonerName(t *testing.T) {
	type args struct {
		summonerName string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Test 1 - Cleanup summoner name - uppercase, spaces, unterlines and dashes",
			args: args{
				summonerName: "Test Su_-mmOner",
			},
			want: `testsummoner`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CleanUpSummonerName(tt.args.summonerName); got != tt.want {
				t.Errorf("CleanUpSummonerName() = %v, want %v", got, tt.want)
			}
		})
	}
}
