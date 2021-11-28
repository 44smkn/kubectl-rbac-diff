package policy

import (
	"fmt"

	"github.com/44smkn/kubectl-role-diff/pkg/kubernetes"
	"github.com/44smkn/kubectl-role-diff/pkg/model"
)

var (
	defaultHeader = []string{"apiGroup", "resource", "create", "delete", "deletecollection", "get", "list", "patch", "update", "watch"}
)

type PolicyTable struct {
	Header   []string
	Contents []PolicyTableRow
}

func NewPolicyTable(policyTableRows []PolicyTableRow) *PolicyTable {
	return &PolicyTable{
		Header:   defaultHeader,
		Contents: policyTableRows,
	}
}

func (p *PolicyTable) RenderWithoutHeader() [][]string {
	table := make([][]string, 0)
	for _, row := range p.Contents {
		table = append(table, row.Render())
	}
	return table
}

type PolicyTableRow struct {
	apiGroup             string
	resource             string
	verbCreate           *bool
	verbDelete           *bool
	verbDeleteCollection *bool
	verbGet              *bool
	verbList             *bool
	verbPatch            *bool
	verbUpdate           *bool
	verbWatch            *bool
}

func NewPolicyTableRow(definition model.APIResource, declaredVerbs map[string]bool) *PolicyTableRow {
	row := PolicyTableRow{
		apiGroup: definition.Group,
		resource: definition.Name,
	}
	for _, v := range definition.Verbs {
		switch v {
		case "create":
			_, ok := declaredVerbs["create"]
			row.verbCreate = &ok
		case "delete":
			_, ok := declaredVerbs["delete"]
			row.verbDelete = &ok
		case "deletecollection":
			_, ok := declaredVerbs["deletecollection"]
			row.verbDeleteCollection = &ok
		case "get":
			_, ok := declaredVerbs["get"]
			row.verbGet = &ok
		case "list":
			_, ok := declaredVerbs["list"]
			row.verbGet = &ok
		case "patch":
			_, ok := declaredVerbs["list"]
			row.verbGet = &ok
		case "update":
			_, ok := declaredVerbs["update"]
			row.verbUpdate = &ok
		case "watch":
			_, ok := declaredVerbs["watch"]
			row.verbWatch = &ok
		default:
			// Log
		}
	}
	return &row
}

func (p *PolicyTableRow) Render() []string {
	return []string{
		p.apiGroup,
		p.resource,
		convVerbPermission(p.verbCreate),
		convVerbPermission(p.verbDelete),
		convVerbPermission(p.verbDeleteCollection),
		convVerbPermission(p.verbGet),
		convVerbPermission(p.verbList),
		convVerbPermission(p.verbPatch),
		convVerbPermission(p.verbUpdate),
		convVerbPermission(p.verbWatch),
	}
}

func convVerbPermission(verbPermission *bool) string {
	if verbPermission == nil {
		return "/"
	}
	if *verbPermission {
		return "OK"
	}
	return "NG"
}

type PolicyTableGenerator interface {
	Generate(manifest []byte) (*PolicyTable, error)
}

func NewPolicyTableGenerator(apiResources []model.APIResource) PolicyTableGenerator {
	return &defaultPolicyTableGenerator{apiResourcesDefinitions: apiResources}
}

type defaultPolicyTableGenerator struct {
	apiResourcesDefinitions []model.APIResource
}

func (g *defaultPolicyTableGenerator) Generate(manifest []byte) (*PolicyTable, error) {
	rows := make([]PolicyTableRow, 0, 100)
	rrp := kubernetes.NewRoleResourceParser(manifest)
	policies, err := rrp.Parse()
	if err != nil {
		return nil, fmt.Errorf("Failed to parse manifest: %w", err)
	}
	for _, d := range g.apiResourcesDefinitions {
		p := policies[d.APIGroupResource().String()]
		elem := NewPolicyTableRow(d, p.Verbs)
		rows = append(rows, *elem)
	}

	return NewPolicyTable(rows), nil
}
