package v2alpha1

import "k8s.io/apimachinery/pkg/runtime/schema"

func init() {
	SchemeBuilder.Register(&Integration{}, &IntegrationList{})
}

func Resource(resource string) schema.GroupResource {
	return SchemeGroupVersion.WithResource(resource).GroupResource()
}
