package rtorrentexporter

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"strings"
)

const headersFormat = "CONTENT_LENGTH\x00%d\x00SCGI\x001\x00"
const statusHeaderPrefix = "Status:"

type body struct {
	io.ReadCloser
	conn net.Conn
}

func (b *body) Close() error {
	defer b.conn.Close()
	return b.ReadCloser.Close()
}

// SCGITransport implements RoundTripper for the 'scgi' protocol.
//
// This is a minimal, rTorrent specific implementation.
type SCGITransport struct {
	// DialContext specifies the dial function for creating connections.
	// If DialContext is nil then the transport dials using net.Dialer.
	DialContext func(ctx context.Context, network, addr string) (net.Conn, error)
}

// RoundTrip implements the net/http.RoundTripper interface.
//
// A new connection is established for each handled request.
func (t SCGITransport) RoundTrip(req *http.Request) (*http.Response, error) {
	conn, err := t.dial(req.Context(), req.URL)
	if err != nil {
		return nil, err
	}

	if err = t.writeRequest(req, conn); err != nil {
		conn.Close()
		return nil, err
	}

	resp, err := t.readResponse(req, conn)
	if err != nil {
		conn.Close()
		return nil, err
	}
	// Close the connection after body is consumed
	resp.Body = &body{
		ReadCloser: resp.Body,
		conn:       conn,
	}
	return resp, nil
}

var zeroDialer net.Dialer

// dial creates a TCP or UNIX socket connection, depending on the URL format.
func (t SCGITransport) dial(ctx context.Context, url *url.URL) (net.Conn, error) {
	if url.Scheme != "scgi" {
		return nil, fmt.Errorf("unsupported protocol scheme %q", url.Scheme)
	}

	var network, addr string
	if url.Host != "" {
		network = "tcp"
		addr = url.Host
	} else {
		network = "unix"
		addr = url.Path
	}

	if t.DialContext != nil {
		return t.DialContext(ctx, network, addr)
	}
	return zeroDialer.DialContext(ctx, network, addr)
}

// writeRequest formats and writes the SCGI request.
func (t SCGITransport) writeRequest(req *http.Request, conn net.Conn) error {
	writer := bufio.NewWriter(conn)
	headers := fmt.Sprintf(headersFormat, req.ContentLength)

	if _, err := writer.WriteString(fmt.Sprintf("%d:%s,", len(headers), headers)); err != nil {
		return err
	}

	if req.Body != nil {
		defer req.Body.Close()
		if _, err := io.Copy(writer, req.Body); err != nil {
			return err
		}
	}

	return writer.Flush()
}

// readResponse reads and parses the SCGI response.
func (t SCGITransport) readResponse(req *http.Request, conn net.Conn) (*http.Response, error) {
	reader := bufio.NewReader(conn)
	status, err := reader.ReadString('\n')
	if err != nil {
		return nil, err
	}

	if !strings.HasPrefix(status, statusHeaderPrefix) {
		return nil, fmt.Errorf("expected %q header, received %q", statusHeaderPrefix, status)
	}
	status = strings.Replace(status, statusHeaderPrefix, req.Proto, 1)
	reader = bufio.NewReader(io.MultiReader(strings.NewReader(status), reader))
	return http.ReadResponse(reader, req)
}
