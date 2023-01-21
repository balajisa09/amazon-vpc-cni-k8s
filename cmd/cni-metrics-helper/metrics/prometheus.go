package metrics

import (
	"net/http"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/log"
	"github.com/prometheus/client_golang/prometheus"
	"strconv"
	"time"
)

var (
	awsAPILatency = prometheus.NewSummaryVec(
		prometheus.SummaryOpts{
			Name: "awscni_aws_api_latency_ms",
			Help: "AWS API call latency in ms",
		},
		[]string{"api", "error", "status"},
	)
	awsAPIErr = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "awscni_aws_api_error_count",
			Help: "The number of times AWS API returns an error",
		},
		[]string{"api", "error"},
	)
	awsUtilsErr = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "awscni_aws_utils_error_count",
			Help: "The number of errors not handled in awsutils library",
		},
		[]string{"fn", "error"},
	)
	ec2ApiReq = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "awscni_ec2api_req_count",
			Help: "The number of requests made to EC2 APIs by CNI",
		},
		[]string{"fn"},
	)
	ec2ApiErr = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "awscni_ec2api_error_count",
			Help: "The number of failed EC2 APIs requests",
		},
		[]string{"fn"},
	)
	ipamdErr = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "awscni_ipamd_error_count",
			Help: "The number of errors encountered in ipamd",
		},
		[]string{"fn"},
	)
	ipamdActionsInprogress = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "awscni_ipamd_action_inprogress",
			Help: "The number of ipamd actions in progress",
		},
		[]string{"fn"},
	)
	enisMax = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "awscni_eni_max",
			Help: "The maximum number of ENIs that can be attached to the instance, accounting for unmanaged ENIs",
		},
	)
	ipMax = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "awscni_ip_max",
			Help: "The maximum number of IP addresses that can be allocated to the instance",
		},
	)
	reconcileCnt = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "awscni_reconcile_count",
			Help: "The number of times ipamd reconciles on ENIs and IP/Prefix addresses",
		},
		[]string{"fn"},
	)
	addIPCnt = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "awscni_add_ip_req_count",
			Help: "The number of add IP address requests",
		},
		[]string{"fn"},
	)
	delIPCnt = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "awscni_del_ip_req_count",
			Help: "The number of delete IP address requests",
		},
		[]string{"reason"},
	)
	podENIErr = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "awscni_pod_eni_error_count",
			Help: "The number of errors encountered for pod ENIs",
		},
		[]string{"fn"},
	)

	enis = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "awscni_eni_allocated",
			Help: "The number of ENIs allocated",
		},
	)
	totalIPs = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "awscni_total_ip_addresses",
			Help: "The total number of IP addresses",
		},
	)
	assignedIPs = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "awscni_assigned_ip_addresses",
			Help: "The number of IP addresses assigned to pods",
		},
	)
	forceRemovedENIs = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "awscni_force_removed_enis",
			Help: "The number of ENIs force removed while they had assigned pods",
		},
		[]string{"fn"},
	)
	forceRemovedIPs = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "awscni_force_removed_ips",
			Help: "The number of IPs force removed while they had assigned pods",
		},
		[]string{"fn"},
	)
	totalPrefixes = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "awscni_total_ipv4_prefixes",
			Help: "The total number of IPv4 prefixes",
		},
	)
	ipsPerCidr = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "awscni_assigned_ip_per_cidr",
			Help: "The total number of IP addresses assigned per cidr",
		},
		[]string{"cidr"},
	)
)

const (
	// metricsPort is the port for prometheus metrics for cni-metrics-helper
	prometheusMetricsPort = 2112
)

func setupPrometheusMetricsServer() *http.Server {
	serveMux := http.NewServeMux()
	serveMux.Handle("/metrics", promhttp.Handler())
	server := &http.Server{
		Addr:         ":" + strconv.Itoa(prometheusMetricsPort),
		Handler:      serveMux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}
	return server
}

func StartPrometheusMetricsServer(){
	log.Info("Starting prometehus metrics server for cni-metrics-helper")
	prometheusServer := setupPrometheusMetricsServer()
	prometheusServer.ListenAndServe()
}

func init(){
	prometheus.MustRegister(awsAPILatency)
	prometheus.MustRegister(awsAPIErr)
	prometheus.MustRegister(awsUtilsErr)
	prometheus.MustRegister(ec2ApiReq)
	prometheus.MustRegister(ec2ApiErr)
	prometheus.MustRegister(ipamdErr)
	prometheus.MustRegister(ipamdActionsInprogress)
	prometheus.MustRegister(enisMax)
	prometheus.MustRegister(ipMax)
	prometheus.MustRegister(reconcileCnt)
	prometheus.MustRegister(addIPCnt)
	prometheus.MustRegister(delIPCnt)
	prometheus.MustRegister(podENIErr)
	prometheus.MustRegister(enis)
	prometheus.MustRegister(totalIPs)
	prometheus.MustRegister(assignedIPs)
	prometheus.MustRegister(forceRemovedENIs)
	prometheus.MustRegister(forceRemovedIPs)
	prometheus.MustRegister(totalPrefixes)
	prometheus.MustRegister(ipsPerCidr)
}
