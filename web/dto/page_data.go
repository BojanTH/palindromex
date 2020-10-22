package dto

// PageData contains all data necessary for the view
type PageData map[string]interface{}

func (data PageData) SetFlashes(flashes map[string][]string) {
	data["Flashes"] = flashes
}