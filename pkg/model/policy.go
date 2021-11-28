package model

import "fmt"

type APIResource struct {
	Name  string
	Group string
	Verbs []string
}

type PolicyRule struct {
	Verbs         map[string]bool
	APIGroups     []string
	Resources     []string
	ResourceNames []string
}

type APIGroupResource struct {
	APIGroup string
	Resource string
}

func NewPolicyRule(verbs, apiGroups, resouces, resourceNames []string) PolicyRule {
	if isCoreResource(apiGroups) {
		apiGroups = append(apiGroups, "")
	}
	return PolicyRule{
		Verbs:         convertToSet(verbs),
		APIGroups:     apiGroups,
		Resources:     resouces,
		ResourceNames: resourceNames,
	}
}

func isCoreResource(apiGroups []string) bool {
	return len(apiGroups) == 0
}

func (r *PolicyRule) APIGroupResources() []APIGroupResource {
	apiGroupResources := make([]APIGroupResource, 0)
	for _, apiGroup := range r.APIGroups {
		for _, resource := range r.Resources {
			elem := APIGroupResource{
				APIGroup: apiGroup,
				Resource: resource,
			}
			apiGroupResources = append(apiGroupResources, elem)
		}
	}
	return apiGroupResources
}

func (ar *APIGroupResource) String() string {
	return fmt.Sprintf("%s/%s", ar.APIGroup, ar.Resource)
}

func (r *APIResource) APIGroupResource() *APIGroupResource {
	return &APIGroupResource{
		APIGroup: r.Group,
		Resource: r.Name,
	}
}
