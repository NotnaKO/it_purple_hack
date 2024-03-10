package price_retrival

import (
	"fmt"
	"slices"
)

var IDToCategoryNodeMap = map[uint64]*CategoryNode{}

func GetCategoriesTree() *CategoryNode {
	// Создаем корневую категорию - ROOT
	rootNode := NewCategory("ROOT")

	for _, item := range rawCategories {
		categoryNode := NewCategory(item.head)

		for _, subCategory := range item.children {
			subCategoryNode := NewCategory(subCategory)
			categoryNode.AddChild(subCategoryNode)
		}

		rootNode.AddChild(categoryNode)
	}

	return rootNode
}

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

type rawCategoryItem struct {
	head     string
	children []string
}
