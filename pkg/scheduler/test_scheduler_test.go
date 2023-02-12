package scheduler

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	k8sv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	v1 "github.com/kubeshop/testkube-operator/apis/executor/v1"
	testsv3 "github.com/kubeshop/testkube-operator/apis/tests/v3"
	executorsclientv1 "github.com/kubeshop/testkube-operator/client/executors/v1"
	testsclientv3 "github.com/kubeshop/testkube-operator/client/tests/v3"
	"github.com/kubeshop/testkube/pkg/api/v1/testkube"
	"github.com/kubeshop/testkube/pkg/executor/client"
	"github.com/kubeshop/testkube/pkg/log"
)

func TestParamsNilAssign(t *testing.T) {

	t.Run("merge two maps", func(t *testing.T) {

		p1 := map[string]testkube.Variable{"p1": testkube.NewBasicVariable("p1", "1")}
		p2 := map[string]testkube.Variable{"p2": testkube.NewBasicVariable("p2", "2")}

		out := mergeVariables(p1, p2)

		assert.Equal(t, 2, len(out))
		assert.Equal(t, "1", out["p1"].Value)
	})

	t.Run("merge two maps with override", func(t *testing.T) {

		p1 := map[string]testkube.Variable{"p1": testkube.NewBasicVariable("p1", "1")}
		p2 := map[string]testkube.Variable{"p1": testkube.NewBasicVariable("p1", "2")}

		out := mergeVariables(p1, p2)

		assert.Equal(t, 1, len(out))
		assert.Equal(t, "2", out["p1"].Value)
	})

	t.Run("merge with nil map", func(t *testing.T) {

		p2 := map[string]testkube.Variable{"p2": testkube.NewBasicVariable("p2", "2")}

		out := mergeVariables(nil, p2)

		assert.Equal(t, 1, len(out))
		assert.Equal(t, "2", out["p2"].Value)
	})

}

func TestGetExecuteOptions(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockTestsClient := testsclientv3.NewMockInterface(mockCtrl)
	mockExecutorsClient := executorsclientv1.NewMockInterface(mockCtrl)

	sc := Scheduler{testsClient: mockTestsClient, executorsClient: mockExecutorsClient, logger: log.DefaultLogger}

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
	mockExecutorTypes := "cypress"
	mockExecutor := v1.Executor{
		TypeMeta:   metav1.TypeMeta{},
		ObjectMeta: metav1.ObjectMeta{Namespace: "testkube", Name: "cypress"},
		Spec: v1.ExecutorSpec{
			Types:            []string{mockExecutorTypes},
			ExecutorType:     "job",
			URI:              "",
			Image:            "cypress",
			Args:             nil,
			Command:          []string{"run"},
			ImagePullSecrets: []k8sv1.LocalObjectReference{{Name: "secret-name1"}, {Name: "secret-name2"}},
			Features:         nil,
			ContentTypes:     nil,
			JobTemplate:      "",
			Meta:             nil,
		},
	}

	mockTestsClient.EXPECT().Get("id").Return(&mockTest, nil).Times(1)
	mockExecutorsClient.EXPECT().GetByType(mockExecutorTypes).Return(&mockExecutor, nil)

	req := testkube.ExecutionRequest{
		Name:             "id-1",
		Number:           1,
		ExecutionLabels:  map[string]string{"label": "value"},
		Namespace:        "namespace",
		VariablesFile:    "",
		Variables:        map[string]testkube.Variable{"var": testkube.Variable{Name: "one"}},
		Command:          []string{},
		Args:             []string{},
		Image:            "executor-image",
		ImagePullSecrets: []testkube.LocalObjectReference{},
		Envs: map[string]string{
			"env": "var",
		},
		SecretEnvs: map[string]string{
			"secretEnv": "secretVar",
		},
		Sync:                  false,
		HttpProxy:             "",
		HttpsProxy:            "",
		Uploads:               []string{},
		ActiveDeadlineSeconds: 10,
		ArtifactRequest:       &testkube.ArtifactRequest{},
		JobTemplate:           "",
		PreRunScript:          "",
		ScraperTemplate:       "",
	}

	got, err := sc.getExecuteOptions("namespace", "id", req)
	assert.NoError(t, err)

	want := client.ExecuteOptions{
		ID:                   "",
		TestName:             "id",
		Namespace:            "namespace",
		TestSpec:             mockTest.Spec,
		ExecutorName:         "cypress",
		ExecutorSpec:         mockExecutor.Spec,
		Request:              req,
		Sync:                 false,
		Labels:               map[string]string(nil),
		ImageOverride:        "executor-image",
		ImagePullSecretNames: []string{"secret-name1", "secret-name2"},
	}

	assert.Equal(t, want, got)
}
