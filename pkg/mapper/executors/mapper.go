package executors

import (
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	executorv1 "github.com/kubeshop/testkube-operator/apis/executor/v1"
	"github.com/kubeshop/testkube/pkg/api/v1/testkube"
)

// MapCRDToAPI maps Executor CRD to OpenAPI spec Executor
func MapCRDToAPI(item executorv1.Executor) testkube.ExecutorUpsertRequest {
	return testkube.ExecutorUpsertRequest{
		Name:             item.Name,
		Namespace:        item.Namespace,
		Labels:           item.Labels,
		ExecutorType:     item.Spec.ExecutorType,
		Types:            item.Spec.Types,
		Uri:              item.Spec.URI,
		Image:            item.Spec.Image,
		ImagePullSecrets: mapImagePullSecretsToAPI(item.Spec.ImagePullSecrets),
		Command:          item.Spec.Command,
		Args:             item.Spec.Args,
		JobTemplate:      item.Spec.JobTemplate,
		Features:         MapFeaturesToAPI(item.Spec.Features),
		ContentTypes:     MapContentTypesToAPI(item.Spec.ContentTypes),
		Meta:             MapMetaToAPI(item.Spec.Meta),
	}
}

// MapAPIToCRD maps OpenAPI spec ExecutorUpsertRequest to CRD Executor
func MapAPIToCRD(request testkube.ExecutorUpsertRequest) executorv1.Executor {
	return executorv1.Executor{
		ObjectMeta: metav1.ObjectMeta{
			Name:      request.Name,
			Namespace: request.Namespace,
			Labels:    request.Labels,
		},
		Spec: executorv1.ExecutorSpec{
			ExecutorType:     request.ExecutorType,
			Types:            request.Types,
			URI:              request.Uri,
			Image:            request.Image,
			ImagePullSecrets: mapImagePullSecretsToCRD(request.ImagePullSecrets),
			Command:          request.Command,
			Args:             request.Args,
			JobTemplate:      request.JobTemplate,
			Features:         MapFeaturesToCRD(request.Features),
			ContentTypes:     MapContentTypesToCRD(request.ContentTypes),
			Meta:             MapMetaToCRD(request.Meta),
		},
	}
}

// MapExecutorCRDToExecutorDetails maps CRD Executor to OpemAPI spec ExecutorDetails
func MapExecutorCRDToExecutorDetails(item executorv1.Executor) testkube.ExecutorDetails {
	return testkube.ExecutorDetails{
		Name: item.Name,
		Executor: &testkube.Executor{
			ExecutorType:     item.Spec.ExecutorType,
			Image:            item.Spec.Image,
			ImagePullSecrets: mapImagePullSecretsToAPI(item.Spec.ImagePullSecrets),
			Command:          item.Spec.Command,
			Args:             item.Spec.Args,
			Types:            item.Spec.Types,
			Uri:              item.Spec.URI,
			JobTemplate:      item.Spec.JobTemplate,
			Labels:           item.Labels,
			Features:         MapFeaturesToAPI(item.Spec.Features),
			ContentTypes:     MapContentTypesToAPI(item.Spec.ContentTypes),
			Meta:             MapMetaToAPI(item.Spec.Meta),
		},
	}
}

func mapImagePullSecretsToCRD(secrets []testkube.LocalObjectReference) []v1.LocalObjectReference {
	var res []v1.LocalObjectReference
	for _, secret := range secrets {
		res = append(res, v1.LocalObjectReference{Name: secret.Name})
	}
	return res
}

func mapImagePullSecretsToAPI(secrets []v1.LocalObjectReference) []testkube.LocalObjectReference {
	var res []testkube.LocalObjectReference
	for _, secret := range secrets {
		res = append(res, testkube.LocalObjectReference{Name: secret.Name})
	}
	return res
}

func MapFeaturesToCRD(features []string) (out []executorv1.Feature) {
	for _, feature := range features {
		out = append(out, executorv1.Feature(feature))
	}
	return out
}

func MapFeaturesToAPI(features []executorv1.Feature) (out []string) {
	for _, feature := range features {
		out = append(out, string(feature))
	}
	return out
}

func MapContentTypesToCRD(contentTypes []string) (out []executorv1.ScriptContentType) {
	for _, contentType := range contentTypes {
		out = append(out, executorv1.ScriptContentType(contentType))
	}
	return out
}

func MapMetaToCRD(meta *testkube.ExecutorMeta) *executorv1.ExecutorMeta {
	if meta == nil {
		return nil
	}

	return &executorv1.ExecutorMeta{
		IconURI:  meta.IconURI,
		DocsURI:  meta.DocsURI,
		Tooltips: meta.Tooltips,
	}
}

func MapContentTypesToAPI(contentTypes []executorv1.ScriptContentType) (out []string) {
	for _, contentType := range contentTypes {
		out = append(out, string(contentType))
	}
	return out
}

func MapMetaToAPI(meta *executorv1.ExecutorMeta) *testkube.ExecutorMeta {
	if meta == nil {
		return nil
	}

	return &testkube.ExecutorMeta{
		IconURI:  meta.IconURI,
		DocsURI:  meta.DocsURI,
		Tooltips: meta.Tooltips,
	}
}

// MapUpdateToSpec maps ExecutorUpdateRequest to Executor CRD spec
func MapUpdateToSpec(request testkube.ExecutorUpdateRequest, executor *executorv1.Executor) *executorv1.Executor {
	var fields = []struct {
		source      *string
		destination *string
	}{
		{
			request.Name,
			&executor.Name,
		},
		{
			request.Namespace,
			&executor.Namespace,
		},
		{
			request.ExecutorType,
			&executor.Spec.ExecutorType,
		},
		{
			request.Image,
			&executor.Spec.Image,
		},
		{
			request.Uri,
			&executor.Spec.URI,
		},
		{
			request.JobTemplate,
			&executor.Spec.JobTemplate,
		},
	}

	for _, field := range fields {
		if field.source != nil {
			*field.destination = *field.source
		}
	}

	var slices = []struct {
		source      *[]string
		destination *[]string
	}{
		{
			request.Command,
			&executor.Spec.Command,
		},
		{
			request.Args,
			&executor.Spec.Args,
		},
		{
			request.Types,
			&executor.Spec.Types,
		},
	}

	for _, slice := range slices {
		if slice.source != nil {
			*slice.destination = *slice.source
		}
	}

	if request.Labels != nil {
		executor.Labels = *request.Labels
	}

	if request.ImagePullSecrets != nil {
		executor.Spec.ImagePullSecrets = mapImagePullSecretsToCRD(*request.ImagePullSecrets)
	}

	if request.Features != nil {
		executor.Spec.Features = MapFeaturesToCRD(*request.Features)
	}

	if request.ContentTypes != nil {
		executor.Spec.ContentTypes = MapContentTypesToCRD(*request.ContentTypes)
	}

	if request.Meta != nil {
		if (*request.Meta) == nil {
			executor.Spec.Meta = nil
			return executor
		}

		if (*request.Meta).IsEmpty() {
			executor.Spec.Meta = nil
			return executor
		}

		if executor.Spec.Meta == nil {
			executor.Spec.Meta = &executorv1.ExecutorMeta{}
		}

		if (*request.Meta).IconURI != nil {
			executor.Spec.Meta.IconURI = *(*request.Meta).IconURI
		}

		if (*request.Meta).DocsURI != nil {
			executor.Spec.Meta.DocsURI = *(*request.Meta).DocsURI
		}

		if (*request.Meta).Tooltips != nil {
			executor.Spec.Meta.Tooltips = *(*request.Meta).Tooltips
		}
	}

	return executor
}
