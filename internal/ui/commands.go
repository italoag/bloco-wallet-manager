package ui

import (
	"context"
	"fmt"
	"strings"

	"blocowallet/pkg/config"

	tea "github.com/charmbracelet/bubbletea"
)

// createFriendlyErrorMsg creates user-friendly error messages
func createFriendlyErrorMsg(operation string, err error) errorMsg {
	errStr := err.Error()
	errLower := strings.ToLower(errStr)

	// Check for common network-related errors
	if strings.Contains(errLower, "connection") ||
		strings.Contains(errLower, "timeout") ||
		strings.Contains(errLower, "network") ||
		strings.Contains(errLower, "no such host") ||
		strings.Contains(errLower, "connection refused") ||
		strings.Contains(errLower, "dial tcp") {
		return errorMsg(fmt.Sprintf("🌐 Connection Problem: Unable to %s. Please check your internet connection and try again.", operation))
	}

	// Check for RPC errors
	if strings.Contains(errLower, "rpc") || strings.Contains(errLower, "endpoint") {
		return errorMsg(fmt.Sprintf("🔌 RPC Error: Unable to %s. The blockchain network may be temporarily unavailable.", operation))
	}

	// Check for authentication errors
	if strings.Contains(errLower, "unauthorized") || strings.Contains(errLower, "forbidden") {
		return errorMsg(fmt.Sprintf("🔐 Authentication Error: Unable to %s. Please check your credentials.", operation))
	}

	// Check for rate limiting
	if strings.Contains(errLower, "rate limit") || strings.Contains(errLower, "too many requests") {
		return errorMsg(fmt.Sprintf("⏱️ Rate Limited: Too many requests. Please wait a moment before trying to %s again.", operation))
	}

	// Generic error with context
	return errorMsg(fmt.Sprintf("❌ Error: Unable to %s. %s", operation, errStr))
}

// loadWalletsCmd loads all wallets from the service
func (m Model) loadWalletsCmd() tea.Cmd {
	return tea.Cmd(func() tea.Msg {
		wallets, err := m.walletService.List(context.Background())
		if err != nil {
			return createFriendlyErrorMsg("load wallets", err)
		}
		return walletsLoadedMsg(wallets)
	})
}

// getBalanceCmd gets the balance for a given address
func (m Model) getBalanceCmd(address string) tea.Cmd {
	return tea.Cmd(func() tea.Msg {
		balance, err := m.walletService.GetBalance(context.Background(), address)
		if err != nil {
			return createFriendlyErrorMsg("fetch wallet balance", err)
		}
		return balanceLoadedMsg(balance)
	})
}

// getMultiBalanceCmd gets the balance for a given address across all active networks
func (m Model) getMultiBalanceCmd(address string) tea.Cmd {
	return tea.Cmd(func() tea.Msg {
		multiBalance, err := m.walletService.GetMultiNetworkBalance(context.Background(), address)
		if err != nil {
			return createFriendlyErrorMsg("fetch wallet balances", err)
		}
		return multiBalanceLoadedMsg(multiBalance)
	})
}

// createWalletCmd creates a new wallet with the given name and password
func (m Model) createWalletCmd(name, password string) tea.Cmd {
	return tea.Cmd(func() tea.Msg {
		_, err := m.walletService.CreateWalletWithMnemonic(context.Background(), name, password)
		if err != nil {
			return errorMsg(err.Error())
		}
		return walletCreatedMsg{}
	})
}

// importWalletCmd imports a wallet from a mnemonic phrase
func (m Model) importWalletCmd(name, password, mnemonic string) tea.Cmd {
	return tea.Cmd(func() tea.Msg {
		_, err := m.walletService.ImportWalletFromMnemonic(context.Background(), name, mnemonic, password)
		if err != nil {
			return errorMsg(err.Error())
		}
		return walletCreatedMsg{}
	})
}

// importWalletFromPrivateKeyCmd imports a wallet from a private key
func (m Model) importWalletFromPrivateKeyCmd(name, password, privateKey string) tea.Cmd {
	return tea.Cmd(func() tea.Msg {
		_, err := m.walletService.ImportWalletFromPrivateKey(context.Background(), name, privateKey, password)
		if err != nil {
			return errorMsg(err.Error())
		}
		return walletCreatedMsg{}
	})
}

// Message types for network operations
type networkErrorMsg string

// addNetworkCmd adds a custom network using ChainList API with retry
func (m Model) addNetworkCmd(name, chainIDStr, rpcEndpoint string) tea.Cmd {
	return tea.Cmd(func() tea.Msg {
		// Parse chain ID
		var chainID int
		if _, err := fmt.Sscanf(chainIDStr, "%d", &chainID); err != nil {
			return networkErrorMsg("Invalid chain ID: " + err.Error())
		}

		// Get chain info with retry mechanism
		chainInfo, workingRPC, err := m.chainListService.GetChainInfoWithRetry(chainID)
		if err != nil {
			friendlyErr := createFriendlyErrorMsg("fetch chain information", err)
			return networkErrorMsg(string(friendlyErr))
		}

		// Use the working RPC from chainlist if no custom RPC provided
		finalRPC := rpcEndpoint
		if finalRPC == "" {
			finalRPC = workingRPC
		}

		// Create network configuration
		network := config.Network{
			Name:        name,
			RPCEndpoint: finalRPC,
			ChainID:     int64(chainInfo.ChainID),
			Symbol:      chainInfo.NativeCurrency.Symbol,
			Explorer:    "",
			IsActive:    false,
			IsCustom:    true,
		}

		// Set explorer if available
		if len(chainInfo.Explorers) > 0 {
			network.Explorer = chainInfo.Explorers[0].URL
		}

		// Generate a unique key
		key := fmt.Sprintf("custom_%d", chainID)

		return networkAddedMsg{
			key:     key,
			network: network,
		}
	})
}

// searchNetworksCmd searches for network suggestions by name
func (m Model) searchNetworksCmd(query string) tea.Cmd {
	return tea.Cmd(func() tea.Msg {
		suggestions, err := m.chainListService.SearchNetworksByName(query)
		if err != nil {
			return createFriendlyErrorMsg("search networks", err)
		}
		return networkSuggestionsMsg(suggestions)
	})
}

// loadChainInfoByIDCmd loads chain info using the new retry mechanism
func (m Model) loadChainInfoByIDCmd(chainID int) tea.Cmd {
	return tea.Cmd(func() tea.Msg {
		chainInfo, rpcURL, err := m.chainListService.GetChainInfoWithRetry(chainID)
		if err != nil {
			return createFriendlyErrorMsg("load chain information", err)
		}
		return chainInfoLoadedMsg{
			chainInfo: chainInfo,
			rpcURL:    rpcURL,
		}
	})
}
