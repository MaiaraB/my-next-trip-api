package ssmodels

type Carrier struct {
	ID          int `json:"Id"`
	Code        string
	Name        string
	ImageURL    string `json:"ImageUrl"`
	DisplayCode string
}

func SearchCarrierByID(list []Carrier, id int) Carrier {
	var idElement Carrier
	for i := range list {
		currentID := list[i].ID
		if currentID == id {
			idElement = list[i]
			break
		}
	}
	return idElement
}
