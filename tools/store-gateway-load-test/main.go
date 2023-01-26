package main

import (
	"context"
	"flag"
	"fmt"
	"math"
	"math/rand"
	"os"
	"strconv"
	"time"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/grafana/dskit/dns"
	"github.com/grafana/dskit/flagext"
	"github.com/grafana/dskit/kv"
	"github.com/grafana/dskit/kv/codec"
	"github.com/grafana/dskit/kv/memberlist"
	"github.com/grafana/dskit/ring"
	"github.com/grafana/dskit/services"
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/prometheus/model/labels"
	"github.com/prometheus/prometheus/promql/parser"
	"golang.org/x/sync/errgroup"

	"github.com/grafana/mimir/pkg/mimir"
	"github.com/grafana/mimir/pkg/querier"
	"github.com/grafana/mimir/pkg/storage/bucket"
	"github.com/grafana/mimir/pkg/storage/tsdb/bucketindex"
	"github.com/grafana/mimir/pkg/storegateway"
	"github.com/grafana/mimir/pkg/storegateway/storepb"
	util_log "github.com/grafana/mimir/pkg/util/log"
)

type Config struct {
	Mimir mimir.Config `yaml:"inline"`

	// Tested config.
	TesterUserID                      string
	TesterRequestMinTime              flagext.Time
	TesterRequestMaxTime              flagext.Time
	TesterRequestMinRange             time.Duration
	TesterRequestMaxRange             time.Duration
	TesterRequestMetricNameRegex      string
	TesterRequestPromQLMatchers       flagext.StringSliceCSV
	TesterRequestRandomLabelName      string
	TesterRequestSkipChunks           bool
	TesterConcurrency                 int
	TesterComparisonAuthoritativeZone string
}

func (c *Config) RegisterFlags(f *flag.FlagSet, logger log.Logger) {
	c.Mimir.RegisterFlags(f, logger)

	f.StringVar(&c.TesterUserID, "tester.user-id", "anonymous", "The user ID to run queries.")
	f.Var(&c.TesterRequestMinTime, "tester.request-min-time", fmt.Sprintf("The minimum time to query. The supported time format is %q.", time.RFC3339))
	f.Var(&c.TesterRequestMaxTime, "tester.request-max-time", fmt.Sprintf("The maximum time to query. The supported time format is %q.", time.RFC3339))
	f.DurationVar(&c.TesterRequestMinRange, "tester.request-min-range", 24*time.Hour, "The minimum time range to query within the configured min and max time.")
	f.DurationVar(&c.TesterRequestMaxRange, "tester.request-max-range", 7*24*time.Hour, "The maximum time range to query within the configured min and max time.")
	f.StringVar(&c.TesterRequestMetricNameRegex, "tester.request-metric-name-regex", "up", "The metric name regex to use in the request to store-gateways.")
	f.Var(&c.TesterRequestPromQLMatchers, "tester.request-promql-matchers", `Additional PromQL matchers. Format is comma-separated individual matchers: pod=~"a.*",cluster!="123"`)
	f.StringVar(&c.TesterRequestRandomLabelName, "tester.request-label-name-with-random-values", "", "Label name to add to all requests' matchers that will have a random value.")
	f.BoolVar(&c.TesterRequestSkipChunks, "tester.request-skip-chunks", false, "True to request series without chunks.")
	f.IntVar(&c.TesterConcurrency, "tester.concurrency", 1, "The number of concurrent requests to run.")
	f.StringVar(&c.TesterComparisonAuthoritativeZone, "tester.comparison-authoritative-zone", "zone-c", "The name of the zone to compare results against. This should be the zone expected to return the expected results.")
}

func (c *Config) Validate(logger log.Logger) error {
	return c.Mimir.Validate(logger)
}

func main() {
	ctx := context.Background()
	logger := log.NewLogfmtLogger(os.Stdout)
	reg := prometheus.NewRegistry()

	// IMPORTANT: we assume store-gateway shuffle sharding is disabled (or the tenant blocks are sharded across all store-gateways).
	limits := &noBlocksStoreLimits{}

	// Parse the config.
	cfg := &Config{}
	cfg.RegisterFlags(flag.CommandLine, logger)
	if err := flag.CommandLine.Parse(os.Args[1:]); err != nil {
		fmt.Fprintln(flag.CommandLine.Output(), "Run with -help to get a list of available parameters")
		os.Exit(1)
	}

	if err := cfg.Validate(logger); err != nil {
		panic(err)
	}

	// Init memberlist.
	memberlistKV, err := initMemberlistKV(cfg, reg)
	if err != nil {
		panic(err)
	}

	if err := services.StartAndAwaitRunning(ctx, memberlistKV); err != nil {
		panic(err)
	}

	// Init the bucket client.
	bucketClient, err := bucket.NewClient(context.Background(), cfg.Mimir.BlocksStorage.Bucket, "querier", logger, reg)
	if err != nil {
		panic(errors.Wrap(err, "failed to create bucket client"))
	}

	// Init the finder.
	finder := querier.NewBucketIndexBlocksFinder(querier.BucketIndexBlocksFinderConfig{
		IndexLoader: bucketindex.LoaderConfig{
			CheckInterval:         time.Minute,
			UpdateOnStaleInterval: cfg.Mimir.BlocksStorage.BucketStore.SyncInterval,
			UpdateOnErrorInterval: cfg.Mimir.BlocksStorage.BucketStore.BucketIndex.UpdateOnErrorInterval,
			IdleTimeout:           cfg.Mimir.BlocksStorage.BucketStore.BucketIndex.IdleTimeout,
		},
		MaxStalePeriod:           cfg.Mimir.BlocksStorage.BucketStore.BucketIndex.MaxStalePeriod,
		IgnoreDeletionMarksDelay: cfg.Mimir.BlocksStorage.BucketStore.IgnoreDeletionMarksDelay,
	}, bucketClient, limits, logger, reg)

	if err := services.StartAndAwaitRunning(ctx, finder); err != nil {
		panic(err)
	}

	// Init the selector.
	storesRingCfg := cfg.Mimir.StoreGateway.ShardingRing.ToRingConfig()
	storesRingBackend, err := kv.NewClient(
		storesRingCfg.KVStore,
		ring.GetCodec(),
		kv.RegistererWithKVName(prometheus.WrapRegistererWithPrefix("cortex_", reg), "querier-store-gateway"),
		logger,
	)
	if err != nil {
		panic(errors.Wrap(err, "failed to create store-gateway ring backend"))
	}

	storesRing, err := ring.NewWithStoreClientAndStrategy(storesRingCfg, storegateway.RingNameForClient, storegateway.RingKey, storesRingBackend, ring.NewIgnoreUnhealthyInstancesReplicationStrategy(), prometheus.WrapRegistererWithPrefix("cortex_", reg), logger)
	if err != nil {
		panic(errors.Wrap(err, "failed to create store-gateway ring client"))
	}

	if err := services.StartAndAwaitRunning(ctx, storesRing); err != nil {
		panic(err)
	}

	selector := newStoreGatewaySelector(storesRing, cfg.Mimir.Querier.StoreGatewayClient, limits, logger, reg)

	matchers, err := parsePromQLMatchers(cfg.TesterRequestPromQLMatchers)
	if err != nil {
		panic(errors.Wrap(err, "parsing additional matchers matchers"))
	}

	// Request.
	matchers = append(matchers, storepb.LabelMatcher{
		Type:  storepb.LabelMatcher_RE,
		Name:  labels.MetricName,
		Value: cfg.TesterRequestMetricNameRegex,
	})
	if cfg.TesterRequestRandomLabelName != "" {
		matchers = append(matchers, storepb.LabelMatcher{
			Type:  storepb.LabelMatcher_EQ,
			Name:  cfg.TesterRequestRandomLabelName,
			Value: "",
		})
	}

	logger.Log("msg", "config", "matchers", fmt.Sprintf("%v", matchers))

	t := newTester(cfg.TesterUserID, finder, selector, cfg.TesterComparisonAuthoritativeZone, logger)
	g, _ := errgroup.WithContext(ctx)

	for c := 0; c < cfg.TesterConcurrency; c++ {
		// Compare results only on 1 request, to reduce memory utilization.
		compareResults := c == 0
		myMatchers := make([]storepb.LabelMatcher, len(matchers))
		copy(myMatchers, matchers)

		g.Go(func() error {
			for {
				start, end := getRandomRequestTimeRange(cfg)
				if cfg.TesterRequestRandomLabelName != "" {
					myMatchers[len(myMatchers)-1].Value = strconv.FormatInt(rand.Int63n(math.MaxInt64), 10)
				}
				//level.Info(logger).Log("msg", "request", "start", time.UnixMilli(start).String(), "end", time.UnixMilli(end).String(), "metric name regexp", cfg.TesterMetricNameRegex)

				if err := t.sendRequestToAllStoreGatewayZonesAndCompareResults(ctx, start, end, myMatchers, cfg.TesterRequestSkipChunks, compareResults); err != nil {
					level.Error(logger).Log("msg", "failed to run test", "err", err)
				}
			}
		})
	}

	if err := g.Wait(); err != nil {
		level.Error(logger).Log("msg", "test execution failed", "err", err)
	}
}

func getRandomRequestTimeRange(cfg *Config) (start, end int64) {
	// Get a random time range duration, honoring the configured min and max range.
	timeRangeDurationMillis := (cfg.TesterRequestMinRange + time.Duration(rand.Int63n(int64(cfg.TesterRequestMaxRange-cfg.TesterRequestMinRange)))).Milliseconds()

	// Get a random min timestamp, honoring the configured min and max time.
	minTimeMillis := time.Time(cfg.TesterRequestMinTime).UnixMilli()
	maxTimeMillis := time.Time(cfg.TesterRequestMaxTime).UnixMilli()

	if timeRangeDurationMillis < maxTimeMillis-minTimeMillis {
		start = minTimeMillis + rand.Int63n(maxTimeMillis-minTimeMillis-timeRangeDurationMillis)
	} else {
		start = minTimeMillis
	}

	if start+timeRangeDurationMillis <= maxTimeMillis {
		end = start + timeRangeDurationMillis
	} else {
		end = maxTimeMillis
	}

	return
}

func parsePromQLMatchers(ms []string) ([]storepb.LabelMatcher, error) {
	var promMatchers []*labels.Matcher
	for _, m := range ms {
		parsed, err := parser.ParseMetricSelector("{" + m + "}")
		if err != nil {
			return nil, err
		}
		promMatchers = append(promMatchers, parsed...)

	}
	return storepb.PromMatchersToMatchers(promMatchers...)
}

type noBlocksStoreLimits struct{}

func (l *noBlocksStoreLimits) S3SSEType(userID string) string {
	return ""
}

func (l *noBlocksStoreLimits) S3SSEKMSKeyID(userID string) string {
	return ""
}

func (l *noBlocksStoreLimits) S3SSEKMSEncryptionContext(userID string) string {
	return ""
}

func (l *noBlocksStoreLimits) MaxLabelsQueryLength(userID string) time.Duration {
	return 0
}

func (l *noBlocksStoreLimits) MaxChunksPerQuery(userID string) int {
	return 0
}

func (l *noBlocksStoreLimits) StoreGatewayTenantShardSize(userID string) int {
	return 0
}

func initMemberlistKV(cfg *Config, reg prometheus.Registerer) (services.Service, error) {
	cfg.Mimir.MemberlistKV.MetricsRegisterer = reg
	cfg.Mimir.MemberlistKV.Codecs = []codec.Codec{
		ring.GetCodec(),
	}
	dnsProviderReg := prometheus.WrapRegistererWithPrefix(
		"cortex_",
		prometheus.WrapRegistererWith(
			prometheus.Labels{"component": "memberlist"},
			reg,
		),
	)
	dnsProvider := dns.NewProvider(util_log.Logger, dnsProviderReg, dns.GolangResolverType)
	memberlistKV := memberlist.NewKVInitService(&cfg.Mimir.MemberlistKV, util_log.Logger, dnsProvider, reg)

	// Update the config.
	cfg.Mimir.Distributor.DistributorRing.KVStore.MemberlistKV = memberlistKV.GetMemberlistKV
	cfg.Mimir.Ingester.IngesterRing.KVStore.MemberlistKV = memberlistKV.GetMemberlistKV
	cfg.Mimir.StoreGateway.ShardingRing.KVStore.MemberlistKV = memberlistKV.GetMemberlistKV
	cfg.Mimir.Compactor.ShardingRing.KVStore.MemberlistKV = memberlistKV.GetMemberlistKV
	cfg.Mimir.Ruler.Ring.KVStore.MemberlistKV = memberlistKV.GetMemberlistKV
	cfg.Mimir.Alertmanager.ShardingRing.KVStore.MemberlistKV = memberlistKV.GetMemberlistKV
	cfg.Mimir.QueryScheduler.ServiceDiscovery.SchedulerRing.KVStore.MemberlistKV = memberlistKV.GetMemberlistKV

	return memberlistKV, nil
}