package utils

import (
	"testing"
)

func TestEncryptPassword(t *testing.T) {
	type args struct {
		password string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "should encrypt password",
			args: args{
				password: "password",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := EncryptPassword(tt.args.password)
			if (err != nil) != tt.wantErr {
				t.Errorf("EncryptPassword() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err == nil {
				if err := ComparePasswordHash(got, tt.args.password); err != nil {
					t.Errorf("ComparePasswordHash() error = %v", err)
				}
			}
		})
	}
}
