package kubernetes

import (
	"fmt"

	"github.com/44smkn/kubectl-role-diff/pkg/model"
	"gopkg.in/yaml.v2"
	rbacv1 "k8s.io/api/rbac/v1"
	rbacutil "k8s.io/kubectl/pkg/util/rbac"
)

type RoleResourceParser interface {
	Parse() (map[string]model.PolicyRule, error)
}

func NewRoleResourceParser(manifest []byte) RoleResourceParser {
	return &defaultRoleResourceParser{
		manifest: manifest,
	}
}

type defaultRoleResourceParser struct {
	manifest []byte
}

func (d *defaultRoleResourceParser) Parse() (map[string]model.PolicyRule, error) {
	role := rbacv1.Role{}
	err := yaml.Unmarshal(d.manifest, &role)
	if err != nil {
		return nil, fmt.Errorf("Failed to unmarshal manifest: %w", err)
	}

	rawRules, err := rbacutil.CompactRules(role.Rules)
	if err != nil {
		return nil, fmt.Errorf("Failed to execute rules compaction: %w", err)
	}

	rules := make(map[string]model.PolicyRule)
	for _, r := range rawRules {
		policyRule := model.NewPolicyRule(r.Verbs, r.APIGroups, r.Resources, r.ResourceNames)
		apiGroupResources := policyRule.APIGroupResources()
		for _, agr := range apiGroupResources {
			rules[agr.String()] = policyRule
		}
	}

	return rules, nil
}
