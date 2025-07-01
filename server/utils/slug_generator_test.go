package utils

import (
	"testing"
)

func TestGenerateSlug(t *testing.T) {
	type args struct {
		title string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "should generate slug correctly",
			args: args{
				title: "Hello World",
			},
			want: "hello-world",
		},
		{
			name: "should handle multiple spaces",
			args: args{
				title: "Hello  World",
			},
			want: "hello--world",
		},
		{
			name: "should handle special characters",
			args: args{
				title: "Hello World!",
			},
			want: "hello-world",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GenerateSlug(tt.args.title); got != tt.want {
				t.Errorf("GenerateSlug() = %v, want %v", got, tt.want)
			}
		})
	}
}
