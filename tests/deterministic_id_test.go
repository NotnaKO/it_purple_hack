package tests

import (
	"github.com/stretchr/testify/require"
	retrival "price_retrieval"
	"testing"
)

func TestDeterministicCategories(t *testing.T) {
	retrival.GetCategoriesTree()
	answer := retrival.ShowCategoryTree(false)
	for i := 0; i < 10; i++ {
		clear(retrival.IDToCategoryNodeMap)
		retrival.CategoryID = 0
		retrival.GetCategoriesTree()
		require.Equal(t, answer, retrival.ShowCategoryTree(false))
	}
}

func TestDeterministicLocation(t *testing.T) {
	retrival.GetLocationsTree()
	answer := retrival.ShowLocationTree(false)
	for i := 0; i < 5; i++ {
		retrival.LocationID = 0
		retrival.GetLocationsTree()
		require.Equal(t, answer, retrival.ShowLocationTree(false))
	}
}
