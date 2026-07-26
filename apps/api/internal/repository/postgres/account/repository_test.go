package account

import (
	"slices"
	"testing"
)

func TestAccountSearchTokensKeepsDistinctiveAccountWord(t *testing.T) {
	tokens := accountSearchTokens("Prospect RS Kariadi, memiliki stage desire, dengan target deal 50 juta")
	if !slices.Contains(tokens, "kariadi") {
		t.Fatalf("expected account search tokens to include kariadi, got %#v", tokens)
	}
}

func TestAccountSearchTokensDropsNoisyCRMWords(t *testing.T) {
	tokens := accountSearchTokens("buat deal stage desire untuk RSUP Dr Kariadi")
	for _, noisy := range []string{"buat", "deal", "stage", "rsup", "dr"} {
		if slices.Contains(tokens, noisy) {
			t.Fatalf("expected noisy token %q to be dropped, got %#v", noisy, tokens)
		}
	}
}
