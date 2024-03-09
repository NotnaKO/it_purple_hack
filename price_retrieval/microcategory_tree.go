package main

import "fmt"

var IDToCategoryNodeMap = map[uint64]*CategoryNode{}

func GetCategoriesTree() *CategoryNode {
	// Создаем корневую категорию - ROOT
	rootNode := NewCategory("ROOT")

	for category, subCategories := range rawCategories {
		categoryNode := NewCategory(category)

		for _, subCategory := range subCategories {
			subCategoryNode := NewCategory(subCategory)
			categoryNode.AddChild(subCategoryNode)
		}

		rootNode.AddChild(categoryNode)
	}

	return rootNode
}

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

var rawCategories = map[string][]string{
	"Бытовая электроника":           {"Товары для компьютера", "Фототехника", "Телефоны", "Планшеты и электронные книги", "Оргтехника и расходники", "Ноутбуки", "Настольные компьютеры", "Игры, приставки и программы", "Аудио и видео"},
	"Готовый бизнес и оборудование": {"Готовый бизнес", "Оборудование для бизнеса"},
	"Для дома и дачи":               {"Мебель и интерьер", "Ремонт и строительство", "Продукты питания", "Растения", "Бытовая техника", "Посуда и товары для кухни"},
	"Животные":                      {"Другие животные", "Товары для животных", "Птицы", "Аквариум", "Кошки", "Собаки"},
	"Личные вещи":                   {"Детская одежда и обувь", "Одежда, обувь, аксессуары", "Товары для детей и игрушки", "Часы и украшения", "Красота и здоровье"},
	"Недвижимость":                  {"Недвижимость за рубежом", "Квартиры", "Коммерческая недвижимость", "Гаражи и машиноместа", "Земельные участки", "Дома, дачи, коттеджи", "Комнаты"},
	"Работа":                        {"Резюме", "Вакансии"},
	"Транспорт":                     {"Автомобили", "Запчасти и аксессуары", "Грузовики и спецтехника", "Водный транспорт", "Мотоциклы и мототехника"},
	"Услуги":                        {"Предложения услуг"},
	"Хобби и отдых":                 {"Охота и рыбалка", "Спорт и отдых", "Коллекционирование", "Книги и журналы", "Велосипеды", "Музыкальные инструменты", "Билеты и путешествия"},
}
