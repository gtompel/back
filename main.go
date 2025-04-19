package main

import (
	"fmt"
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"github.com/BurntSushi/toml"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var (
	localizer *i18n.Localizer
	db        *gorm.DB
)

func init() {
	bundle := i18n.NewBundle(language.Russian)
	bundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)
	bundle.MustLoadMessageFile("ru.toml")
	localizer = i18n.NewLocalizer(bundle, language.Russian.String())
}

func localize(id string) string {
	msg, err := localizer.Localize(&i18n.LocalizeConfig{
		MessageID: id,
	})
	if err != nil {
		fmt.Printf("Localization error for ID %s: %v\n", id, err)
		return id
	}
	return msg
}

func main() {
	a := app.New()
	w := a.NewWindow(localize("leave_management_app"))

	var err error
	db, err = gorm.Open(sqlite.Open("leave_management.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database: ", err)
	}

	err = db.AutoMigrate(&LeaveRequest{}, &Employee{})
	if err != nil {
		log.Fatal("Failed to migrate database: ", err)
	}

	// Начинаем с отображения таблицы отпусков
	showLeaveRequests(w)

	w.Resize(fyne.NewSize(800, 600))
	w.ShowAndRun()
}
