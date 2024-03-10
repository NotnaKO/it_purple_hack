package tests

import (
	"github.com/stretchr/testify/require"
	retrival "price_retrieval"
	"testing"
)

func TestDeterministicCategories(t *testing.T) {
	_, err := retrival.BuildCategoryTreeFromFile("../data/category_tree.json")
	require.NoError(t, err)
	answer := retrival.ShowCategoryTree(false)
	for i := 0; i < 10; i++ {
		clear(retrival.IDToCategoryNodeMap)
		retrival.CategoryID = 0
		_, err = retrival.BuildCategoryTreeFromFile("../data/category_tree.json")
		require.NoError(t, err)
		require.Equal(t, answer, retrival.ShowCategoryTree(false))
	}
}

func TestDeterministicLocation(t *testing.T) {
	_, err := retrival.BuildLocationTreeFromFile("../data/locations_tree.json")
	require.NoError(t, err)
	answer := retrival.ShowLocationTree(false)
	for i := 0; i < 5; i++ {
		retrival.LocationID = 0
		_, err := retrival.BuildLocationTreeFromFile("../data/locations_tree.json")
		require.NoError(t, err)
		require.Equal(t, answer, retrival.ShowLocationTree(false))
	}
}
