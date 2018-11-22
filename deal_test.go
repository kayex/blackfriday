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
			exp: "68ba7efd1fcea2d1b263b8d27e49bf38",
		},
		{
			deal: &Deal{
				Product: "RTX 2080",
				Category: "Graphics Cards",
				Price: 899,
				Score: 10,
				URL: "https://example.com/store/rtx2080",
			},
			exp: "68ba7efd1fcea2d1b263b8d27e49bf38",
		},
	}

	for _, c := range cases {
		act := c.deal.Digest()

		if act != c.exp {
			t.Errorf("Expected %#v.Digest() to return %q, %q given", c.deal, c.exp, act)
		}
	}
}

func TestNewDealNotification(t *testing.T) {
	vendor := "ACME Computers"

	cases := []struct{
		deal *Deal
		exp string
	}{
		{
			deal: &Deal{
				Product: "RTX 2080",
				Price: 899,
				Vendor: &vendor,
				URL: "https://example.com/store/rtx2080",
			},
			exp: "RTX 2080 *899 kr* <https://example.com/store/rtx2080|ACME Computers>",
		},
		{
			deal: &Deal{
			Product: "RTX 2080",
			Price: 899,
			URL: "https://example.com/store/rtx2080",
		},
			exp: "RTX 2080 *899 kr* <https://example.com/store/rtx2080|example.com>",
		},
		{
			deal: &Deal{
				Product: "RTX 2080",
				URL: "https://example.com/store/rtx2080",
			},
			exp: "RTX 2080 <https://example.com/store/rtx2080|example.com>",
		},
		{
			deal: &Deal{
				Product: "RTX 2080",
				URL: "this-url-looks-funky",
			},
			exp: "RTX 2080 <this-url-looks-funky|this-url-looks-funky>",
		},
		{
			deal: &Deal{
				Product: "RTX 2080",
				URL: "https://www.example.com/store/rtx2080",
			},
			exp: "RTX 2080 <https://www.example.com/store/rtx2080|example.com>",
		},
	}

	for _, c := range cases {
		act := NewDealNotification(c.deal)

		if act != c.exp {
			t.Errorf("Expected NewDealNotification(%#v) to return %q, got %q", c.deal, c.exp, act)
		}
	}
}
