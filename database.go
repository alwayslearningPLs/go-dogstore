package main

import (
	"context"
	"sync"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
  db *gorm.DB
  dbOnce sync.Once
)

func configureDB(dns string) {
  dbOnce.Do(func() {
    var err error
    db, err = gorm.Open(postgres.Open(dns), &gorm.Config{}) 
    if err != nil {
      panic(err)
    }
  })
}

func getInstance() *gorm.DB {
  return db
}

func getInstanceWithCtx(ctx context.Context) *gorm.DB {
  return getInstance().WithContext(ctx)
}

func getInstanceDry() *gorm.DB {
  return db.Session(&gorm.Session{DryRun: true})
}
