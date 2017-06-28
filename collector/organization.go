package collector

import (
	"context"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/shurcooL/githubql"
)

const namespace = "github"

// OrganizationCollector collects metrics about the account.
type OrganizationCollector struct {
	logger log.Logger
	client *githubql.Client
	orgs   string

	created   *prometheus.Desc
	diskUsage *prometheus.Desc
	forks     *prometheus.Desc
	issues    *prometheus.Desc
	//pullRequests *prometheus.Desc
	//pushed       *prometheus.Desc
	//stargazers   *prometheus.Desc
	//watchers     *prometheus.Desc
}

type (
	organizationQuery struct {
		Organization struct {
			Login        githubql.String
			Repositories struct {
				Nodes []organizationRepositoryNode
			} `graphql:"repositories(first: 100)"`
		} `graphql:"organization(login: $organization)"`
	}
	organizationRepositoryNode struct {
		Name      githubql.String
		DiskUsage githubql.Int
		CreatedAt githubql.DateTime
		PushedAt  githubql.DateTime

		Stargazers struct {
			TotalCount githubql.Int
		}
		Watchers struct {
			TotalCount githubql.Int
		}
		Forks struct {
			TotalCount githubql.Int
		}
		IssuesOpen struct {
			TotalCount githubql.Int
		} `graphql:"issuesOpen: issues(states: OPEN)"`
		IssuesClosed struct {
			TotalCount githubql.Int
		} `graphql:"issuesClosed: issues(states: CLOSED)"`
		PullRequestsOpen struct {
			TotalCount githubql.Int
		} `graphql:"PullRequestsOpen: pullRequests(states: OPEN)"`
		PullRequestsClosed struct {
			TotalCount githubql.Int
		} `graphql:"PullRequestsClosed: pullRequests(states: CLOSED)"`
		PullRequestsMerged struct {
			TotalCount githubql.Int
		} `graphql:"PullRequestsMerged: pullRequests(states: MERGED)"`
	}
)

// NewOrganizationCollector returns a new OrganizationCollector.
func NewOrganizationCollector(logger log.Logger, client *githubql.Client, orgs string) *OrganizationCollector {
	return &OrganizationCollector{
		logger: logger,
		client: client,
		orgs:   orgs,

		created: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "repo", "created"),
			"Unix timestamp of when the repo was created",
			[]string{"owner", "name"}, nil,
		),
		diskUsage: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "repo", "disk_usage_bytes"),
			"Bytes of the repository used on disk",
			[]string{"owner", "name"}, nil,
		),
		forks: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "repo", "forks"),
			"Number of forks of that repo",
			[]string{"owner", "name"}, nil,
		),
		issues: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "repo", "issues"),
			"Number of issues with a state of open or closed",
			[]string{"owner", "name", "state"}, nil,
		),
	}
}

// Describe sends the super-set of all possible descriptors of metrics
// collected by this Collector.
func (c *OrganizationCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.created
	ch <- c.diskUsage
	ch <- c.forks
	ch <- c.issues
	//ch <- c.pullRequests
	//ch <- c.pushed
	//ch <- c.stargazers
	//ch <- c.watchers
}

// Collect is called by the Prometheus registry when collecting metrics.
func (c *OrganizationCollector) Collect(ch chan<- prometheus.Metric) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	variables := map[string]interface{}{
		"organization": githubql.String(c.orgs),
	}

	var query organizationQuery
	if err := c.client.Query(ctx, &query, variables); err != nil {
		level.Warn(c.logger).Log("msg", "failed to execute organization query successfully", "err", err)
		return
	}

	for _, repo := range query.Organization.Repositories.Nodes {
		ch <- prometheus.MustNewConstMetric(
			c.created,
			prometheus.GaugeValue,
			float64(repo.CreatedAt.Unix()),
			string(query.Organization.Login), string(repo.Name),
		)
		ch <- prometheus.MustNewConstMetric(
			c.diskUsage,
			prometheus.GaugeValue,
			float64(repo.DiskUsage),
			string(query.Organization.Login), string(repo.Name),
		)
		ch <- prometheus.MustNewConstMetric(
			c.forks,
			prometheus.GaugeValue,
			float64(repo.Forks.TotalCount),
			string(query.Organization.Login), string(repo.Name),
		)
		ch <- prometheus.MustNewConstMetric(
			c.issues,
			prometheus.GaugeValue,
			float64(repo.IssuesOpen.TotalCount),
			string(query.Organization.Login), string(repo.Name), "open",
		)
		ch <- prometheus.MustNewConstMetric(
			c.issues,
			prometheus.GaugeValue,
			float64(repo.IssuesClosed.TotalCount),
			string(query.Organization.Login), string(repo.Name), "closed",
		)
	}
}
