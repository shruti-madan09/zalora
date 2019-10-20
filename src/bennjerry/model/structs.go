package model

// Used to define schema of table product
type Product struct {
	Allergy                string
	Description            string
	ImageClosed            string
	ImageOpened            string
	Name                   string
	ProductId              string
	Story                  string
	DietaryCertificationId int
	Id                     int
	IsInActive             int8
}

// Used to define schema of tables sourcingvalue, ingredient, dietarycertification
type Property struct {
	Id   int
	Name string
}

// Used for relation table of product and its property (e.g. sourcingvalue, ingredient, dietarycertification)
type ProductProperty struct {
	ProductId    int
	PropertyId   int
	PropertyName string
}
