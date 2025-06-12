package ui

import (
	"blocowallet/internal/blockchain"
)

// networkSuggestionItem implements list.Item for the autocomplete
// This will be used in the AddNetworkComponent

type networkSuggestionItem struct {
	suggestion blockchain.NetworkSuggestion
}

func (i networkSuggestionItem) Title() string       { return i.suggestion.Name }
func (i networkSuggestionItem) Description() string { return i.suggestion.Symbol }
func (i networkSuggestionItem) FilterValue() string { return i.suggestion.Name }
