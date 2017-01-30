// Command rtorrent_exporter provides a Prometheus exporter for rTorrent.
package main

import (
	"log"
	"net/http"

	"github.com/mdlayher/rtorrent"
	"github.com/mdlayher/rtorrent_exporter"
	"github.com/namsral/flag"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	telemetryAddr = flag.String("telemetry_addr", ":9135", "host:port for rTorrent exporter")
	metricsPath   = flag.String("telemetry_path", "/metrics", "URL path for surfacing collected metrics")

	rtorrentAddr     = flag.String("rtorrent_addr", "", "address of rTorrent XML-RPC server")
	rtorrentUsername = flag.String("rtorrent_username", "", "[optional] username used for HTTP Basic authentication with rTorrent XML-RPC server")
	rtorrentPassword = flag.String("rtorrent_password", "", "[optional] password used for HTTP Basic authentication with rTorrent XML-RPC server")
)

func main() {
	flag.Parse()

	if *rtorrentAddr == "" {
		log.Fatal("address of rTorrent XML-RPC server must be specified with '-rtorrent_addr' flag")
	}

	// Optionally enable HTTP Basic authentication
	var rt http.RoundTripper
	if u, p := *rtorrentUsername, *rtorrentPassword; u != "" && p != "" {
		rt = &authRoundTripper{
			Username: u,
			Password: p,
		}
	}

	c, err := rtorrent.New(*rtorrentAddr, rt)
	if err != nil {
		log.Fatalf("cannot create rTorrent client: %v", err)
	}

	prometheus.MustRegister(rtorrentexporter.New(c))

	http.Handle(*metricsPath, prometheus.Handler())
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, *metricsPath, http.StatusMovedPermanently)
	})

	log.Printf("starting rTorrent exporter on %q for server %q (authentication: %v)",
		*telemetryAddr, *rtorrentAddr, rt != nil)

	if err := http.ListenAndServe(*telemetryAddr, nil); err != nil {
		log.Fatalf("cannot start rTorrent exporter: %s", err)
	}
}

var _ http.RoundTripper = &authRoundTripper{}

// An authRoundTripper is a http.RoundTripper which adds HTTP Basic authentication
// to each HTTP request.
type authRoundTripper struct {
	Username string
	Password string
}

func (rt *authRoundTripper) RoundTrip(r *http.Request) (*http.Response, error) {
	r.SetBasicAuth(rt.Username, rt.Password)
	return http.DefaultTransport.RoundTrip(r)
}
