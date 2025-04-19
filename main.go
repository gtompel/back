package main

import (
	"fmt"
	"log"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/BurntSushi/toml"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
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

type Employee struct {
	ID         uint   `gorm:"primaryKey"`
	LastName   string `gorm:"not null"`
	FirstName  string `gorm:"not null"`
	MiddleName string
}

type LeaveRequest struct {
	ID         uint      `gorm:"primaryKey"`
	EmployeeID uint      `gorm:"not null"`
	Employee   Employee  `gorm:"foreignKey:EmployeeID"`
	Type       string    `gorm:"not null"`
	StartDate  time.Time `gorm:"not null"`
	EndDate    time.Time `gorm:"not null"`
	Reason     string    `gorm:"type:text"`
	Status     string    `gorm:"default:'pending'"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
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

func showLeaveRequests(w fyne.Window) {
	var results []struct {
		LastName  string
		FirstName string
		StartDate time.Time
		EndDate   time.Time
		Type      string
	}

	db.Table("leave_requests").
		Select("employees.last_name, employees.first_name, leave_requests.start_date, leave_requests.end_date, leave_requests.type").
		Joins("join employees on leave_requests.employee_id = employees.id").
		Scan(&results)

	data := [][]string{{localize("last_name"), localize("first_name"), localize("start_date"), localize("end_date"), localize("leave_type")}}
	for _, r := range results {
		data = append(data, []string{r.LastName, r.FirstName, r.StartDate.Format("02.01.2006"), r.EndDate.Format("02.01.2006"), r.Type})
	}

	table := widget.NewTable(
		func() (int, int) { return len(data), 5 },
		func() fyne.CanvasObject { return widget.NewLabel("template") },
		func(id widget.TableCellID, o fyne.CanvasObject) {
			label := o.(*widget.Label)
			label.SetText(data[id.Row][id.Col])
		},
	)

	// Добавляем кнопку "+" для добавления нового отпуска
	addButton := widget.NewButton("+", func() {
		showAddLeaveForm(w)
	})

	// Создаем контейнер с таблицей и кнопкой
	content := container.NewBorder(nil, container.NewHBox(addButton), nil, nil, table)

	w.SetContent(content)
}

func showAddLeaveForm(w fyne.Window) {
	lastName := widget.NewEntry()
	lastName.SetPlaceHolder(localize("last_name"))

	firstName := widget.NewEntry()
	firstName.SetPlaceHolder(localize("first_name"))

	middleName := widget.NewEntry()
	middleName.SetPlaceHolder(localize("middle_name"))

	// Создаем метку для placeholder
	selectedLeaveType := localize("select_leave_type")
	leaveTypeLabel := widget.NewLabel(selectedLeaveType)

	// Создаем PopUpMenu для выбора типа отпуска
	menu := fyne.NewMenu("",
		&fyne.MenuItem{Label: localize("paid"),
			Action: func() {
				selectedLeaveType := localize("paid")
				leaveTypeLabel.SetText(selectedLeaveType)
				fmt.Println(fmt.Sprintf(localize("selected_leave_type"), selectedLeaveType))
			}},
		&fyne.MenuItem{Label: localize("unpaid"),
			Action: func() {
				selectedLeaveType := localize("unpaid")
				leaveTypeLabel.SetText(selectedLeaveType)
				fmt.Println(fmt.Sprintf(localize("selected_leave_type"), selectedLeaveType))
			}},
		&fyne.MenuItem{Label: localize("sick"),
			Action: func() {
				selectedLeaveType := localize("sick")
				leaveTypeLabel.SetText(selectedLeaveType)
				fmt.Println(fmt.Sprintf(localize("selected_leave_type"), selectedLeaveType))
			}},
	)

	// Правильное объявление переменной leaveTypeButton
	var leaveTypeButton *widget.Button
	leaveTypeButton = widget.NewButton(selectedLeaveType, func() {
		// Получаем позицию кнопки относительно окна
		buttonPos := fyne.CurrentApp().Driver().AbsolutePositionForObject(leaveTypeButton)
		// Получаем размер кнопки
		buttonSize := leaveTypeButton.Size()

		// Вычисляем позицию для PopUpMenu (подя кнопкой)
		popUpPos := fyne.NewPos(buttonPos.X, buttonPos.Y+buttonSize.Height)

		widget.NewPopUpMenu(menu, w.Canvas()).ShowAtPosition(popUpPos)
	})

	startDate := widget.NewEntry()
	startDate.SetPlaceHolder(localize("start_date_format"))
	startDate.OnChanged = func(s string) {
		if len(s) != 10 || s[2] != '.' || s[5] != '.' {
			startDate.SetText(localize("invalid_date_format"))
		} else {
			_, err := time.Parse("02.01.2006", s)
			if err != nil {
				startDate.SetText(localize("invalid_date"))
			}
		}
	}

	endDate := widget.NewEntry()
	endDate.SetPlaceHolder(localize("end_date_format"))
	endDate.OnChanged = func(s string) {
		if len(s) != 10 || s[2] != '.' || s[5] != '.' {
			endDate.SetText(localize("invalid_date_format"))
		} else {
			_, err := time.Parse("02.01.2006", s)
			if err != nil {
				endDate.SetText(localize("invalid_date"))
			}
		}
	}

	reason := widget.NewMultiLineEntry()
	reason.SetPlaceHolder(localize("reason"))

	submitButton := widget.NewButton(localize("submit"), func() {
		fmt.Println("Submitting leave request...")

		start, err := time.Parse("02.01.2006", startDate.Text)
		if err != nil {
			fmt.Println("Error parsing start date: ", err)
			return
		}

		end, err := time.Parse("02.01.2006", endDate.Text)
		if err != nil {
			fmt.Println("Error parsing end date: ", err)
			return
		}

		employee := Employee{
			LastName:   lastName.Text,
			FirstName:  firstName.Text,
			MiddleName: middleName.Text,
		}

		result := db.Create(&employee)
		if result.Error != nil {
			fmt.Println("Failed to create employee: ", result.Error)
			return
		}

		leave := LeaveRequest{
			EmployeeID: employee.ID,
			Type:       leaveTypeLabel.Text,
			StartDate:  start,
			EndDate:    end,
			Reason:     reason.Text,
			Status:     "pending",
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		}

		result = db.Create(&leave)
		if result.Error != nil {
			fmt.Println("Failed to create leave request: ", result.Error)
			return
		}

		fmt.Println("Leave request submitted successfully!")
		showLeaveRequests(w)
	})

	// Кнопка отмены для возврата к таблице
	cancelButton := widget.NewButton(localize("cancel"), func() {
		showLeaveRequests(w)
	})

	content := container.NewVBox(
		widget.NewLabel(localize("last_name")+":"),
		lastName,
		widget.NewLabel(localize("first_name")+":"),
		firstName,
		widget.NewLabel(localize("middle_name")+":"),
		middleName,
		widget.NewLabel(localize("leave_type")+":"),
		leaveTypeButton,
		widget.NewLabel(localize("start_date")+":"),
		startDate,
		widget.NewLabel(localize("end_date")+":"),
		endDate,
		widget.NewLabel(localize("reason")+":"),
		reason,
		container.NewHBox(submitButton, cancelButton),
	)

	w.SetContent(content)
}
