package plateauapi

import (
	"context"
	"testing"

	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
)

func TestMerger_Node(t *testing.T) {
	r1 := &r{
		nodes: []Node{
			n{ID: "1", Year: 2020, Name: "a"},
			n{ID: "2", Year: 2019},
			n{ID: "3", Year: 2021},
		},
	}
	r2 := &r{
		nodes: []Node{
			n{ID: "1", Year: 2020},
			n{ID: "2", Year: 2021},
			n{ID: "3", Year: 2020},
		},
	}
	r := NewMerger(r1, r2)

	t.Run("ok", func(t *testing.T) {
		res, err := r.Node(context.Background(), "1")
		assert.NoError(t, err)
		assert.Equal(t, n{ID: "1", Year: 2020, Name: "a"}, res)
	})

	t.Run("not found", func(t *testing.T) {
		res, err := r.Node(context.Background(), "4")
		assert.NoError(t, err)
		assert.Nil(t, res)
	})
}

func TestMerger_Nodes(t *testing.T) {
	r1 := &r{
		nodes: []Node{
			n{ID: "1", Year: 2020, Name: "a"},
			n{ID: "2", Year: 2019},
			n{ID: "3", Year: 2021},
		},
	}
	r2 := &r{
		nodes: []Node{
			n{ID: "1", Year: 2020},
			n{ID: "2", Year: 2021},
			n{ID: "3", Year: 2020},
		},
	}
	r := NewMerger(r1, r2)

	t.Run("ok", func(t *testing.T) {
		res, err := r.Nodes(context.Background(), []ID{"1", "2", "3"})
		assert.NoError(t, err)
		assert.Equal(t, []Node{
			n{ID: "1", Year: 2020, Name: "a"},
			n{ID: "2", Year: 2021},
			n{ID: "3", Year: 2021},
		}, res)
	})

	t.Run("not found", func(t *testing.T) {
		res, err := r.Nodes(context.Background(), []ID{"10", "11", "12"})
		assert.NoError(t, err)
		assert.Equal(t, []Node{nil, nil, nil}, res)
	})
}

func TestMergeResults(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		results1 := []*PlateauDatasetType{
			{ID: "1", Year: 2020, Name: "a"},
			{ID: "1", Year: 2020},
			{ID: "2", Year: 2019},
			{ID: "2", Year: 2021, Order: 10},
			{ID: "3", Year: 2021},
		}
		expected1 := []*PlateauDatasetType{
			{ID: "1", Year: 2020, Name: "a"},
			{ID: "3", Year: 2021},
			{ID: "2", Year: 2021, Order: 10},
		}
		res1 := mergeResults(results1)
		assert.Equal(t, expected1, res1)
	})

	t.Run("empty", func(t *testing.T) {
		results3 := []*PlateauDataset{}
		expected3 := []*PlateauDataset{}
		res3 := mergeResults(results3)
		assert.Equal(t, expected3, res3)
	})
}

func TestSortNodes(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		nodes := []n{}
		sortNodes(nodes)
		assert.Empty(t, nodes)
	})

	t.Run("single node", func(t *testing.T) {
		nodes := []n{
			{ID: "1", Year: 2020},
		}
		sortNodes(nodes)
		expected := []n{
			{ID: "1", Year: 2020},
		}
		assert.Equal(t, expected, nodes)
	})

	t.Run("multiple nodes", func(t *testing.T) {
		nodes := []n{
			{ID: "3", Year: 2020},
			{ID: "2", Year: 2019},
			{ID: "1", Year: 2020},
			{ID: "3", Year: 2021},
		}
		sortNodes(nodes)
		expected := []n{
			{ID: "1", Year: 2020},
			{ID: "2", Year: 2019},
			{ID: "3", Year: 2021},
			{ID: "3", Year: 2020},
		}
		assert.Equal(t, expected, nodes)
	})
}

func TestGetLatestYearNode(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		results := []PlateauDataset{}
		res := getLatestYearNode(results)
		assert.Equal(t, PlateauDataset{}, res)
	})

	t.Run("single node", func(t *testing.T) {
		results := []PlateauDataset{
			{ID: "1", Year: 2020},
		}
		res := getLatestYearNode(results)
		assert.Equal(t, PlateauDataset{ID: "1", Year: 2020}, res)
	})

	t.Run("duplicated nodes", func(t *testing.T) {
		results := []PlateauDataset{
			{ID: "1", Year: 2020, Name: "a"},
			{ID: "1", Year: 2020, Name: "b"},
		}
		res := getLatestYearNode(results)
		assert.Equal(t, PlateauDataset{ID: "1", Year: 2020, Name: "a"}, res)
	})

	t.Run("multiple nodes", func(t *testing.T) {
		dt := PlateauDatasetType{ID: "1", Year: 2021}
		assert.True(t, isPresent(dt))

		results := []Node{nil, dt}
		res := getLatestYearNode(results)
		assert.Equal(t, dt, res)
	})

	t.Run("multiple nodes 2", func(t *testing.T) {
		results := []PlateauDataset{
			{ID: "1", Year: 2020, Name: "a"},
			{ID: "2", Year: 2019},
			{ID: "3", Year: 2021},
		}
		res := getLatestYearNode(results)
		assert.Equal(t, PlateauDataset{ID: "3", Year: 2021}, res)
	})
}

func TestZip(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		res := zip([]int(nil))
		assert.Empty(t, res)
	})

	t.Run("singleSlice", func(t *testing.T) {
		input := []int{1, 2, 3}
		expected := [][]int{{1}, {2}, {3}}
		res := zip(input)
		assert.Equal(t, expected, res)
	})

	t.Run("multipleSlices", func(t *testing.T) {
		input1 := []int{1, 2, 3}
		input2 := []int{4, 5}
		input3 := []int{7, 8, 9}
		expected := [][]int{
			{1, 4, 7},
			{2, 5, 8},
			{3, 0, 9},
		}
		res := zip(input1, input2, input3)
		assert.Equal(t, expected, res)
	})
}

type n struct {
	ID    string
	Name  string
	Year  int
	Order int
}

func (n n) IsNode() {}

func (n n) GetID() ID {
	return ID(n.ID)
}

func (n n) GetYear() int {
	return n.Year
}

func (n n) GetOder() int {
	return n.Order
}

type r struct {
	nodes []Node
	Repo
}

func (r r) Node(ctx context.Context, id ID) (Node, error) {
	n, _ := lo.Find(r.nodes, func(n Node) bool {
		return n.GetID() == id
	})
	return n, nil
}

func (r r) Nodes(ctx context.Context, ids []ID) ([]Node, error) {
	res := lo.Map(ids, func(id ID, _ int) Node {
		n, _ := lo.Find(r.nodes, func(n Node) bool {
			return n.GetID() == id
		})
		return n
	})
	return res, nil
}
