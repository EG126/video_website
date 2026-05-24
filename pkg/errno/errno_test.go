package errno

import (
	"testing"

	"github.com/pkg/errors"
)

func TestGetCode_KnownErrors(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want int32
	}{
		{"Success", Success, 0},
		{"ParamError", ParamError, 10001},
		{"UserNotExist", UserNotExist, 10002},
		{"PasswordError", PasswordError, 10003},
		{"TokenExpired", TokenExpired, 10007},
		{"DBError", DBError, 20001},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetCode(tt.err); got != tt.want {
				t.Errorf("GetCode() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetCode_UnknownError(t *testing.T) {
	unknownErr := errors.New("unknown error")
	code := GetCode(unknownErr)
	if code != 20001 {
		t.Errorf("GetCode() for unknown error = %v, want 20001", code)
	}
}

func TestGetCode_WrappedError(t *testing.T) {
	wrappedErr := errors.Wrap(UserNotExist, "database error")
	code := GetCode(wrappedErr)
	if code != 10002 {
		t.Errorf("GetCode() for wrapped error = %v, want 10002", code)
	}
}
