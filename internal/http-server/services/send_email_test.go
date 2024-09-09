package services

import (
	"testing"
)

func TestEmailSending(t *testing.T) {
	cases := []struct {
		name   string
		email  string
		expErr bool
	}{
		{
			name:   "Success",
			email:  "moolikov@mail.ru",
			expErr: false,
		},
		{
			name:   "Empty email",
			email:  "",
			expErr: true,
		},
		{
			name:   "Invalid email",
			email:  "moolikov.mail",
			expErr: true,
		},
	}

	for _, tc := range cases {
		tc := tc
		err := SendMail(tc.email)
		if err != nil && !tc.expErr {
			t.Errorf("Incorrect result in test %s: expected error - %t, get error - %s.", tc.name, tc.expErr, err.Error())
		}

	}
}
