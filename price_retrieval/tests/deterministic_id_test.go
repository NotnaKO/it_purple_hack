package tests

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"os"
	"path"
	"testing"
	"trees"
)

func TestDeterministicCategories(t *testing.T) {
	rootPath, err := os.Getwd()
	require.NoError(t, err)
	pathToSource := path.Join(rootPath, "../../data/category_tree.json")
	fmt.Println(pathToSource)
	_, err = trees.BuildCategoryTreeFromFile(pathToSource)
	require.NoError(t, err)
	answer := trees.ShowCategoryTree(false)
	for i := 0; i < 10; i++ {
		clear(trees.IDToCategoryNodeMap)
		trees.CategoryID = 0
		_, err = trees.BuildCategoryTreeFromFile(pathToSource)
		require.NoError(t, err)
		require.Equal(t, answer, trees.ShowCategoryTree(false))
	}
}

func TestDeterministicLocation(t *testing.T) {
	rootPath, err := os.Getwd()
	require.NoError(t, err)
	pathToSource := path.Join(rootPath, "../../data/locations_tree.json")
	_, err = trees.BuildLocationTreeFromFile(pathToSource)
	require.NoError(t, err)
	answer := trees.ShowLocationTree(false)
	for i := 0; i < 5; i++ {
		trees.LocationID = 0
		_, err := trees.BuildLocationTreeFromFile(pathToSource)
		require.NoError(t, err)
		require.Equal(t, answer, trees.ShowLocationTree(false))
	}
}
