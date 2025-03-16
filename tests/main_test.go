package tests

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"testing"
	"time"

	"github.com/AnthonyHewins/td"
	"golang.org/x/oauth2"
	"gopkg.in/yaml.v2"
)

var c *controller

type config struct {
	Key          string `yaml:"key"`
	Secret       string `yaml:"secret"`
	RefreshToken string `yaml:"refresh_token"`

	CustomerID string `yaml:"customerID"`
	SessionID  string `yaml:"sessionID"`

	Timeout  time.Duration `yaml:"timeout"`
	LogLevel int           `yaml:"logLevel"`
}

type controller struct {
	*td.WS
	timeout          time.Duration
	equityChan       chan *td.Equity
	optionChan       chan *td.Option
	futureChan       chan *td.Future
	futureOptionChan chan *td.FutureOption
	chartEquityChan  chan *td.ChartEquity
	futureChartChan  chan *td.ChartFuture
}

func makeHandler[X any](ctx context.Context, channel chan X) func(X) {
	return func(x X) {
		select {
		case <-ctx.Done():
			return
		case channel <- x:
		}
	}
}

func (c controller) ctx() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), c.timeout)
}

func readToken() (oauth2.Token, error) {
	b, err := os.ReadFile("./token.json")
	if err != nil {
		return oauth2.Token{}, err
	}

	var t oauth2.Token
	if err := json.Unmarshal(b, &t); err != nil {
		return oauth2.Token{}, err
	}

	if time.Now().After(t.Expiry) {
		return oauth2.Token{}, fmt.Errorf("token expired")
	}

	return t, nil
}

func writeToken(t oauth2.Token) error {
	buf, err := json.Marshal(t)
	if err != nil {
		return err
	}

	return os.WriteFile("./token.json", buf, 0777)
}

func newController(ctx context.Context) (*controller, error) {
	b, err := os.ReadFile(os.Getenv("CONFIG"))
	if err != nil {
		return nil, err
	}

	var conf config
	if err := yaml.Unmarshal(b, &conf); err != nil {
		return nil, err
	}

	logger := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
		Level:     slog.Level(conf.LogLevel),
	})

	t, err := readToken()
	if err != nil {
		fmt.Println("failed reading token from cache:", err)
	}

	hc := td.New(
		td.ProdURL,
		td.AuthUrl,
		conf.Key,
		conf.Secret,
		td.WithClientLogger(logger),
		td.WithHTTPAccessToken(t.AccessToken),
	)

	if t.AccessToken == "" {
		t, err = hc.Authenticate(ctx, conf.RefreshToken)
		if err != nil {
			return nil, err
		}
	}

	if err = writeToken(t); err != nil {
		fmt.Println("failed writing token to cache:", err)
	}

	c = &controller{
		timeout:          conf.Timeout,
		equityChan:       make(chan *td.Equity, 1),
		optionChan:       make(chan *td.Option, 1),
		futureChan:       make(chan *td.Future, 1),
		futureOptionChan: make(chan *td.FutureOption, 1),
		chartEquityChan:  make(chan *td.ChartEquity, 1),
		futureChartChan:  make(chan *td.ChartFuture, 1),
	}

	c.WS, err = td.NewSocket(
		ctx,
		nil,
		hc,
		t.RefreshToken,
		td.WithLogger(logger),
		td.WithTimeout(conf.Timeout),
		td.WithEquityHandler(makeHandler(ctx, c.equityChan)),
		td.WithOptionHandler(makeHandler(ctx, c.optionChan)),
		td.WithFutureHandler(makeHandler(ctx, c.futureChan)),
		td.WithFutureOptionHandler(makeHandler(ctx, c.futureOptionChan)),
		td.WithChartEquityHandler(makeHandler(ctx, c.chartEquityChan)),
		td.WithChartFutureHandler(makeHandler(ctx, c.futureChartChan)),
	)

	if err != nil {
		return nil, err
	}

	return c, nil
}

func TestMain(m *testing.M) {
	if os.Getenv("INT_TESTS") != "1" {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()

	var err error
	c, err = newController(ctx)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	exit := 1
	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()

		c.Close(ctx)
		os.Exit(exit)
	}()

	exit = m.Run()
}
