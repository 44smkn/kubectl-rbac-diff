package kubernetes

import (
	"fmt"

	"github.com/44smkn/kubectl-role-diff/pkg/model"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/discovery"
)

type ServerResourceFetcher interface {
	Fetch() ([]model.APIResource, error)
}

type defaultServerResourceFetcher struct {
	discovery.ServerResourcesInterface
}

func NewServerResourceFetcher(client discovery.ServerResourcesInterface) ServerResourceFetcher {
	return &defaultServerResourceFetcher{
		client,
	}
}

func (s *defaultServerResourceFetcher) Fetch() ([]model.APIResource, error) {
	lists, err := s.ServerPreferredResources()
	if err != nil {
		return nil, fmt.Errorf("Failed to fetch resources list from server: %w", err)
	}

	apiResources := make([]model.APIResource, 0, 100)
	for _, list := range lists {
		gv, err := schema.ParseGroupVersion(list.GroupVersion)
		if err != nil {
			return nil, fmt.Errorf("failed to parse GroupVersion: %w", err)
		}
		for _, resource := range list.APIResources {
			elem := model.APIResource{
				Name:  resource.Name,
				Group: gv.Group,
				Verbs: resource.Verbs,
			}
			apiResources = append(apiResources, elem)
		}
	}
	return apiResources, nil
}
