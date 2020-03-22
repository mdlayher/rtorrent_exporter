package rtorrentexporter

import (
	"bytes"
	"context"
	"net"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestSCGITransport(t *testing.T) {
	transport := &SCGITransport{}

	t.Run("dial", func(t *testing.T) {
		var dialNetwork string
		var dialAddr string
		transport.DialContext = func(_ context.Context, network, addr string) (net.Conn, error) {
			dialNetwork = network
			dialAddr = addr
			return nil, nil
		}

		t.Run("withHost", func(t *testing.T) {
			requestURL, _ := url.Parse("scgi://host:8080")
			_, err := transport.dial(context.TODO(), requestURL)
			if err != nil {
				t.Fatalf("dial returned: %s", err)
			}

			assertEqual(t, "tcp", dialNetwork)
			assertEqual(t, "host:8080", dialAddr)
		})

		t.Run("withoutHost", func(t *testing.T) {
			requestURL, _ := url.Parse("scgi:///socket")
			_, err := transport.dial(context.TODO(), requestURL)
			if err != nil {
				t.Fatalf("dial returned: %s", err)
			}

			assertEqual(t, "unix", dialNetwork)
			assertEqual(t, "/socket", dialAddr)
		})

		t.Run("withBadScheme", func(t *testing.T) {
			requestURL, _ := url.Parse("http://")
			_, err := transport.dial(context.TODO(), requestURL)
			if err == nil {
				t.Fatalf("expected dial to return an error")
			}

			assertEqual(t, `unsupported protocol scheme "http"`, err.Error())
		})
	})

	t.Run("writeRequest", func(t *testing.T) {
		channel := make(chan []byte)
		server := func(conn net.Conn) {
			buffer := make([]byte, 128)
			n, _ := conn.Read(buffer)
			conn.Close()
			channel <- buffer[:n]
		}

		t.Run("withBody", func(t *testing.T) {
			input, output := net.Pipe()
			defer input.Close()
			go server(output)

			body := `<?xml version="1.0" encoding="UTF-8"?>`
			request := httptest.NewRequest("POST", "scgi://", bytes.NewBufferString(body))
			if err := transport.writeRequest(request, input); err != nil {
				t.Fatalf("writeRequest returned: %s", err)
			}

			assertEqual(t, "25:CONTENT_LENGTH\x0038\x00SCGI\x001\x00,"+body, string(<-channel))
		})

		t.Run("withoutBody", func(t *testing.T) {
			input, output := net.Pipe()
			defer input.Close()
			go server(output)

			request := httptest.NewRequest("POST", "scgi://", nil)
			if err := transport.writeRequest(request, input); err != nil {
				t.Fatalf("writeRequest returned: %s", err)
			}

			assertEqual(t, "24:CONTENT_LENGTH\x000\x00SCGI\x001\x00,", string(<-channel))
		})
	})

	t.Run("readResponse", func(t *testing.T) {
		server := func(conn net.Conn, data []byte) {
			conn.Write(data)
			conn.Close()
		}

		t.Run("withStatus", func(t *testing.T) {
			input, output := net.Pipe()
			defer output.Close()
			go server(input, []byte("Status: 200 OK\r\nContent-Length: 0\r\n\r\n"))

			request := httptest.NewRequest("POST", "scgi://", nil)
			response, err := transport.readResponse(request, output)
			if err != nil {
				t.Fatalf("readResponse returned: %s", err)
			}

			assertEqual(t, "200 OK", response.Status)
		})

		t.Run("withoutStatus", func(t *testing.T) {
			input, output := net.Pipe()
			defer output.Close()
			go server(input, []byte("HTTP/1.1 200 OK\r\n"))

			request := httptest.NewRequest("POST", "scgi://", nil)
			_, err := transport.readResponse(request, output)
			if err == nil {
				t.Fatalf("expected readResponse to return an error")
			}

			assertEqual(t, `expected "Status:" header, received "HTTP/1.1 200 OK\r\n"`, err.Error())
		})
	})
}

func assertEqual(t *testing.T, expected, actual interface{}) {
	if expected != actual {
		t.Fatalf("%q != %q", expected, actual)
	}
}
