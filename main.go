package main

import (
	"context"
	"net/http"
	"os"
	"runtime"
	"strings"
	"time"

	arg "github.com/alexflint/go-arg"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/joho/godotenv"
	"github.com/metalmatze/githubql_exporter/collector"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/shurcooL/githubql"
	"golang.org/x/oauth2"
)

var (
	// Version of githubql_exporter.
	Version string
	// Revision or Commit this binary was built from.
	Revision string
	// BuildDate this binary was built.
	BuildDate string
	// GoVersion running this binary.
	GoVersion = runtime.Version()
	// StartTime has the time this was started.
	StartTime = time.Now()
)

// Config gets its content from env and passes it on to different packages
type Config struct {
	Debug       bool   `arg:"env:DEBUG"`
	GitHubToken string `arg:"env:GITHUB_TOKEN"`
	Orgs        string `arg:"env:ORGS"`
	WebAddr     string `arg:"env:WEB_ADDR"`
	WebPath     string `arg:"env:WEB_PATH"`
}

// Token returns a token or an error.
func (c Config) Token() oauth2.TokenSource {
	return oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: c.GitHubToken},
	)
}

func main() {
	_ = godotenv.Load()

	c := Config{
		WebPath: "/metrics",
		WebAddr: ":9276",
	}
	arg.MustParse(&c)

	if c.GitHubToken == "" {
		panic("GITHUB_TOKEN is required")
	}

	filterOption := level.AllowInfo()
	if c.Debug {
		filterOption = level.AllowDebug()
	}

	logger := log.NewLogfmtLogger(log.NewSyncWriter(os.Stderr))
	logger = level.NewFilter(logger, filterOption)
	logger = log.With(logger,
		"ts", log.DefaultTimestampUTC,
		"caller", log.DefaultCaller,
	)

	level.Info(logger).Log(
		"msg", "starting githubql_exporter",
		"version", Version,
		"revision", Revision,
		"buildDate", BuildDate,
		"goVersion", GoVersion,
	)

	httpClient := oauth2.NewClient(context.Background(), c.Token())
	client := githubql.NewClient(httpClient)

	organizations := strings.Split(c.Orgs, ",")

	prometheus.MustRegister(collector.NewOrganizationCollector(logger, client, organizations))

	http.Handle(c.WebPath, promhttp.Handler())
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`<html>
			<head><title>GitHubQL Exporter</title></head>
			<body>
			<h1>GitHubQL Exporter</h1>
			<p><a href="` + c.WebPath + `">Metrics</a></p>
			</body>
			</html>`))
	})

	level.Info(logger).Log("msg", "listening", "addr", c.WebAddr)
	if err := http.ListenAndServe(c.WebAddr, nil); err != nil {
		level.Error(logger).Log("msg", "http listen and serve error", "err", err)
		os.Exit(1)
	}
}
