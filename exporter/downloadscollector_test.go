package rtorrentexporter

import (
	"regexp"
	"strings"
	"testing"
)

func TestDownloadsCollector(t *testing.T) {
	// Info-hashes are usually 40 bytes, but we will make ours only
	// 4 for easier testing
	hashA := strings.Repeat("A", 4)
	hashB := strings.Repeat("B", 4)
	hashC := strings.Repeat("C", 4)

	var tests = []struct {
		desc    string
		ds      *testDownloadsSource
		matches []*regexp.Regexp
	}{
		{
			desc: "one download",
			ds: &testDownloadsSource{
				downloads: []string{
					hashA,
				},
				files: map[string]string{
					hashA: "foo",
				},
				rate: 1024,
			},
			matches: []*regexp.Regexp{
				regexp.MustCompile(`rtorrent_downloads 1`),
				regexp.MustCompile(`rtorrent_downloads_started 1`),
				regexp.MustCompile(`rtorrent_downloads_stopped 1`),
				regexp.MustCompile(`rtorrent_downloads_complete 1`),
				regexp.MustCompile(`rtorrent_downloads_incomplete 1`),
				regexp.MustCompile(`rtorrent_downloads_hashing 1`),
				regexp.MustCompile(`rtorrent_downloads_seeding 1`),
				regexp.MustCompile(`rtorrent_downloads_leeching 1`),
				regexp.MustCompile(`rtorrent_downloads_active 1`),

				regexp.MustCompile(`rtorrent_downloads_download_rate_bytes{info_hash="AAAA",name="foo"} 1024`),
				regexp.MustCompile(`rtorrent_downloads_upload_rate_bytes{info_hash="AAAA",name="foo"} 1024`),
			},
		},
		{
			desc: "three downloads",
			ds: &testDownloadsSource{
				downloads: []string{
					hashA,
					hashB,
					hashC,
				},
				files: map[string]string{
					hashA: "foo",
					hashB: "bar",
					hashC: "baz",
				},
				rate: 2048,
			},
			matches: []*regexp.Regexp{
				regexp.MustCompile(`rtorrent_downloads 3`),
				regexp.MustCompile(`rtorrent_downloads_started 3`),
				regexp.MustCompile(`rtorrent_downloads_stopped 3`),
				regexp.MustCompile(`rtorrent_downloads_complete 3`),
				regexp.MustCompile(`rtorrent_downloads_incomplete 3`),
				regexp.MustCompile(`rtorrent_downloads_hashing 3`),
				regexp.MustCompile(`rtorrent_downloads_seeding 3`),
				regexp.MustCompile(`rtorrent_downloads_leeching 3`),
				regexp.MustCompile(`rtorrent_downloads_active 3`),

				regexp.MustCompile(`rtorrent_downloads_download_rate_bytes{info_hash="AAAA",name="foo"} 2048`),
				regexp.MustCompile(`rtorrent_downloads_upload_rate_bytes{info_hash="AAAA",name="foo"} 2048`),
				regexp.MustCompile(`rtorrent_downloads_download_rate_bytes{info_hash="BBBB",name="bar"} 2048`),
				regexp.MustCompile(`rtorrent_downloads_upload_rate_bytes{info_hash="BBBB",name="bar"} 2048`),
				regexp.MustCompile(`rtorrent_downloads_download_rate_bytes{info_hash="CCCC",name="baz"} 2048`),
				regexp.MustCompile(`rtorrent_downloads_upload_rate_bytes{info_hash="CCCC",name="baz"} 2048`),
			},
		},
	}

	for i, tt := range tests {
		t.Logf("[%02d] test %q", i, tt.desc)

		out := testCollector(t, NewDownloadsCollector(tt.ds))

		for j, m := range tt.matches {
			t.Logf("\t[%02d:%02d] match: %s", i, j, m.String())

			if !m.Match(out) {
				t.Fatal("\toutput failed to match regex")
			}
		}
	}
}

var _ DownloadsSource = &testDownloadsSource{}

type testDownloadsSource struct {
	downloads []string
	files     map[string]string
	rate      int
}

func (ds *testDownloadsSource) All() ([]string, error)        { return ds.downloads, nil }
func (ds *testDownloadsSource) Started() ([]string, error)    { return ds.downloads, nil }
func (ds *testDownloadsSource) Stopped() ([]string, error)    { return ds.downloads, nil }
func (ds *testDownloadsSource) Complete() ([]string, error)   { return ds.downloads, nil }
func (ds *testDownloadsSource) Incomplete() ([]string, error) { return ds.downloads, nil }
func (ds *testDownloadsSource) Hashing() ([]string, error)    { return ds.downloads, nil }
func (ds *testDownloadsSource) Seeding() ([]string, error)    { return ds.downloads, nil }
func (ds *testDownloadsSource) Leeching() ([]string, error)   { return ds.downloads, nil }
func (ds *testDownloadsSource) Active() ([]string, error)     { return ds.downloads, nil }

func (ds *testDownloadsSource) BaseFilename(infoHash string) (string, error) {
	return ds.files[infoHash], nil
}
func (ds *testDownloadsSource) DownloadRate(_ string) (int, error) { return ds.rate, nil }
func (ds *testDownloadsSource) UploadRate(_ string) (int, error)   { return ds.rate, nil }
