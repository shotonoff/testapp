package pow

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPOW(t *testing.T) {
	msg := []byte("goodbye world")
	testCases := []struct {
		difficulty int
		wantHash   string
		wantNonce  int
	}{
		{
			difficulty: 5,
			wantHash:   "0485f35b5c7bbbaaeb6dff0b2028a1a579cd5a690b8fecef93be5bbd69d765b3",
			wantNonce:  2,
		},
		{
			difficulty: 10,
			wantHash:   "0007da7094b16b0e156c25b61f7be358d3d68d6813dc848530ef833a42ccbd89",
			wantNonce:  3391,
		},
		{
			difficulty: 15,
			wantHash:   "0000d5fa6efc34e90c08c77b8eec60b3ff297495087816d77abec5e0f0d1090b",
			wantNonce:  38556,
		},
	}
	for i, tc := range testCases {
		t.Run(fmt.Sprintf("test-case %d", i), func(t *testing.T) {
			nonce, hash, err := Compute(msg, tc.difficulty)
			require.NoError(t, err)
			require.Equal(t, tc.wantHash, fmt.Sprintf("%x", hash))
			require.Equal(t, tc.wantNonce, nonce)
			require.True(t, Verify(msg, hash, tc.difficulty, nonce))
		})
	}
}
