package tests

import (
	"testing"

	"github.com/AnthonyHewins/td"
)

func TestEquitySub(t *testing.T) {
	ctx, cancel := c.ctx()
	defer cancel()

	resp, err := c.SetEquitySubscription(ctx, &td.EquityReq{
		Symbols: []string{"AAPL"},
		Fields:  []td.EquityField{td.EquityFieldHighPrice},
	})

	if err != nil {
		t.Errorf("failed making chart equity subscription: %s", err)
	}

	if resp.Msg == "SUBS command succeeded" {
		_, err := c.UnsubEquitySubscription(ctx, "AAPL")
		if err != nil {
			t.Errorf("failed unsubbing from AAPL subscription, but was able to sub: %s", err)
		}
	}
}
