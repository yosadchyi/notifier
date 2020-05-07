package notification_test

import (
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/yosadchyi/notifier/pkg/notification"
)

func TestHttpFunc(t *testing.T) {
	type given struct {
		statusCode int
	}
	type expected struct {
		nonEmptyLog  bool
		logSubstring string
	}
	cases := []struct {
		name string
		given
		expected
	}{
		{
			name: "WHEN notification sent AND 201 status code received THEN no messages logged",
			given: given{
				statusCode: http.StatusCreated,
			},
			expected: expected{
				nonEmptyLog: false,
			},
		},
		{
			name: "WHEN notification sent AND non-2xx status code received THEN no messages logged",
			given: given{
				statusCode: http.StatusNotFound,
			},
			expected: expected{
				nonEmptyLog:  true,
				logSubstring: "status code",
			},
		},
	}

	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			logout := &strings.Builder{}
			log.SetOutput(logout)
			defer log.SetOutput(os.Stderr)

			mux := http.NewServeMux()
			mux.HandleFunc("/notify", func(w http.ResponseWriter, r *http.Request) {
				w.Header().Add("Content-Type", "text/plain")
				w.WriteHeader(test.statusCode)
			})
			server := httptest.NewServer(mux)
			defer server.Close()

			url := server.URL + "/notify"
			notifyFunc := notification.HttpFunc(url)
			notifyFunc("TEST")

			if test.expected.nonEmptyLog && logout.Len() == 0 {
				t.Errorf("Expected that log will contain substring %s, but got empty log", test.expected.logSubstring)
			}
			if !test.expected.nonEmptyLog && logout.Len() > 0 {
				t.Errorf("Unexpected log message: %s", logout)
			}
		})
	}
}
