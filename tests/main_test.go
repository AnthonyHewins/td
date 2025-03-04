package tests

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"testing"
	"time"

	"github.com/AnthonyHewins/td"
	"github.com/google/uuid"
	"gopkg.in/yaml.v2"
)

type config struct {
	Timeout  time.Duration
	LogLevel string
}

type controller struct {
	*td.WS
	timeout time.Duration
}

func (c controller) ctx() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), c.timeout)
}

func newController() (*controller, error) {
	b, err := os.ReadFile(os.Getenv("CONFIG"))
	if err != nil {
		return nil, err
	}

	var c config
	if err := yaml.Unmarshal(b, &c); err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), c.Timeout)
	defer cancel()

	s, err := td.NewSocket(
		ctx,
		"",
		nil,
		td.WSCreds{
			CustomerID: "",
			SessionID:  uuid.UUID{},
			Token:      &td.Token{},
		},
		td.WithTimeout(c.Timeout),
		td.WithLogger(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			AddSource: true,
			Level:     slog.LevelDebug,
		})),
	)

	if err != nil {
		return nil, err
	}

	return &controller{WS: s, timeout: c.Timeout}, nil
}

func TestMain(m *testing.M) {
	if os.Getenv("INT_TESTS") != "1" {
		return
	}

	c, err := newController()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	exit := 1
	defer func() {
		c.Close()
		os.Exit(exit)
	}()

	exit = m.Run()
}
