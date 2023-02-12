package triggers

import (
	"context"
	"testing"
	"time"

	"github.com/kubeshop/testkube/pkg/repository/result"
	"github.com/kubeshop/testkube/pkg/repository/testresult"

	"github.com/golang/mock/gomock"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes/fake"

	executorv1 "github.com/kubeshop/testkube-operator/apis/executor/v1"
	testsv3 "github.com/kubeshop/testkube-operator/apis/tests/v3"

	testtriggersv1 "github.com/kubeshop/testkube-operator/apis/testtriggers/v1"
	executorsclientv1 "github.com/kubeshop/testkube-operator/client/executors/v1"
	testsclientv3 "github.com/kubeshop/testkube-operator/client/tests/v3"
	testsourcesv1 "github.com/kubeshop/testkube-operator/client/testsources/v1"
	testsuitesv2 "github.com/kubeshop/testkube-operator/client/testsuites/v2"
	faketestkube "github.com/kubeshop/testkube-operator/pkg/clientset/versioned/fake"
	"github.com/kubeshop/testkube/internal/app/api/metrics"
	"github.com/kubeshop/testkube/pkg/api/v1/testkube"
	"github.com/kubeshop/testkube/pkg/config"
	"github.com/kubeshop/testkube/pkg/event"
	"github.com/kubeshop/testkube/pkg/event/bus"
	"github.com/kubeshop/testkube/pkg/executor/client"
	"github.com/kubeshop/testkube/pkg/log"
	"github.com/kubeshop/testkube/pkg/scheduler"
	"github.com/kubeshop/testkube/pkg/secret"

	v1 "github.com/kubeshop/testkube-operator/apis/testtriggers/v1"
	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestService_Run(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()
	testMetrics := metrics.NewMetrics()

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockResultRepository := result.NewMockRepository(mockCtrl)
	mockTestResultRepository := testresult.NewMockRepository(mockCtrl)

	mockExecutorsClient := executorsclientv1.NewMockInterface(mockCtrl)
	mockTestsClient := testsclientv3.NewMockInterface(mockCtrl)
	mockTestSuitesClient := testsuitesv2.NewMockInterface(mockCtrl)
	mockTestSourcesClient := testsourcesv1.NewMockInterface(mockCtrl)
	mockSecretClient := secret.NewMockInterface(mockCtrl)
	configMap := config.NewMockRepository(mockCtrl)

	mockExecutor := client.NewMockExecutor(mockCtrl)

	mockEventEmitter := event.NewEmitter(bus.NewEventBusMock())

	mockTest := testsv3.Test{
		ObjectMeta: metav1.ObjectMeta{Namespace: "testkube", Name: "some-test"},
		Spec: testsv3.TestSpec{
			Type_: "cypress",
			ExecutionRequest: &testsv3.ExecutionRequest{
				Name:   "some-custom-execution",
				Number: 1,
				Image:  "test-image",
			},
		},
	}
	mockTestsClient.EXPECT().Get("some-test").Return(&mockTest, nil).AnyTimes()
	var mockNextExecutionNumber int32 = 1
	mockResultRepository.EXPECT().GetNextExecutionNumber(gomock.Any(), "some-test").Return(mockNextExecutionNumber, nil)
	mockExecutionResult := testkube.ExecutionResult{Status: testkube.ExecutionStatusRunning}
	mockExecution := testkube.Execution{Name: "test-execution-1"}
	mockExecution.ExecutionResult = &mockExecutionResult
	mockResultRepository.EXPECT().GetByNameAndTest(gomock.Any(), "some-custom-execution-1", "some-test").Return(mockExecution, nil)
	mockSecretUUID := "b524c2f6-6bcf-4178-87c1-1aa2b2abb5dc"
	mockTestsClient.EXPECT().GetCurrentSecretUUID("some-test").Return(mockSecretUUID, nil)
	mockExecutorTypes := "cypress"
	mockExecutorV1 := executorv1.Executor{
		TypeMeta:   metav1.TypeMeta{},
		ObjectMeta: metav1.ObjectMeta{Namespace: "testkube", Name: "cypress"},
		Spec: executorv1.ExecutorSpec{
			Types:            []string{mockExecutorTypes},
			ExecutorType:     "job",
			URI:              "",
			Image:            "cypress",
			Args:             nil,
			Command:          []string{"run"},
			ImagePullSecrets: nil,
			Features:         nil,
			ContentTypes:     nil,
			JobTemplate:      "",
			Meta:             nil,
		},
	}
	mockExecutorsClient.EXPECT().GetByType(mockExecutorTypes).Return(&mockExecutorV1, nil).AnyTimes()
	mockResultRepository.EXPECT().Insert(gomock.Any(), gomock.Any()).Return(nil)
	mockResultRepository.EXPECT().StartExecution(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
	mockExecutor.EXPECT().Execute(gomock.Any(), gomock.Any(), gomock.Any()).Return(&mockExecutionResult, nil)
	mockResultRepository.EXPECT().UpdateResult(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)

	mockTestExecution := testkube.Execution{
		Id:              "test-suite-execution-1",
		ExecutionResult: &testkube.ExecutionResult{Status: testkube.ExecutionStatusPassed},
	}
	mockResultRepository.EXPECT().Get(gomock.Any(), gomock.Any()).Return(mockTestExecution, nil)

	testLogger := log.DefaultLogger

	sched := scheduler.NewScheduler(
		testMetrics,
		mockExecutor,
		mockExecutor,
		mockResultRepository,
		mockTestResultRepository,
		mockExecutorsClient,
		mockTestsClient,
		mockTestSuitesClient,
		mockTestSourcesClient,
		mockSecretClient,
		mockEventEmitter,
		testLogger,
		configMap,
	)

	mockLeaseBackend := NewMockLeaseBackend(mockCtrl)
	testClusterID := "testkube-api"
	testIdentifier := "test-host-1"
	mockLeaseBackend.EXPECT().TryAcquire(gomock.Any(), testIdentifier, testClusterID).Return(true, nil).AnyTimes()

	fakeTestkubeClientset := faketestkube.NewSimpleClientset()
	fakeClientset := fake.NewSimpleClientset()
	s := NewService(
		sched,
		fakeClientset,
		fakeTestkubeClientset,
		mockTestSuitesClient,
		mockTestsClient,
		mockResultRepository,
		mockTestResultRepository,
		mockLeaseBackend,
		testLogger,
		configMap,
		WithClusterID(testClusterID),
		WithIdentifier(testIdentifier),
		WithScraperInterval(50*time.Millisecond),
		WithLeaseCheckerInterval(50*time.Millisecond),
	)

	s.Run(ctx)

	time.Sleep(100 * time.Millisecond)

	testNamespace := "testkube"
	testTrigger := testtriggersv1.TestTrigger{
		ObjectMeta: metav1.ObjectMeta{Namespace: testNamespace, Name: "test-trigger-1"},
		Spec: testtriggersv1.TestTriggerSpec{
			Resource:         "pod",
			ResourceSelector: testtriggersv1.TestTriggerSelector{Name: "test-pod"},
			Event:            "created",
			Action:           "run",
			Execution:        "test",
			TestSelector:     testtriggersv1.TestTriggerSelector{Name: "some-test"},
		},
	}
	createdTestTrigger, err := fakeTestkubeClientset.TestsV1().TestTriggers(testNamespace).Create(ctx, &testTrigger, metav1.CreateOptions{})
	assert.NotNil(t, createdTestTrigger)
	assert.NoError(t, err)

	time.Sleep(10 * time.Millisecond)

	assert.Len(t, s.triggerStatus, 1)
	key := newStatusKey(testNamespace, "test-trigger-1")
	assert.Contains(t, s.triggerStatus, key)

	testPod := corev1.Pod{ObjectMeta: metav1.ObjectMeta{Namespace: testNamespace, Name: "test-pod", CreationTimestamp: metav1.Now()}}
	_, err = fakeClientset.CoreV1().Pods(testNamespace).Create(ctx, &testPod, metav1.CreateOptions{})
	assert.NoError(t, err)

	<-ctx.Done()
}

func TestService_addTrigger(t *testing.T) {
	t.Parallel()

	s := Service{triggerStatus: make(map[statusKey]*triggerStatus)}

	testTrigger := v1.TestTrigger{
		ObjectMeta: metav1.ObjectMeta{Name: "test-trigger-1", Namespace: "testkube"},
	}
	s.addTrigger(&testTrigger)

	assert.Len(t, s.triggerStatus, 1)
	key := newStatusKey("testkube", "test-trigger-1")
	assert.NotNil(t, s.triggerStatus[key])
}

func TestService_removeTrigger(t *testing.T) {
	t.Parallel()

	s := Service{triggerStatus: make(map[statusKey]*triggerStatus)}

	testTrigger1 := v1.TestTrigger{
		ObjectMeta: metav1.ObjectMeta{Name: "test-trigger-1", Namespace: "testkube"},
	}
	testTrigger2 := v1.TestTrigger{
		ObjectMeta: metav1.ObjectMeta{Name: "test-trigger-2", Namespace: "testkube"},
	}
	s.addTrigger(&testTrigger1)
	s.addTrigger(&testTrigger2)

	assert.Len(t, s.triggerStatus, 2)

	s.removeTrigger(&testTrigger1)

	assert.Len(t, s.triggerStatus, 1)
	key := newStatusKey("testkube", "test-trigger-2")
	assert.NotNil(t, s.triggerStatus[key])
	deletedKey := newStatusKey("testkube", "test-trigger-1")
	assert.Nil(t, s.triggerStatus[deletedKey])
}

func TestService_updateTrigger(t *testing.T) {
	t.Parallel()

	s := Service{triggerStatus: make(map[statusKey]*triggerStatus)}

	oldTestTrigger := v1.TestTrigger{
		ObjectMeta: metav1.ObjectMeta{Namespace: "testkube", Name: "test-trigger-1"},
		Spec:       v1.TestTriggerSpec{Event: "created"},
	}
	s.addTrigger(&oldTestTrigger)

	newTestTrigger := v1.TestTrigger{
		ObjectMeta: metav1.ObjectMeta{Namespace: "testkube", Name: "test-trigger-1"},
		Spec:       v1.TestTriggerSpec{Event: "modified"},
	}

	s.updateTrigger(&newTestTrigger)

	assert.Len(t, s.triggerStatus, 1)
	key := newStatusKey("testkube", "test-trigger-1")
	assert.NotNil(t, s.triggerStatus[key])
}
