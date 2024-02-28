package error

import (
	"errors"
	"strings"
	"testing"
)

func TestWrapIfNotNil(t *testing.T) {
	type args struct {
		message string
		err     error
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{name: "nil error", args: args{message: "maybe an error", err: nil}, wantErr: false},
		{name: "non nil error", args: args{message: "maybe an error", err: errors.New("a bad thing")}, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var err error
			if err = WrapIfNotNil(tt.args.message, tt.args.err); (err != nil) != tt.wantErr {
				t.Errorf("WrapIfNotNil() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				if b := strings.Contains(err.Error(), tt.args.message); !b {
					t.Errorf("Non nil error did not contain the required message ")
				}
				if b := strings.Contains(err.Error(), tt.args.err.Error()); !b {
					t.Errorf("Non nil error did not contain the required error text")
				}

			}
		})
	}
}
