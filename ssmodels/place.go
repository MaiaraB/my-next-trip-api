package ssmodels

type Place struct {
	ID       int `json:"Id"`
	ParentID int `json:"ParentId"`
	Code     string
	Type     string
	Name     string
}

func SearchPlaceByID(list []Place, id int) Place {
	var idElement Place
	for i := range list {
		currentID := list[i].ID
		if currentID == id {
			idElement = list[i]
			break
		}
	}
	return idElement
}
