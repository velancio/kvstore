package util

import "testing"

func TestValidateKvPair(t *testing.T) {
	tests := []struct {
		name       string
		key        string
		value      string
		wantErr    bool
		wantErrMsg string
	}{
		{
			name:    "valid key-value pair",
			key:     "test-key",
			value:   "test-value",
			wantErr: false,
		},
		{
			name:       "empty key",
			key:        "",
			value:      "test-value",
			wantErr:    true,
			wantErrMsg: "key cannot be empty",
		},
		{
			name:       "empty value",
			key:        "test-key",
			value:      "",
			wantErr:    true,
			wantErrMsg: "value cannot be empty",
		},
		{
			name:       "invalid key characters",
			key:        "test-key!",
			value:      "test-value",
			wantErr:    true,
			wantErrMsg: "key contains invalid characters",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateKvPair(tt.key, tt.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateKvPair() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && err.Error() != tt.wantErrMsg {
				t.Errorf("ValidateKvPair() error message = %v, want %v", err, tt.wantErrMsg)
			}
		})
	}
}
