# td

TD ameritrade API client

```shell
go get github.com/AnthonyHewins/td # requires >= go1.24
```

## Usage

Websocket:

```go
s, err := td.NewSocket(
	ctx, // parent context; canceling this will close connection
	"", // uri
	nil, // websocket connection options, if desired
	td.WSCreds{ // credentials from the user preferences endpoint
		CustomerID: "",
		SessionID:  uuid.UUID{},
		Token:      &td.Token{},
	},
	td.WithTimeout(c.Timeout),
	td.WithLogger(/* logger configured for slog, if desired */),
)
```