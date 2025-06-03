package ui

import (
	"testing"
)

func TestValidateWalletInput(t *testing.T) {
	tests := []struct {
		name    string
		input   [2]string
		wantErr bool
	}{
		{"valid", [2]string{"wallet", "123456"}, false},
		{"empty name", [2]string{"", "123456"}, true},
		{"empty password", [2]string{"wallet", ""}, true},
		{"short password", [2]string{"wallet", "123"}, true},
	}
	for _, tt := range tests {
		err := ValidateWalletInput(tt.input[0], tt.input[1])
		if (err != nil) != tt.wantErr {
			t.Errorf("%s: got error %v, wantErr %v", tt.name, err, tt.wantErr)
		}
	}
}

func TestValidateImportInput(t *testing.T) {
	err := ValidateImportInput("wallet", "123456", "word1 word2 word3 word4 word5 word6 word7 word8 word9 word10 word11 word12")
	if err != nil {
		t.Errorf("valid mnemonic: got error %v", err)
	}
	if ValidateImportInput("wallet", "123456", "") == nil {
		t.Error("empty mnemonic: expected error")
	}
	if ValidateImportInput("wallet", "123456", "word1 word2") == nil {
		t.Error("short mnemonic: expected error")
	}
}

func TestValidatePrivateKeyInput(t *testing.T) {
	validKey := "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef"
	err := ValidatePrivateKeyInput("wallet", "123456", validKey)
	if err != nil {
		t.Errorf("valid key: got error %v", err)
	}
	if ValidatePrivateKeyInput("wallet", "123456", "") == nil {
		t.Error("empty key: expected error")
	}
	if ValidatePrivateKeyInput("wallet", "123456", "xyz") == nil {
		t.Error("invalid key: expected error")
	}
}
