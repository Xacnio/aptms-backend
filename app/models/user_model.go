package models

type User struct {
	ID            uint   `json:"id" validate:"min=1"`
	Email         string `json:"email" validate:"email"`
	Password      string `json:"-"`
	Name          string `json:"name"`
	Surname       string `json:"surname"`
	PhoneNumber   string `json:"phone_number" db:"phone_number"`
	Type          uint8  `json:"type" db:"type"`
	CreatedBy     uint   `json:"created_by" db:"created_by"`
	CreatedByName string `json:"created_by_name,omitempty" db:"created_by_name,omitempty"`
	Rank          string `json:"rank,omitempty" db:"rank,omitempty"`
}

type BuildingUser struct {
	User        User   `json:"user"`
	OwnedFlats  string `json:"owned_flats"`
	RentedFlats string `json:"rented_flats"`
}

type ManagerUserBuilding struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type ManagerUser struct {
	User      User                  `json:"user"`
	Buildings []ManagerUserBuilding `json:"buildings"`
}
