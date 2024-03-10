package main

import (
	"encoding/json"
	"fmt"
	"os"
)

var IDToCategoryNodeMap = map[uint64]*CategoryNode{}
var categoryID uint64

// CategoryNode представляет собой узел дерева локаций
type CategoryNode struct {
	ID     uint64
	Name   string
	Parent *CategoryNode
}

// NewCategory Создает новый узел локации
func NewCategory(name string) *CategoryNode {
	categoryID++
	ptr := &CategoryNode{
		ID: categoryID, Name: name,
	}
	IDToCategoryNodeMap[ptr.ID] = ptr
	return ptr
}

// AddChild Добавляет дочернюю локацию к родительской категории
func (l *CategoryNode) AddChild(child *CategoryNode) {
	child.Parent = l
}

type JSONCategory struct {
	Name     string   `json:"name"`
	Children []string `json:"children"`
}

func BuildCategoryTreeFromFile(filename string) (*CategoryNode, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var categories []JSONCategory
	if err := json.NewDecoder(file).Decode(&categories); err != nil {
		return nil, err
	}

	// Create the root node
	rootNode := NewCategory("ROOT")

	// Map to store category nodes by name
	categoryNodeMap := make(map[string]*CategoryNode)

	// Iterate over JSON categories and construct the tree
	for _, jsonCategory := range categories {
		categoryNode := getCategoryNode(jsonCategory.Name, categoryNodeMap)
		for _, childName := range jsonCategory.Children {
			childNode := getCategoryNode(childName, categoryNodeMap)
			categoryNode.AddChild(childNode)
		}
		rootNode.AddChild(categoryNode)
	}

	return rootNode, nil
}

func getCategoryNode(name string, categoryNodeMap map[string]*CategoryNode) *CategoryNode {
	if node, ok := categoryNodeMap[name]; ok {
		return node
	}
	node := NewCategory(name)
	categoryNodeMap[name] = node
	return node
}

func PrintCategoryTree() {
	for ID, ptr := range IDToCategoryNodeMap {
		var curChildID []uint64
		for _, child := range IDToCategoryNodeMap {
			if child.Parent == ptr {
				curChildID = append(curChildID, child.ID)
			}
		}
		fmt.Printf("Node's %d name %s,  children: %v\n", ID, ptr.Name, curChildID)
	}
}
