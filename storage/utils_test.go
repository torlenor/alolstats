package storage

import (
	"net/url"
	"testing"
)

func Test_extractURLStringParameter(t *testing.T) {
	type args struct {
		parameters    url.Values
		parameterName string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "Test 1 - Valid parameter",
			args: args{
				parameters: url.Values{
					"param1": []string{"someValue"},
				},
				parameterName: "param1",
			},
			want:    "someValue",
			wantErr: false,
		},
		{
			name: "Test 2 - Valid parameter, but too many occurrences",
			args: args{
				parameters: url.Values{
					"param1": []string{"someValue", "someValue2"},
				},
				parameterName: "param1",
			},
			want:    "",
			wantErr: true,
		},
		{
			name: "Test 3 - Invalid/non-existing parameter",
			args: args{
				parameters: url.Values{
					"param1": []string{"someValue"},
				},
				parameterName: "param2",
			},
			want:    "",
			wantErr: true,
		},
		{
			name: "Test 4 - Invalid/non-existing parameter, because defined but no value",
			args: args{
				parameters: url.Values{
					"param1": []string{},
				},
				parameterName: "param1",
			},
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := extractURLStringParameter(tt.args.parameters, tt.args.parameterName)
			if (err != nil) != tt.wantErr {
				t.Errorf("extractURLStringParameter() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("extractURLStringParameter() = %v, want %v", got, tt.want)
			}
		})
	}
}
