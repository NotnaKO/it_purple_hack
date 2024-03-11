package trees

import (
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"os"
	"slices"
)

var IDToCategoryNodeMap = map[uint64]*CategoryNode{}
var CategoryID uint64

// CategoryNode представляет собой узел дерева локаций
type CategoryNode struct {
	ID     uint64
	Name   string
	Parent *CategoryNode
}

// NewCategory Создает новый узел локации
func NewCategory(name string) *CategoryNode {
	CategoryID++
	ptr := &CategoryNode{
		ID: CategoryID, Name: name,
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
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			logrus.Error(err)
		}
	}(file)

	var categories []JSONCategory
	decoder := json.NewDecoder(file)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&categories); err != nil {
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

func ShowCategoryTree(printNeed bool) []string {
	var answer []string
	for ID, ptr := range IDToCategoryNodeMap {
		var curChildID []uint64
		for _, child := range IDToCategoryNodeMap {
			if child.Parent == ptr {
				curChildID = append(curChildID, child.ID)
			}
		}
		slices.Sort(curChildID)
		answer = append(answer, fmt.Sprintf("Node's %d name %s,  children: %v", ID, ptr.Name, curChildID))
	}
	slices.Sort(answer)
	if printNeed {
		for _, s := range answer {
			fmt.Println(s)
		}
	}
	return answer
}
