package tests

import (
	"testing"

	"github.com/AnthonyHewins/td"
)

func TestChartFutureSub(t *testing.T) {
	ctx, cancel := c.ctx()
	defer cancel()

	resp, err := c.AddChartFutureSubscription(ctx, &td.ChartFutureReq{
		Symbols: []string{"/ESQ25"},
		Fields:  td.ChartFutureFieldValues(),
	})

	if err != nil {
		t.Errorf("failed making chart equity subscription: %s", err)
		return
	}

	if resp.Msg == "ADD command succeeded" {
		_, err := c.UnsubEquitySubscription(ctx, "AAPL")
		if err != nil {
			t.Errorf("failed unsubbing from AAPL subscription, but was able to sub: %s", err)
		}
	}
	t.Error("asd")
}
