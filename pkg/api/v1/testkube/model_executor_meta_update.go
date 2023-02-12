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

// Executor meta update data
type ExecutorMetaUpdate struct {
	// URI for executor icon
	IconURI *string `json:"iconURI,omitempty"`
	// URI for executor docs
	DocsURI *string `json:"docsURI,omitempty"`
	// executor tooltips
	Tooltips *map[string]string `json:"tooltips,omitempty"`
}
