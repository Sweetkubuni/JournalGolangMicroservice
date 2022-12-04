package controller

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/Sweetkubuni/journal/api/models"
	"github.com/labstack/echo/v4"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type JournalHandlers struct {
	DB     *gorm.DB
	NextId int64
}

func NewJournalHandlers(DbHost, DbPort, DbUser, DbName, DbPassword string) (*JournalHandlers, error) {

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=America/Nassau", DbHost, DbUser, DbPassword, DbName, DbPort)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		fmt.Printf("Cannot connect to postgresql database")
		log.Fatal("This is the error:", err)
		return nil, err
	}

	db.Debug().AutoMigrate(&models.JournalEntry{})

	jh := &JournalHandlers{
		DB:     db,
		NextId: 1000,
	}

	return jh, nil
}

func (j *JournalHandlers) GetJournal(ctx echo.Context) error {

	journals, _ := models.FindAllJournals(j.DB)

	return ctx.JSON(http.StatusOK, journals)
}

func (j *JournalHandlers) PostJournal(ctx echo.Context) error {
	//TODO: get form data and save audio file to audio directory
	//store audio file path (/static/feifniweiewoifew.ogg to file)
	var journalEntry models.JournalEntry
	journalEntry.Prepare()

	// Read form fields
	journalEntry.Title = ctx.FormValue("title")

	file, err := ctx.FormFile("file")
	if err != nil {
		return err
	}
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	filename := "/audio/" + file.Filename
	// Destination
	dst, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer dst.Close()

	// Copy
	if _, err = io.Copy(dst, src); err != nil {
		return err
	}

	journalEntry.Audio = filename

	journalEntry.Save(j.DB)

	return nil
}
