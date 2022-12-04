package models

import (
	"errors"
	"html"
	"strings"
	"time"

	"gorm.io/gorm"
)

type JournalEntry struct {
	gorm.Model
	Journal
}

func (j *JournalEntry) Prepare() {
	j.Title = html.EscapeString(strings.TrimSpace(j.Title))
	j.CreatedAt = time.Now()
	j.UpdatedAt = time.Now()
}

func (j *JournalEntry) Validate(action string) error {
	if j.Title == "" {
		return errors.New("required Title")
	}

	return nil
}

func (j *JournalEntry) Save(db *gorm.DB) (*JournalEntry, error) {
	err := db.Debug().Create(&j).Error
	if err != nil {
		return &JournalEntry{}, err
	}
	return j, nil
}

func FindAllJournals(db *gorm.DB) (*[]JournalEntry, error) {
	var err error
	journals := []JournalEntry{}
	if err != nil {
		return &[]JournalEntry{}, err
	}
	return &journals, err
}

func FindJournalByID(db *gorm.DB, uid uint32) (*JournalEntry, error) {
	var err error
	var j JournalEntry
	err = db.Debug().Model(JournalEntry{}).Where("id = ?", uid).Take(&j).Error
	if err != nil {
		return &JournalEntry{}, err
	}

	if gorm.ErrRecordNotFound == err {
		return &JournalEntry{}, errors.New("Journal Not Found")
	}
	return &j, err
}
