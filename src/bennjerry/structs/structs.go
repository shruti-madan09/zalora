package structs

type IceCreamDataStruct struct {
	AllergyInfo           string   `json:"allergy_info"`
	Description           string   `json:"description"`
	DietaryCertifications string   `json:"dietary_certifications"`
	ImageClosed           string   `json:"image_closed"`
	ImageOpened           string   `json:"image_open"`
	Name                  string   `json:"name"`
	ProductId             string   `json:"productId"`
	Story                 string   `json:"story"`
	Id                    int      `json:"id"`
	SourcingValues        []string `json:"sourcing_values"`
	Ingredients           []string `json:"ingredients"`
}

type CreateUpdateDeleteResponse struct {
	Message string `json:"message"`
	Id      int    `json:"id"`
	Success bool   `json:"success"`
}

type ReadResponse struct {
	Message string              `json:"message"`
	Success bool                `json:"success"`
	Data    *IceCreamDataStruct `json:"data"`
}
