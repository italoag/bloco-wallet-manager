package blockchain

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// RPCEndpoint represents an RPC endpoint from ChainList API
type RPCEndpoint struct {
	URL          string `json:"url"`
	Tracking     string `json:"tracking"`
	IsOpenSource bool   `json:"isOpenSource"`
}

// ChainInfo represents chain information from ChainList API
type ChainInfo struct {
	ChainID        int    `json:"chainId"`
	Name           string `json:"name"`
	NativeCurrency struct {
		Name     string `json:"name"`
		Symbol   string `json:"symbol"`
		Decimals int    `json:"decimals"`
	} `json:"nativeCurrency"`
	RPC       []RPCEndpoint `json:"rpc"`
	Explorers []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"explorers"`
}

// ChainListService handles interaction with ChainList API
type ChainListService struct {
	client  *http.Client
	baseURL string
}

// NewChainListService creates a new ChainList service
func NewChainListService() *ChainListService {
	return &ChainListService{
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
		baseURL: "https://chainlist.org",
	}
}

// GetChainInfo fetches chain information by chain ID
func (s *ChainListService) GetChainInfo(chainID int) (*ChainInfo, error) {
	url := fmt.Sprintf("%s/rpcs.json", s.baseURL)

	resp, err := s.client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch chain list: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var chains []ChainInfo
	if err := json.Unmarshal(body, &chains); err != nil {
		return nil, fmt.Errorf("failed to parse JSON response: %w", err)
	}

	// Find chain by ID
	for _, chain := range chains {
		if chain.ChainID == chainID {
			return &chain, nil
		}
	}

	return nil, fmt.Errorf("chain with ID %d not found", chainID)
}

// ValidateRPCEndpoint checks if an RPC endpoint is accessible
func (s *ChainListService) ValidateRPCEndpoint(rpcURL string) error {
	if rpcURL == "" {
		return fmt.Errorf("RPC URL cannot be empty")
	}

	// Create a simple JSON-RPC request to check if the endpoint is alive
	reqBody := `{"jsonrpc":"2.0","method":"eth_chainId","params":[],"id":1}`

	resp, err := s.client.Post(rpcURL, "application/json",
		strings.NewReader(reqBody))
	if err != nil {
		return fmt.Errorf("RPC endpoint is not accessible: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("RPC endpoint returned status: %d", resp.StatusCode)
	}

	return nil
}

// GetChainIDFromRPC attempts to get chain ID from RPC endpoint
func (s *ChainListService) GetChainIDFromRPC(rpcURL string) (int, error) {
	reqBody := `{"jsonrpc":"2.0","method":"eth_chainId","params":[],"id":1}`

	resp, err := s.client.Post(rpcURL, "application/json",
		strings.NewReader(reqBody))
	if err != nil {
		return 0, fmt.Errorf("failed to call RPC: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("RPC call failed with status: %d", resp.StatusCode)
	}

	var result struct {
		Result string `json:"result"`
		Error  *struct {
			Code    int    `json:"code"`
			Message string `json:"message"`
		} `json:"error"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return 0, fmt.Errorf("failed to decode response: %w", err)
	}

	if result.Error != nil {
		return 0, fmt.Errorf("RPC error: %s", result.Error.Message)
	}

	// Convert hex chain ID to int
	chainID, err := strconv.ParseInt(result.Result, 0, 64)
	if err != nil {
		return 0, fmt.Errorf("failed to parse chain ID: %w", err)
	}

	return int(chainID), nil
}
