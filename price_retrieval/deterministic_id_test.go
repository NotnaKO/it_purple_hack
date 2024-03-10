package main

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestDeterministicCategories(t *testing.T) {
	_, err := BuildCategoryTreeFromFile("../data/category_tree.json")
	require.NoError(t, err)
	answer := ShowCategoryTree(false)
	for i := 0; i < 10; i++ {
		clear(IDToCategoryNodeMap)
		CategoryID = 0
		_, err = BuildCategoryTreeFromFile("../data/category_tree.json")
		require.NoError(t, err)
		require.Equal(t, answer, ShowCategoryTree(false))
	}
}

func TestDeterministicLocation(t *testing.T) {
	_, err := BuildLocationTreeFromFile("../data/locations_tree.json")
	require.NoError(t, err)
	answer := ShowLocationTree(false)
	for i := 0; i < 5; i++ {
		LocationID = 0
		_, err := BuildLocationTreeFromFile("../data/locations_tree.json")
		require.NoError(t, err)
		require.Equal(t, answer, ShowLocationTree(false))
	}
}
