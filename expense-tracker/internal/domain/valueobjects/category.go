package valueobjects

type Category string

const (
	Groceries   Category = "groceries"
	Leisure     Category = "leisure"
	Electronics Category = "electronics"
	Utilities   Category = "utilities"
	Clothing    Category = "clothing"
	Health      Category = "health"
	Others      Category = "others"
)

func (c Category) IsValid() bool {
	switch c {
	case Groceries, Leisure, Electronics, Utilities, Clothing, Health, Others:
		return true
	default:
		return false
	}
}

func GetAllCategories() []Category {
	return []Category{Groceries, Leisure, Electronics, Utilities, Clothing, Health, Others}
}