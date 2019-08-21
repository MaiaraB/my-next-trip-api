package ssmodels

type Agent struct {
	ID                 int `json:"Id"`
	Name               string
	ImageURL           string `json:"ImageUrl"`
	Status             string
	OptimisedForMobile bool
	Type               string
}

func SearchAgentByID(list []Agent, id int) Agent {
	var idElement Agent
	for i := range list {
		currentID := list[i].ID
		if currentID == id {
			idElement = list[i]
			break
		}
	}
	return idElement
}
