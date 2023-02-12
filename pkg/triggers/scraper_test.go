package triggers

import (
	"context"
	"testing"
	"time"

	"github.com/kubeshop/testkube/pkg/repository/result"
	"github.com/kubeshop/testkube/pkg/repository/testresult"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/kubeshop/testkube/pkg/api/v1/testkube"
	"github.com/kubeshop/testkube/pkg/log"
)

func TestService_runExecutionScraper(t *testing.T) {
	t.Parallel()

	t.Run("completed jobs", func(t *testing.T) {
		t.Parallel()

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		ctx, cancel := context.WithTimeout(context.Background(), 310*time.Millisecond)
		defer cancel()

		mockResultRepository := result.NewMockRepository(mockCtrl)
		mockTestResultRepository := testresult.NewMockRepository(mockCtrl)

		mockResultRepository.EXPECT().Get(gomock.Any(), "test-execution-1").Return(testkube.Execution{}, mongo.ErrNoDocuments)
		testSuiteExecutionStatus := testkube.PASSED_TestSuiteExecutionStatus
		mockTestSuiteExecution := testkube.TestSuiteExecution{Id: "test-suite-execution-1", Status: &testSuiteExecutionStatus}
		mockTestResultRepository.EXPECT().Get(gomock.Any(), "test-suite-execution-1").Return(mockTestSuiteExecution, nil)

		statusKey1 := newStatusKey("testkube", "test-trigger-1")
		statusKey2 := newStatusKey("testkube", "test-trigger-2")
		triggerStatus1 := &triggerStatus{testExecutionIDs: []string{"test-execution-1"}}
		triggerStatus2 := &triggerStatus{testSuiteExecutionIDs: []string{"test-suite-execution-1"}}
		triggerStatusMap := map[statusKey]*triggerStatus{
			statusKey1: triggerStatus1,
			statusKey2: triggerStatus2,
		}
		s := &Service{
			triggerStatus:        triggerStatusMap,
			resultRepository:     mockResultRepository,
			testResultRepository: mockTestResultRepository,
			scraperInterval:      100 * time.Millisecond,
			logger:               log.DefaultLogger,
		}

		s.runExecutionScraper(ctx)

		for testTrigger, status := range s.triggerStatus {
			assert.Falsef(t, status.hasActiveTests(), "TestTrigger V1 %s should not have active tests", testTrigger)
		}
	})

	t.Run("active jobs", func(t *testing.T) {
		t.Parallel()

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		ctx, cancel := context.WithTimeout(context.Background(), 310*time.Millisecond)
		defer cancel()

		mockResultRepository := result.NewMockRepository(mockCtrl)
		mockTestResultRepository := testresult.NewMockRepository(mockCtrl)

		testSuiteExecutionStatus := testkube.RUNNING_TestSuiteExecutionStatus
		mockTestSuiteExecution := testkube.TestSuiteExecution{Id: "test-suite-execution-1", Status: &testSuiteExecutionStatus}
		mockTestResultRepository.EXPECT().Get(gomock.Any(), "test-suite-execution-1").Return(mockTestSuiteExecution, nil).Times(3)

		statusKey1 := newStatusKey("testkube", "test-trigger-1")
		triggerStatus1 := &triggerStatus{testSuiteExecutionIDs: []string{"test-suite-execution-1"}}
		triggerStatusMap := map[statusKey]*triggerStatus{
			statusKey1: triggerStatus1,
		}
		s := &Service{
			triggerStatus:        triggerStatusMap,
			resultRepository:     mockResultRepository,
			testResultRepository: mockTestResultRepository,
			scraperInterval:      100 * time.Millisecond,
			logger:               log.DefaultLogger,
		}

		s.runExecutionScraper(ctx)

		for testTrigger, status := range s.triggerStatus {
			assert.Truef(t, status.hasActiveTests(), "TestTrigger V1 %s should not have finished tests", testTrigger)
		}
	})
}
