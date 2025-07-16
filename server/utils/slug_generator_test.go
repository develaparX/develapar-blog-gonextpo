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
			want: "hello-world",
		},
		{
			name: "should handle special characters",
			args: args{
				title: "Hello World!",
			},
			want: "hello-world",
		},
		{
			name: "should handle empty title",
			args: args{
				title: "",
			},
			want: "",
		},
		{
			name: "should handle title with only special characters",
			args: args{
				title: "!@#$%",
			},
			want: "untitled",
		},
		{
			name: "should handle complex title",
			args: args{
				title: "How to Build a REST API with Go & Gin Framework",
			},
			want: "how-to-build-a-rest-api-with-go-gin-framework",
		},
		{
			name: "should handle title with numbers",
			args: args{
				title: "Top 10 Programming Languages in 2024",
			},
			want: "top-10-programming-languages-in-2024",
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
