package ui

import (
	"testing"
)

func TestValidateNetworkInput(t *testing.T) {
	tests := []struct {
		name, chainID, rpc string
		wantErr            bool
	}{
		{"Ethereum", "1", "https://mainnet.infura.io", false},
		{"", "1", "https://mainnet.infura.io", true},
		{"Ethereum", "", "https://mainnet.infura.io", true},
		{"Ethereum", "1", "", true},
		{"Ethereum", "abc", "https://mainnet.infura.io", true},
		{"Ethereum", "1", "mainnet.infura.io", true},
	}
	for _, tt := range tests {
		err := ValidateNetworkInput(tt.name, tt.chainID, tt.rpc)
		if (err != nil) != tt.wantErr {
			t.Errorf("ValidateNetworkInput(%q, %q, %q): got error %v, wantErr %v", tt.name, tt.chainID, tt.rpc, err, tt.wantErr)
		}
	}
}
