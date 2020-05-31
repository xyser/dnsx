package health

import (
	"github.com/prometheus/client_golang/prometheus"
)

// ClusterManager link: https://blog.csdn.net/u014029783/article/details/80001251
type ClusterManager struct {
	Zone         string
	RecordsCount *prometheus.Desc
	DomainsCount *prometheus.Desc
}

// Simulate prepare the data
func (c *ClusterManager) ReallyExpensiveAssessmentOfTheSystemState() (
	recordsCount int64, domainsCount int64) {

	return 0, 0
}

// Describe simply sends the two Descs in the struct to the channel.
func (c *ClusterManager) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.RecordsCount
	ch <- c.DomainsCount
}

func (c *ClusterManager) Collect(ch chan<- prometheus.Metric) {
	domainsCount, recordsCount := c.ReallyExpensiveAssessmentOfTheSystemState()
	ch <- prometheus.MustNewConstMetric(
		c.DomainsCount,
		prometheus.CounterValue,
		float64(domainsCount),
	)
	ch <- prometheus.MustNewConstMetric(
		c.RecordsCount,
		prometheus.CounterValue,
		float64(recordsCount),
	)
}

// NewClusterManager creates the two Descs OOMCountDesc and RAMUsageDesc. Note
// that the zone is set as a ConstLabel. (It's different in each instance of the
// ClusterManager, but constant over the lifetime of an instance.) Then there is
// a variable label "host", since we want to partition the collected metrics by
// host. Since all Descs created in this way are consistent across instances,
// with a guaranteed distinction by the "zone" label, we can register different
// ClusterManager instances with the same registry.
func NewClusterManager(zone string) *ClusterManager {
	return &ClusterManager{
		Zone: zone,
		DomainsCount: prometheus.NewDesc(
			"message_queue_count",
			"Message Center Timer Count.",
			[]string{},
			prometheus.Labels{"zone": zone},
		),
		RecordsCount: prometheus.NewDesc(
			"timer_count",
			"Message Center Timer Count.",
			[]string{},
			prometheus.Labels{"zone": zone},
		),
	}
}
