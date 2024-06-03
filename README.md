# simplereload

`simplereload` is a small Go module that enables reloading the webpage after web server restart. It is a great companion to [wgo](https://github.com/bokwoon95/wgo), which reloads the server upon file change. So with `wgo` and `simplereload` together, whenever you save a file in your project your Go web server will restart and then the browser page will refresh automatically.

Note: It doesn't seem to work in Firefox (SSE connection "onopen" isn't called). If someone can get it working let me know. There is an issue open for it.

## Features
- Reload middleware for Go web servers
- Based upon standard `net/http`
- Uses Server-Sent Events (SSE) instead of a websocket
- Just 2 small files with no dependencies

## Installation
```go
go get github.com/jazoom/simplereload
```

## Example usage
### ...with `net/http`

```go
package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/jazoom/simplereload"
)

func main() {
	isDev := flag.Bool("dev", false, "Whether to run in development mode")
	flag.Parse()

	mux := http.NewServeMux()

	handler := http.Handler(mux)
	if *isDev {
		mux.Handle("/simplereload", http.HandlerFunc(simplereload.Handler))
		handler = simplereload.Middleware(mux)
	}

	log.Fatal(http.ListenAndServe(":2555", handler))
}
```

### ...with `chi`

```go
package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/jazoom/simplereload"
)

func main() {
	isDev := flag.Bool("dev", false, "Whether to run in development mode")
	flag.Parse()

	r := chi.NewRouter()

	if *isDev {
		r.Use(simplereload.Middleware)
		r.Get("/simplereload", simplereload.Handler)
	}

	log.Fatal(http.ListenAndServe(":2555", r))
}
```

### ...with `echo`
```go
package main

import (
	"flag"
	"net/http"

	"github.com/jazoom/simplereload"
	"github.com/labstack/echo/v4"
)

func main() {
	isDev := flag.Bool("dev", false, "Whether to run in development mode")
	flag.Parse()

	e := echo.New()

	if *isDev {
		e.Use(echo.WrapMiddleware(simplereload.Middleware))
		e.GET("/simplereload", echo.WrapHandler(http.HandlerFunc(simplereload.Handler)))
	}


	e.Logger.Fatal(e.Start(":2555"))
}
```

## Notes
- The middleware injects a `<script>` into the `<head>` of HTML responses served by the web server. Thus, it won't work if for some reason the `<head>` tag is missing. Also, you may experience weird behaviour if you have `"<head>"` written in your response before the actual `<head>` tag, since a simple search-and-replace is used for the injection. But that would be an unusual scenario, and you can probably get away with it by HTML Entity Encoding your string as `"&lt;head&gt;"`.
