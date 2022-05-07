package models

import (
	"database/sql"
	"time"
)

type Building struct {
	ID          int      `json:"id"`
	Name        string   `json:"name"`
	City        City     `json:"city" db:",prefix=cities."`
	District    District `json:"district" db:",prefix=districts."`
	PhoneNumber string   `json:"phoneNumber"  db:"phone_number"`
	TaxNumber   string   `json:"taxNumber"  db:"tax_number"`
	Address     string   `json:"address"`
}

type SystemInfo struct {
	TotalBuildings uint32 `json:"total_buildings" db:"total_buildings"`
	TotalUsers     uint32 `json:"total_users" db:"total_users"`
	TotalFlats     uint32 `json:"total_flats" db:"total_flats"`
}

type BuildingDetails struct {
	CashAmount   float64 `json:"cash_amount" db:"cash_amount"`
	CityName     string  `json:"city_name" db:"city_name"`
	DistrictName string  `json:"district_name" db:"district_name"`
	BlockCount   uint8   `json:"block_count" db:"block_count"`
	FlatCount    uint8   `json:"flat_count" db:"flat_count"`
	UserCount    uint8   `json:"user_count" db:"user_count"`
}

type Block struct {
	ID         int    `json:"id"`
	BuildingID int    `json:"building_id" db:"building_id"`
	Letter     string `json:"letter"`
	DNumber    string `json:"d_number" db:"d_number"`
}

type Flat struct {
	ID          int    `json:"id"`
	Type        uint8  `json:"type"`
	BuildingID  int    `json:"-"`
	BlockID     int    `json:"block_id"`
	BlockLetter string `json:"block_letter"`
	Number      string `json:"number"`
	OwnerID     int    `json:"owner_id"`
	TenantID    int    `json:"tenant_id"`
	OwnerName   string `json:"owner_name"`
	TenantName  string `json:"tenant_name"`
}

func (t *Flat) Default() {
	t.OwnerName = "Yok"
	t.TenantName = "Yok"
}

type Paid struct {
	Type   uint8      `json:"type" db:"paid_type"`
	Time   *time.Time `json:"time" db:"paid_time"`
	Status bool       `json:"status" db:"paid_status"`
}

type Payer struct {
	FullName string `json:"full_name" db:"payer_full_name"`
	Email    string `json:"email" db:"payer_email"`
	Phone    string `json:"phone" db:"payer_phone"`
}

type Revenue struct {
	ID         int           `json:"id"`
	BuildingID int           `json:"building_id" db:"building_id"`
	FlatID     sql.NullInt32 `json:"flat_id" db:"flat_id"`
	FlatName   string        `json:"flat_name" db:"flat_name"`
	RID        int           `json:"rid"`
	Total      float64       `json:"total"`
	Time       *time.Time    `json:"time"`
	Paid       Paid          `json:"paid" db:",prefix="`
	Payer      Payer         `json:"payer" db:",prefix="`
	Details    string        `json:"details" db:"details"`
}

type ExpenseTo struct {
	FullName string `json:"full_name" db:"to_name"`
	Email    string `json:"email" db:"to_email"`
	Phone    string `json:"phone" db:"to_phone"`
}

type Expense struct {
	ID         int        `json:"id"`
	BuildingID int        `json:"building_id" db:"building_id"`
	EID        int        `json:"eid"`
	Total      float64    `json:"total"`
	Time       *time.Time `json:"time"`
	Paid       Paid       `json:"paid" db:",prefix="`
	To         ExpenseTo  `json:"to" db:",prefix="`
	Details    string     `json:"details" db:"details"`
}
