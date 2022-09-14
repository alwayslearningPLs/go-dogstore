package main

import "time"

type breed struct {
	ID          int    `gorm:"column:breed_id;primaryKey;not null" json:"-"`
	Name        string `gorm:"column:name;not null;unique" json:"name"`
	Description string `gorm:"column:description;not null" json:"description"`
}

type owner struct {
	ID       string    `gorm:"column:owner_id;primaryKey;default:uuid_generate_v4()" json:"id"`
	Name     string    `gorm:"column:name;not null" json:"name"`
	Surname  string    `gorm:"column:surname;not null" json:"surname"`
	Birthday time.Time `gorm:"column:birthday;not null" json:"birthday"`
}

type dog struct {
	ID   int    `gorm:"column:dog_id;primaryKey;not null" json:"-"`
	Name string `gorm:"column:name;not null;index:dogs_name_owner_id_key,unique" json:"name"`

	BreedID int    `json:"-"`
	Breed   breed  `json:"breed"`
	OwnerID string `gorm:"index:dogs_name_owner_id_key,unique" json:"-"`
	Owner   owner  `json:"owner"`
}

type ResponseWrap struct {
	Code int    `json:"code"`
	Msg  string `json:"message"`
}
