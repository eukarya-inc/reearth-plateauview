package plateauapi

import (
	"context"
	"fmt"
	"sort"

	"github.com/reearth/reearthx/util"
	"github.com/samber/lo"
)

// Merger is a repository that merges multiple repositories. It returns the latest year data.
type Merger struct {
	repos []Repo
}

func NewMerger(repos ...Repo) *Merger {
	return &Merger{repos: repos}
}

var _ Repo = (*Merger)(nil)

func (m *Merger) Name() string {
	return fmt.Sprintf("merger(%s)", lo.Map(m.repos, func(r Repo, _ int) string {
		if r == nil {
			return "nil"
		}
		return r.Name()
	}))
}

func (m *Merger) Node(ctx context.Context, id ID) (Node, error) {
	nodes, err := getRepoResults(m.repos, func(r Repo) (Node, error) {
		return r.Node(ctx, id)
	})
	if err != nil {
		return nil, err
	}

	return getLatestYearNode(nodes), nil
}

func (m *Merger) Nodes(ctx context.Context, ids []ID) ([]Node, error) {
	nodes, err := getRepoResults(m.repos, func(r Repo) ([]Node, error) {
		return r.Nodes(ctx, ids)
	})
	if err != nil {
		return nil, err
	}

	res := lo.Map(zip(nodes...), func(n []Node, _ int) Node {
		return getLatestYearNode(n)
	})
	return res, nil
}

func (m *Merger) Area(ctx context.Context, code AreaCode) (Area, error) {
	areas, err := getRepoResults(m.repos, func(r Repo) (Area, error) {
		return r.Area(ctx, code)
	})
	if err != nil {
		return nil, err
	}

	return getLatestYearNode(areas), nil
}

func (m *Merger) Areas(ctx context.Context, input *AreasInput) ([]Area, error) {
	areas, err := getFlattenRepoResults(m.repos, func(r Repo) ([]Area, error) {
		return r.Areas(ctx, input)
	})
	if err != nil {
		return nil, err
	}

	return mergeResults(areas), nil
}

func (m *Merger) DatasetTypes(ctx context.Context, input *DatasetTypesInput) ([]DatasetType, error) {
	dts, err := getFlattenRepoResults(m.repos, func(r Repo) ([]DatasetType, error) {
		return r.DatasetTypes(ctx, input)
	})
	if err != nil {
		return nil, err
	}

	return mergeResults(dts), nil
}

func (m *Merger) Datasets(ctx context.Context, input *DatasetsInput) ([]Dataset, error) {
	datasets, err := getFlattenRepoResults(m.repos, func(r Repo) ([]Dataset, error) {
		return r.Datasets(ctx, input)
	})
	if err != nil {
		return nil, err
	}

	return mergeResults(datasets), nil
}

func (m *Merger) PlateauSpecs(ctx context.Context) ([]*PlateauSpec, error) {
	res, err := getFlattenRepoResults(m.repos, func(r Repo) ([]*PlateauSpec, error) {
		return r.PlateauSpecs(ctx)
	})
	if err != nil {
		return nil, err
	}

	sortNodes(res)
	return res, nil
}

func (m *Merger) Years(ctx context.Context) ([]int, error) {
	years, err := getRepoResults(m.repos, func(r Repo) ([]int, error) {
		return r.Years(ctx)
	})
	if err != nil {
		return nil, err
	}

	res := lo.Uniq(lo.Flatten(years))
	sort.Ints(res)
	return res, nil
}

func getRepoResults[T any](repos []Repo, f func(r Repo) (T, error)) ([]T, error) {
	return util.TryMap(repos, func(r Repo) (_ T, _ error) {
		if r == nil {
			return
		}
		res, err := f(r)
		if err != nil {
			return res, fmt.Errorf("repo %s: %w", r.Name(), err)
		}
		return res, nil
	})
}

func getFlattenRepoResults[T any](repos []Repo, f func(r Repo) ([]T, error)) ([]T, error) {
	res, err := getRepoResults(repos, f)
	if err != nil {
		return nil, err
	}

	return lo.Flatten(res), nil
}

func mergeResults[T IDNode](results []T) []T {
	groups := lo.GroupBy(results, func(n T) string {
		return string(n.GetID())
	})

	res := make([]T, 0, len(groups))
	for _, g := range groups {
		res = append(res, getLatestYearNode(g))
	}

	sortNodes(res)
	return res
}

func sortNodes[T IDNode](nodes []T) {
	sort.Slice(nodes, func(i, j int) bool {
		return nodes[i].GetID() < nodes[j].GetID() || (nodes[i].GetID() == nodes[j].GetID() && getYear(nodes[i]) > getYear(nodes[j]))
	})
}

func getLatestYearNode[T any](results []T) T {
	results = lo.Filter(results, func(a T, _ int) bool {
		return isNodePresent(a)
	})
	return lo.MaxBy(results, func(a, b T) bool {
		return getYear(a) > getYear(b)
	})
}

func zip[T any](a ...[]T) [][]T {
	if len(a) == 0 {
		return nil
	}
	res := make([][]T, len(a[0]))
	for i := range res {
		res[i] = make([]T, len(a))
	}
	for i, aa := range a {
		for j, v := range aa {
			res[j][i] = v
		}
	}
	return res
}

type IDNode interface {
	GetID() ID
}

type YearNode interface {
	GetYear() int
}

func isNodePresent(n any) bool {
	if n, ok := n.(Node); ok {
		return n != nil
	}
	return false
}

func getYear(n any) int {
	if yn, ok := n.(YearNode); ok {
		return yn.GetYear()
	}
	return 0
}
