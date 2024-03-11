package trees

import (
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"os"
	"slices"
)

var IDToLocationNodeMap = map[uint64]*LocationNode{}
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

type JSONLocation struct {
	Name     string   `json:"name"`
	Children []string `json:"children"`
}

func BuildLocationTreeFromFile(filename string) (*LocationNode, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			logrus.Error("Error with closing file:", err)
		}
	}(file)

	var locations []JSONLocation
	decoder := json.NewDecoder(file)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&locations); err != nil {
		return nil, err
	}

	// Create the root node
	rootNode := NewLocation("Все регионы")

	// Map to store location nodes by name
	locationNodeMap := make(map[string]*LocationNode)

	// Iterate over JSON locations and construct the tree
	for _, jsonLocation := range locations {
		locationNode := getLocationNode(jsonLocation.Name, locationNodeMap)
		for _, childName := range jsonLocation.Children {
			childNode := getLocationNode(childName, locationNodeMap)
			locationNode.AddChild(childNode)
		}
		rootNode.AddChild(locationNode)
	}

	return rootNode, nil
}

func getLocationNode(name string, locationNodeMap map[string]*LocationNode) *LocationNode {
	if node, ok := locationNodeMap[name]; ok {
		return node
	}
	node := NewLocation(name)
	locationNodeMap[name] = node
	return node
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
