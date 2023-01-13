package deployments

import (
	testtriggersv1 "github.com/kubeshop/testkube-operator/apis/testtriggers/v1"
	appsv1 "k8s.io/api/apps/v1"
)

// MapCRDConditionsToAPI maps Deployment CRD conditions to OpenAPI spec TestTriggerConditions
func MapCRDConditionsToAPI(conditions []appsv1.DeploymentCondition) []testtriggersv1.TestTriggerCondition {
	var results []testtriggersv1.TestTriggerCondition
	for _, condition := range conditions {
		results = append(results, testtriggersv1.TestTriggerCondition{
			Type_:  string(condition.Type),
			Status: (*testtriggersv1.TestTriggerConditionStatuses)(&condition.Status),
		})
	}

	return results
}
