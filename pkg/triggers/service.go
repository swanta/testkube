package triggers

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/kubeshop/testkube/pkg/version"

	"github.com/kubeshop/testkube/pkg/repository/result"
	"github.com/kubeshop/testkube/pkg/repository/testresult"

	"github.com/kubeshop/testkube/pkg/config"
	"github.com/kubeshop/testkube/pkg/scheduler"
	"github.com/kubeshop/testkube/pkg/telemetry"
	"github.com/kubeshop/testkube/pkg/utils"

	testsv3 "github.com/kubeshop/testkube-operator/apis/tests/v3"
	testsuitev2 "github.com/kubeshop/testkube-operator/apis/testsuite/v2"
	testtriggersv1 "github.com/kubeshop/testkube-operator/apis/testtriggers/v1"
	testsclientv3 "github.com/kubeshop/testkube-operator/client/tests/v3"
	testsuitesclientv2 "github.com/kubeshop/testkube-operator/client/testsuites/v2"
	testkubeclientsetv1 "github.com/kubeshop/testkube-operator/pkg/clientset/versioned"

	"go.uber.org/zap"
	"k8s.io/client-go/kubernetes"
)

const (
	defaultScraperInterval        = 5 * time.Second
	defaultLeaseCheckInterval     = 5 * time.Second
	defaultMaxLeaseDuration       = 1 * time.Minute
	defaultConditionsCheckBackoff = 1 * time.Second
	defaultConditionsCheckTimeout = 60 * time.Second
	defaultClusterID              = "testkube-api"
	defaultIdentifierFormat       = "testkube-api-%s"
)

type Service struct {
	informers                     *k8sInformers
	leaseBackend                  LeaseBackend
	identifier                    string
	clusterID                     string
	executor                      ExecutorF
	scraperInterval               time.Duration
	leaseCheckInterval            time.Duration
	maxLeaseDuration              time.Duration
	defaultConditionsCheckTimeout time.Duration
	defaultConditionsCheckBackoff time.Duration
	watchFromDate                 time.Time
	triggerStatus                 map[statusKey]*triggerStatus
	scheduler                     *scheduler.Scheduler
	clientset                     kubernetes.Interface
	testKubeClientset             testkubeclientsetv1.Interface
	testSuitesClient              testsuitesclientv2.Interface
	testsClient                   testsclientv3.Interface
	resultRepository              result.Repository
	testResultRepository          testresult.Repository
	logger                        *zap.SugaredLogger
	configMap                     config.Repository
}

type Option func(*Service)

func NewService(
	scheduler *scheduler.Scheduler,
	clientset kubernetes.Interface,
	testKubeClientset testkubeclientsetv1.Interface,
	testSuitesClient testsuitesclientv2.Interface,
	testsClient testsclientv3.Interface,
	resultRepository result.Repository,
	testResultRepository testresult.Repository,
	leaseBackend LeaseBackend,
	logger *zap.SugaredLogger,
	configMap config.Repository,
	opts ...Option,
) *Service {
	identifier := fmt.Sprintf(defaultIdentifierFormat, utils.RandAlphanum(10))
	s := &Service{
		informers:                     newK8sInformers(clientset, testKubeClientset),
		identifier:                    identifier,
		clusterID:                     defaultClusterID,
		scraperInterval:               defaultScraperInterval,
		leaseCheckInterval:            defaultLeaseCheckInterval,
		maxLeaseDuration:              defaultMaxLeaseDuration,
		defaultConditionsCheckTimeout: defaultConditionsCheckTimeout,
		defaultConditionsCheckBackoff: defaultConditionsCheckBackoff,
		scheduler:                     scheduler,
		clientset:                     clientset,
		testKubeClientset:             testKubeClientset,
		testSuitesClient:              testSuitesClient,
		testsClient:                   testsClient,
		resultRepository:              resultRepository,
		testResultRepository:          testResultRepository,
		leaseBackend:                  leaseBackend,
		logger:                        logger,
		configMap:                     configMap,
		watchFromDate:                 time.Now(),
		triggerStatus:                 make(map[statusKey]*triggerStatus),
	}
	if s.executor == nil {
		s.executor = s.execute
	}

	for _, opt := range opts {
		opt(s)
	}

	return s
}

func WithIdentifier(id string) Option {
	return func(s *Service) {
		s.identifier = id
	}
}

func WithHostnameIdentifier() Option {
	return func(s *Service) {
		identifier, err := os.Hostname()
		if err == nil {
			s.identifier = identifier
		}
	}
}

func WithClusterID(id string) Option {
	return func(s *Service) {
		s.clusterID = id
	}
}

func WithWatchFromDate(from time.Time) Option {
	return func(s *Service) {
		s.watchFromDate = from
	}
}

func WithLeaseCheckerInterval(interval time.Duration) Option {
	return func(s *Service) {
		s.leaseCheckInterval = interval
	}
}

func WithScraperInterval(interval time.Duration) Option {
	return func(s *Service) {
		s.scraperInterval = interval
	}
}

func WithExecutor(executor ExecutorF) Option {
	return func(s *Service) {
		s.executor = executor
	}
}

func (s *Service) Run(ctx context.Context) {
	leaseChan := make(chan bool)

	go s.runLeaseChecker(ctx, leaseChan)

	go s.runWatcher(ctx, leaseChan)

	go s.runExecutionScraper(ctx)
}

func (s *Service) addTrigger(t *testtriggersv1.TestTrigger) {
	key := newStatusKey(t.Namespace, t.Name)
	s.triggerStatus[key] = newTriggerStatus(t)
}

func (s *Service) updateTrigger(target *testtriggersv1.TestTrigger) {
	key := newStatusKey(target.Namespace, target.Name)
	if s.triggerStatus[key] != nil {
		s.triggerStatus[key].testTrigger = target
	} else {
		s.triggerStatus[key] = newTriggerStatus(target)
	}
}

func (s *Service) removeTrigger(target *testtriggersv1.TestTrigger) {
	key := newStatusKey(target.Namespace, target.Name)
	delete(s.triggerStatus, key)
}

func (s *Service) addTest(test *testsv3.Test) {
	ctx := context.Background()
	telemetryEnabled, err := s.configMap.GetTelemetryEnabled(ctx)
	if err != nil {
		s.logger.Debugw("getting telemetry enabled error", "error", err)
	}

	if !telemetryEnabled {
		return
	}

	clusterID, err := s.configMap.GetUniqueClusterId(ctx)
	if err != nil {
		s.logger.Debugw("getting cluster id error", "error", err)
	}

	host, err := os.Hostname()
	if err != nil {
		s.logger.Debugw("getting hostname error", "hostname", host, "error", err)
	}

	var dataSource string
	if test.Spec.Content != nil {
		dataSource = test.Spec.Content.Type_
	}

	out, err := telemetry.SendCreateEvent("testkube_api_create_test", telemetry.CreateParams{
		AppVersion: version.Version,
		DataSource: dataSource,
		Host:       host,
		ClusterID:  clusterID,
		TestType:   test.Spec.Type_,
		TestSource: test.Spec.Source,
	})
	if err != nil {
		s.logger.Debugw("sending create test telemetry event error", "error", err)
	} else {
		s.logger.Debugw("sending create test telemetry event", "output", out)
	}
}

func (s *Service) addTestSuite(testSuite *testsuitev2.TestSuite) {
	ctx := context.Background()
	telemetryEnabled, err := s.configMap.GetTelemetryEnabled(ctx)
	if err != nil {
		s.logger.Debugw("getting telemetry enabled error", "error", err)
	}

	if !telemetryEnabled {
		return
	}

	clusterID, err := s.configMap.GetUniqueClusterId(ctx)
	if err != nil {
		s.logger.Debugw("getting cluster id error", "error", err)
	}

	host, err := os.Hostname()
	if err != nil {
		s.logger.Debugw("getting hostname error", "hostname", host, "error", err)
	}

	out, err := telemetry.SendCreateEvent("testkube_api_create_test_suite", telemetry.CreateParams{
		AppVersion:     version.Version,
		Host:           host,
		ClusterID:      clusterID,
		TestSuiteSteps: int32(len(testSuite.Spec.Before) + len(testSuite.Spec.Steps) + len(testSuite.Spec.After)),
	})
	if err != nil {
		s.logger.Debugw("sending create test suite telemetry event error", "error", err)
	} else {
		s.logger.Debugw("sending create test suite telemetry event", "output", out)
	}
}
