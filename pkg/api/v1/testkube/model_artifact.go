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

// API server artifact
type Artifact struct {
	// artifact file path
	Name string `json:"name,omitempty"`
	// file size in bytes
	Size int32 `json:"size,omitempty"`
	// execution name that produced the artifact
	ExecutionName string `json:"executionName,omitempty"`
}
