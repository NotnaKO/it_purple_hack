package price_retrival

import (
	"fmt"
	"slices"
)

var IDToLocationNodeMap = map[uint64]*LocationNode{}

func GetLocationsTree() *LocationNode {
	// Создаем корневую локацию - Все регионы
	allRegions := NewLocation("Все регионы")

	for _, item := range rawLocations {
		regionNode := NewLocation(item.head)

		for _, city := range item.children {
			cityNode := NewLocation(city)
			regionNode.AddChild(cityNode)
		}
		allRegions.AddChild(regionNode)
	}

	return allRegions
}

var LocationID uint64

// LocationNode представляет собой узел дерева локаций
type LocationNode struct {
	ID     uint64
	Name   string
	Parent *LocationNode
}

// NewLocation Создает новый узел локации
// Придерживается контракта, что такого узла не было
func NewLocation(name string) *LocationNode {
	LocationID++
	ptr := &LocationNode{
		ID:   LocationID,
		Name: name,
	}
	IDToLocationNodeMap[ptr.ID] = ptr
	return ptr
}

// AddChild Добавляет дочернюю локацию к родительской локации
func (l *LocationNode) AddChild(child *LocationNode) {
	child.Parent = l
}

// PrintLocationTree Выводит дерево локаций;
// предполагается, что не используется часто
func ShowLocationTree(printNeed bool) []string {
	var answer []string
	for ID, ptr := range IDToLocationNodeMap {
		var curChildID []uint64
		for _, child := range IDToLocationNodeMap {
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

type rawLocationItem struct {
	head     string
	children []string
}
