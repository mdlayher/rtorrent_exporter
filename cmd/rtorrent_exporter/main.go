// Command rtorrent_exporter provides a Prometheus exporter for rTorrent.
package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/mdlayher/rtorrent"
	"github.com/mdlayher/rtorrent_exporter"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	telemetryAddr = flag.String("telemetry.addr", ":9135", "host:port for rTorrent exporter")
	metricsPath   = flag.String("telemetry.path", "/metrics", "URL path for surfacing collected metrics")

	rtorrentAddr     = flag.String("rtorrent.addr", "", "address of rTorrent XML-RPC server")
	rtorrentUsername = flag.String("rtorrent.username", "", "[optional] username used for HTTP Basic authentication with rTorrent XML-RPC server")
	rtorrentPassword = flag.String("rtorrent.password", "", "[optional] password used for HTTP Basic authentication with rTorrent XML-RPC server")
)

func main() {
	flag.Parse()

	if *rtorrentAddr == "" {
		log.Fatal("address of rTorrent XML-RPC server must be specified with '-rtorrent.addr' flag")
	}

	transport := http.DefaultTransport.(*http.Transport).Clone()
	transport.RegisterProtocol("scgi", &rtorrentexporter.SCGITransport{
		DialContext: transport.DialContext,
	})
	var rt http.RoundTripper = transport

	// Optionally enable HTTP Basic authentication
	auth := false
	if u, p := *rtorrentUsername, *rtorrentPassword; u != "" && p != "" {
		rt = &authRoundTripper{
			Username:  u,
			Password:  p,
			Transport: rt,
		}
		auth = true
	}

	c, err := rtorrent.New(*rtorrentAddr, rt)
	if err != nil {
		log.Fatalf("cannot create rTorrent client: %v", err)
	}

	prometheus.MustRegister(rtorrentexporter.New(c))

	http.Handle(*metricsPath, promhttp.Handler())
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, *metricsPath, http.StatusMovedPermanently)
	})

	log.Printf("starting rTorrent exporter on %q for server %q (authentication: %v)",
		*telemetryAddr, *rtorrentAddr, auth)

	if err := http.ListenAndServe(*telemetryAddr, nil); err != nil {
		log.Fatalf("cannot start rTorrent exporter: %s", err)
	}
}

var _ http.RoundTripper = &authRoundTripper{}

// An authRoundTripper is a http.RoundTripper which adds HTTP Basic authentication
// to each HTTP request.
type authRoundTripper struct {
	Username  string
	Password  string
	Transport http.RoundTripper
}

func (rt *authRoundTripper) RoundTrip(r *http.Request) (*http.Response, error) {
	r.SetBasicAuth(rt.Username, rt.Password)
	return rt.Transport.RoundTrip(r)
}
