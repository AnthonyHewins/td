# td

TD ameritrade API client (now schwab)

```shell
go get github.com/AnthonyHewins/td # requires >= go1.24
```

## Disclaimers

Schwab's API has many things it does that contradict the documentation.
There are many times I had to manually write code to spec and run it only to see that something
was very wrong. As it stands, everything is correct rather than matching what their docs say.
If changes come along server side, this client may be incorrect in certain things

Known areas:
- `ChartFuture` fields are completely out of order

## Usage

The API as it stands is not the nicest. I can't offer the greatest dev experience, but I tried.
Currently the only supported functionality is streaming data. To stream data, you'll need the websocket

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
	td.WithErrHandler(func (err error) {
		// if there is a disconnect, you'll get an error wrapped in net.ErrClosed
		if errors.Is(err, net.ErrClosed) {
			// trigger reconnect, log, whatever
		}

		panic(err)
	})
)
```

## How the socket works

### Under the hood

The websocket acts like a regular client. Under the hood though, it does several things:

- Starts a goroutine to handle pings at regular intervals
- Starts a goroutine to handle constant reads from the websocket, so you're always listening for the next message; when a message is received, it gets pushed to the channel in the below goroutine
- Starts a goroutine whose sole job is to deserialize the message received from the above goroutine and then route it to the correct spot since messages come in out of order

### Calling methods on the socket

Calling methods on the socket works just like any regular code. Call the method, get the response or an error

- Responses are routed from the code mentioned above in the goroutine that routes the response back to you
- Errors that occur from a method call will not propagate to the error handler you pass in

### Observability

These 3 goroutines can witness lots of errors, so it's important that if you want good visibility that you at least use `WihtErrHandler` that routes errors to a handler you make. In addition you can handle server `pong` messages with another handler this package offers

### Disconnects

When there's a disconnect event that the server initiated, that happens in the goroutines that manage the reader of the socket. To propagate that event to your code, it gets passed to the error handler in `WithErrHandler`. Whenever there's a disconnect event, the error sent will be wrapped with `net.ErrClosed` which you can test for using `errors.Is(err, net.ErrClosed)` in standard Go fashion

If you cancel the context passed during creation, no error will be sent because this was initiated by you

## TODOs

- Figure out the absymal documentation on these things:
  - FutureTradingHours field
  - FuturePriceFormat