package simplereload

import (
	"bytes"
	"net/http"
	"strconv"
	"strings"
)

type responseRecorder struct {
	http.ResponseWriter
	body *bytes.Buffer
}

func (r *responseRecorder) Write(b []byte) (int, error) {
	return r.body.Write(b)
}

func (r *responseRecorder) Flush() {
	if flusher, ok := r.ResponseWriter.(http.Flusher); ok {
		flusher.Flush()
	}
}

const (
	route     = "/simplereload"
	sseScript = `
<script>
	(function() {
		const maxRetryInterval = 1000;
		const initialRetryInterval = 100;
		let retryInterval = initialRetryInterval;
		const shouldReload = "simplereloadFlag";

		function connectEventSource() {
			const sse = new EventSource("` + route + `");
			sse.onopen = function(event) {
				console.log("* Connected to Server-Sent Events for hot reload *");
				if (sessionStorage.getItem(shouldReload)) {
					console.log("Reloading page...")
					sessionStorage.removeItem(shouldReload);
					location.reload();
				}
				// No need to reload the page if the flag is not set
				retryInterval = initialRetryInterval;
			};
			sse.onerror = function(event) {
				console.log("* Server-Sent Events connection error. Retrying in " + (retryInterval / 1000) + " seconds... *");
				sse.close();
				setTimeout(() => {
					retryInterval = Math.min(retryInterval * 2, maxRetryInterval); // Exponential backoff
					sessionStorage.setItem(shouldReload, "true");
					connectEventSource();
				}, retryInterval);
			};
		}

		connectEventSource();
	})();
</script>
`
)

// Middleware that injects a Server-Sent Events script into HTML responses, which causes a page refresh after a connection to the server is lost then regained.
func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rec := &responseRecorder{
			ResponseWriter: w,
			body:           new(bytes.Buffer),
		}

		next.ServeHTTP(rec, r)

		contentType := rec.Header().Get("Content-Type")
		if strings.HasPrefix(contentType, "text/html") {
			body := bytes.Replace(rec.body.Bytes(), []byte("<head>"), []byte("<head>"+sseScript), 1)
			rec.Header().Set("Content-Length", strconv.Itoa(len(body)))
			w.Write(body)
		} else {
			w.Write(rec.body.Bytes())
		}
	})
}
