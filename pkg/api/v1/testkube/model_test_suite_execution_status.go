/*
 * Testkube API
 *
 * Testkube provides a Kubernetes-native framework for test definition, execution and results
 *
 * API version: 1.0.0
 * Contact: testkube@kubeshop.io
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */
package testkube

type TestSuiteExecutionStatus string

// List of TestSuiteExecutionStatus
const (
	QUEUED_TestSuiteExecutionStatus   TestSuiteExecutionStatus = "queued"
	RUNNING_TestSuiteExecutionStatus  TestSuiteExecutionStatus = "running"
	PASSED_TestSuiteExecutionStatus   TestSuiteExecutionStatus = "passed"
	FAILED_TestSuiteExecutionStatus   TestSuiteExecutionStatus = "failed"
	ABORTING_TestSuiteExecutionStatus TestSuiteExecutionStatus = "aborting"
	ABORTED_TestSuiteExecutionStatus  TestSuiteExecutionStatus = "aborted"
	TIMEOUT_TestSuiteExecutionStatus  TestSuiteExecutionStatus = "timeout"
)
