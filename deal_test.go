package main

import (
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
			exp: "9eaefbcfe6bb120205f74161ed51428c",
		},
		{
			deal: &Deal{
				Product: "RTX 2080",
				Category: "Graphics Cards",
				Price: 899,
				Score: 10,
				URL: "https://example.com/store/rtx2080",
			},
			exp: "9eaefbcfe6bb120205f74161ed51428c",
		},
	}

	for _, c := range cases {
		act := c.deal.Digest()

		if act != c.exp {
			t.Errorf("Expected %#v.Digest() to return %q, %q given", c.deal, c.exp, act)
		}
	}
}
