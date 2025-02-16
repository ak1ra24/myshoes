package metric

import (
	"context"
	"fmt"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/whywaita/myshoes/internal/config"
	"github.com/whywaita/myshoes/pkg/datastore"
	"github.com/whywaita/myshoes/pkg/gh"
	"github.com/whywaita/myshoes/pkg/runner"
	"github.com/whywaita/myshoes/pkg/starter"
)

const memoryName = "memory"

var (
	memoryStarterMaxRunning = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, memoryName, "starter_max_running"),
		"The number of max running in starter (Config)",
		[]string{"starter"}, nil,
	)
	memoryStarterQueueRunning = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, memoryName, "starter_queue_running"),
		"running queue in starter",
		[]string{"starter"}, nil,
	)
	memoryStarterQueueWaiting = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, memoryName, "starter_queue_waiting"),
		"waiting queue in starter",
		[]string{"starter"}, nil,
	)
	memoryGitHubRateLimitRemaining = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, memoryName, "github_rate_limit_remaining"),
		"The number of rate limit remaining",
		[]string{"scope"}, nil,
	)
	memoryGitHubRateLimitLimiting = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, memoryName, "github_rate_limit_limiting"),
		"The number of rate limit max",
		[]string{"scope"}, nil,
	)
	memoryRunnerMaxConcurrencyDeleting = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, memoryName, "runner_max_concurrency_deleting"),
		"The number of max concurrency deleting in runner (Config)",
		[]string{"runner"}, nil,
	)
	memoryRunnerQueueConcurrencyDeleting = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, memoryName, "runner_queue_concurrency_deleting"),
		"deleting concurrency in runner",
		[]string{"runner"}, nil,
	)
)

// ScraperMemory is scraper implement for memory
type ScraperMemory struct{}

// Name return name
func (ScraperMemory) Name() string {
	return memoryName
}

// Help return help
func (ScraperMemory) Help() string {
	return "Collect from memory"
}

// Scrape scrape metrics
func (ScraperMemory) Scrape(ctx context.Context, ds datastore.Datastore, ch chan<- prometheus.Metric) error {
	if err := scrapeStarterValues(ch); err != nil {
		return fmt.Errorf("failed to scrape starter values: %w", err)
	}
	if err := scrapeGitHubValues(ch); err != nil {
		return fmt.Errorf("failed to scrape GitHub values: %w", err)
	}

	return nil
}

func scrapeStarterValues(ch chan<- prometheus.Metric) error {
	configMax := config.Config.MaxConnectionsToBackend

	const labelStarter = "starter"

	ch <- prometheus.MustNewConstMetric(
		memoryStarterMaxRunning, prometheus.GaugeValue, float64(configMax), labelStarter)

	countRunning := starter.CountRunning
	countWaiting := starter.CountWaiting

	ch <- prometheus.MustNewConstMetric(
		memoryStarterQueueRunning, prometheus.GaugeValue, float64(countRunning), labelStarter)
	ch <- prometheus.MustNewConstMetric(
		memoryStarterQueueWaiting, prometheus.GaugeValue, float64(countWaiting), labelStarter)

	const labelRunner = "runner"
	configRunnerDeletingMax := config.Config.MaxConcurrencyDeleting
	countRunnerDeletingNow := runner.ConcurrencyDeleting

	ch <- prometheus.MustNewConstMetric(
		memoryRunnerMaxConcurrencyDeleting, prometheus.GaugeValue, float64(configRunnerDeletingMax), labelRunner)
	ch <- prometheus.MustNewConstMetric(
		memoryRunnerQueueConcurrencyDeleting, prometheus.GaugeValue, float64(countRunnerDeletingNow), labelRunner)

	return nil
}

func scrapeGitHubValues(ch chan<- prometheus.Metric) error {
	rateLimitRemain := gh.GetRateLimitRemain()
	for scope, remain := range rateLimitRemain {
		ch <- prometheus.MustNewConstMetric(
			memoryGitHubRateLimitRemaining, prometheus.GaugeValue, float64(remain), scope,
		)
	}

	rateLimitLimit := gh.GetRateLimitLimit()
	for scope, limit := range rateLimitLimit {
		ch <- prometheus.MustNewConstMetric(
			memoryGitHubRateLimitLimiting, prometheus.GaugeValue, float64(limit), scope,
		)
	}

	return nil
}

var _ Scraper = ScraperMemory{}
