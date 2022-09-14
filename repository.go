package main

import (
	"context"
	"errors"
	"fmt"

	"gorm.io/gorm"
)

var (
	errOwnerNotFound = errors.New("owner not found")
	errBreedNotFound = errors.New("breed not found")
)

func addBreed(ctx context.Context, b *breed) error {
	return getInstanceWithCtx(ctx).Create(b).Error
}

func getBreed(ctx context.Context, b breed) ([]breed, error) {
	var result []breed

	tx := whereBreed(getInstanceWithCtx(ctx), b).Find(&result)

	return result, tx.Error
}

func delBreed(ctx context.Context, b breed) (rows int64, err error) {
	tx := whereBreed(getInstanceWithCtx(ctx), b).Delete(&b)
	return tx.RowsAffected, tx.Error
}

func whereBreed(db *gorm.DB, b breed) *gorm.DB {
	if b.ID != 0 {
		db = db.Where("breeds.id = ?", b.ID)
	}

	if b.Name != "" {
		db = db.Where("breeds.name = ?", b.Name)
	}

	if b.Description != "" {
		db = db.Where("breeds.description like ?", b.Description)
	}

	return db
}

func addOwner(ctx context.Context, o *owner) error {
	return getInstanceWithCtx(ctx).Create(o).Error
}

func getOwner(ctx context.Context, o owner) ([]owner, error) {
	var result []owner

	tx := whereOwner(getInstanceWithCtx(ctx), o).Find(&result)

	return result, tx.Error
}

func delOwner(ctx context.Context, o owner) (rows int64, err error) {
	tx := whereOwner(getInstanceWithCtx(ctx), o).Delete(&o)
	return tx.RowsAffected, tx.Error
}

func whereOwner(db *gorm.DB, o owner) *gorm.DB {
	if o.ID != "" {
		db = db.Where("owners.owner_id = ?", o.ID)
	}

	if o.Name != "" {
		db = db.Where("owners.name like ?", o.Name)
	}

	if o.Surname != "" {
		db = db.Where("owners.surname like ?", o.Surname)
	}

	return db
}

func addDog(ctx context.Context, p *dog) error {
	t := whereOwner(getInstanceWithCtx(ctx), p.Owner).Find(&p.Owner)
	if err := t.Error; err != nil {
		return err
	} else if p.Owner.Name == "" {
		return errOwnerNotFound
	}

	t = whereBreed(getInstanceWithCtx(ctx), p.Breed).Find(&p.Breed)
	if err := t.Error; err != nil {
		return err
	} else if p.Breed.ID == 0 {
		return errBreedNotFound
	}

	return getInstanceWithCtx(ctx).
		Exec("INSERT INTO dogs (name, breed_id, owner_id) VALUES (?, ?, ?)", p.Name, p.Breed.ID, p.Owner.ID).Error
}

func delDog(ctx context.Context, p dog) (rows int64, err error) {
	tx := getInstanceWithCtx(ctx).
		Exec("DELETE FROM dogs WHERE name = ? AND owner_id = ? AND breed_id IN (SELECT breed_id FROM breeds WHERE name = ?)",
			p.Name, p.Owner.ID, p.Breed.Name)

	return tx.RowsAffected, tx.Error
}

func getDog(ctx context.Context, p dog) ([]dog, error) {
	var result []dog

	tx := whereDog(getInstanceWithCtx(ctx), p).
		Joins("Owner").
		Joins("Breed")

	if p.Owner.ID != "" {
		tx = tx.Where("\"Owner\".owner_id = ?", p.Owner.ID)
	}

	if p.Breed.Name != "" {
		tx = tx.Where("\"Breed\".name = ?", p.Breed.Name)
	}

	tx = tx.Find(&result)

	if tx.DryRun {
		fmt.Println("sql: ", tx.Statement.SQL.String())
	}

	return result, tx.Error
}

func whereDog(db *gorm.DB, p dog) *gorm.DB {
	if p.ID != 0 {
		db = db.Where("dogs.id = ?", p.ID)
	}

	if p.Name != "" {
		db = db.Where("dogs.Name = ?", p.Name)
	}

	return db
}
