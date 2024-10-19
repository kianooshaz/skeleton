package usernamesrv

import (
	"testing"

	"github.com/kianooshaz/skeleton/service/usernamesrv/stores/usernamedb"
)

func TestService_isValidUsername(t *testing.T) {
	type fields struct {
		config  Config
		queries *usernamedb.Queries
	}
	type args struct {
		value string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			name: "Valid username with allowed characters",
			fields: fields{
				config: Config{
					AllowCharacters: "ABCDEGSCX_", // Set of allowed characters
				},
			},
			args: args{
				value: "ABCDE_", // Valid username
			},
			want: true,
		},
		{
			name: "Invalid username with disallowed characters",
			fields: fields{
				config: Config{
					AllowCharacters: "ABCDEGSCX_", // Set of allowed characters
				},
			},
			args: args{
				value: "FREEDOM_", // Contains characters not in the allowed set
			},
			want: false,
		},
		{
			name: "Valid username with minimal allowed characters",
			fields: fields{
				config: Config{
					AllowCharacters: "ABCDEGSCX_", // Set of allowed characters
				},
			},
			args: args{
				value: "A_", // Only allowed characters
			},
			want: true,
		},
		{
			name: "Invalid username with lowercase characters",
			fields: fields{
				config: Config{
					AllowCharacters: "ABCDEGSCX_", // Set of allowed characters
				},
			},
			args: args{
				value: "abcde_", // Contains lowercase letters which are not allowed
			},
			want: false,
		},
		{
			name: "Empty username string",
			fields: fields{
				config: Config{
					AllowCharacters: "ABCDEGSCX_", // Set of allowed characters
				},
			},
			args: args{
				value: "", // Empty string should still be considered valid
			},
			want: true,
		},
		{
			name: "Username with special characters",
			fields: fields{
				config: Config{
					AllowCharacters: "ABCDEGSCX_", // Set of allowed characters
				},
			},
			args: args{
				value: "AB@#_", // Contains special characters not in the allowed set
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Service{
				config:  tt.fields.config,
				queries: tt.fields.queries,
			}
			if got := s.isValidUsername(tt.args.value); got != tt.want {
				t.Errorf("Service.isValidUsername() = %v, want %v", got, tt.want)
			}
		})
	}
}
