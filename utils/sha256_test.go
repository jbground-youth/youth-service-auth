package utils

import (
	"log"
	"testing"
)

func Test_encryptSHA256(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Test Case 1",
			args: args{
				s: "jsjeong",
			},
		},
		{
			name: "Test Case 2",
			args: args{
				s: "test",
			},
		},
		{
			name: "Test Case 3",
			args: args{
				s: "testpassword2",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			log.Print(encryptSHA256(tt.args.s))

		})
	}
}
