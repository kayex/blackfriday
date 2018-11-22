package main

import (
	"fmt"
	"testing"
)

func TestDeal_Digest(t *testing.T) {
	cases := []struct{
		deal *Deal
		exp string
	}{
		{
			deal: &Deal{
				Product: "RTX 2080",
				Category: "Graphics Cards",
				Price: 899,
				Score: 0,
				URL: "https://example.com/store/rtx2080",
			},
			exp: "99914b932bd37a50b983c5e7c90ae93b",
		},
		{
			deal: &Deal{
				Product: "RTX 2080",
				Category: "Graphics Cards",
				Price: 899,
				Score: 10,
				URL: "https://example.com/store/rtx2080",
			},
			exp: "99914b932bd37a50b983c5e7c90ae93b",
		},
	}

	for _, c := range cases {
		act := c.deal.Digest()
		fmt.Printf("hash: %v\n", act)

		if act != c.exp {
			t.Errorf("Expected %#v.Digest() to return %q, %q given", c.deal, c.exp, act)
		}
	}
}
