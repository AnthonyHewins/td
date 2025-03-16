# td

TD ameritrade API client (now schwab)

```shell
go get github.com/AnthonyHewins/td # requires >= go1.24
```

## Usage

The API as it stands is not the nicest. I can't offer the greatest dev experience, but I tried.
Currently the only supported functionality is streaming data. To stream data, you'll need the websocket.
To use the websocket:

1. Get all your credentials in order. As of right now that means client key/secret and also fetching a refresh token every 7 days **manually** (what a horrible pain)
2. Create the HTTPClient whose sole existence is to fetch the information needed to connect to the socket
3. Pass an enormous amount of information to the `NewSocket` function to handle far too much logic than should be needed for a login

Then you're ready to go

```go
hc := td.New( // HTTP client
	td.ProdURL,
	td.AuthUrl,
	"client-key",
	"client-secret",
	// td.WithClientLogger(slog.Handler),
	// td.WithHTTPAccessToken(cachedTokenIfYouHaveIt),
)

t, err = hc.Authenticate(ctx, conf.RefreshToken)
if err != nil { panic("") }

ws, err := td.NewSocket(
	ctx, // master context. When this is canceled, it shuts everything down
	nil, // websocket dial options, not required
	hc, // the HTTP client
	t.RefreshToken,
	// td.WithLogger(logger),
	// td.WithTimeout(timeout),
	// td.WithEquityHandler(),
	// td.WithOptionHandler(),
	// td.WithFutureHandler(),
	// td.WithFutureOptionHandler(),
	// td.WithChartEquityHandler(),
	// td.WithChartFutureHandler(),
)
```

## How the socket works

TODO flesh more out:
The websocket acts very much like a regular struct with methods

## TODOs

- Figure out the absymal documentation on these things:
  - FutureTradingHours field
  - FuturePriceFormat