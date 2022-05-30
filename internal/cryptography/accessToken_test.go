package cryptography

import (
	"testing"
	"time"
)

func TestNewTokenJWT(t *testing.T) {
	type args struct {
		email string
		role  string
		ttl   time.Duration
	}

	var tests = []struct {
		name      string
		secretKey string
		args      args
		wantEmail string
		wantRole  string
		wantErr   bool
	}{
		{name: "Case-1: check token",
			secretKey: "hello",
			args: args{
				email: "test@email.com",
				role:  "user",
				ttl:   time.Minute * 15,
			},
			wantEmail: "test@email.com",
			wantRole:  "user",
			wantErr:   false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := &TokenJWT{
				secretKey: []byte(tt.secretKey),
			}
			got, err := o.CreateToken(tt.args.email, tt.args.role, tt.args.ttl)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateToken error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			email, role, err := o.ParseToken(got)
			if err != nil {
				t.Errorf("got = %v, ParseToken error = %v ", tt.args, err)
			}

			if email != tt.wantEmail || role != tt.wantRole {
				t.Errorf("got = %v, want email: %v, want role: %v", tt.args, tt.wantEmail, tt.wantRole)
			}
		})
	}
}
