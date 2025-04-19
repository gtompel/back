package main

import (
	"fmt"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

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
		data = append(data, []string{
			r.LastName,
			r.FirstName,
			r.StartDate.Format("02.01.2006"),
			r.EndDate.Format("02.01.2006"),
			r.Type,
		})
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
		&fyne.MenuItem{
			Label: localize("paid"),
			Action: func() {
				selectedLeaveType = localize("paid")
				leaveTypeLabel.SetText(selectedLeaveType)
				fmt.Println(fmt.Sprintf(localize("selected_leave_type"), selectedLeaveType))
			},
		},
		&fyne.MenuItem{
			Label: localize("unpaid"),
			Action: func() {
				selectedLeaveType = localize("unpaid")
				leaveTypeLabel.SetText(selectedLeaveType)
				fmt.Println(fmt.Sprintf(localize("selected_leave_type"), selectedLeaveType))
			},
		},
		&fyne.MenuItem{
			Label: localize("sick"),
			Action: func() {
				selectedLeaveType = localize("sick")
				leaveTypeLabel.SetText(selectedLeaveType)
				fmt.Println(fmt.Sprintf(localize("selected_leave_type"), selectedLeaveType))
			},
		},
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

	// Создаем поля для ввода дат без немедленной валидации
	startDate := widget.NewEntry()
	startDate.SetPlaceHolder(localize("start_date_format"))

	endDate := widget.NewEntry()
	endDate.SetPlaceHolder(localize("end_date_format"))

	reason := widget.NewMultiLineEntry()
	reason.SetPlaceHolder(localize("reason"))

	submitButton := widget.NewButton(localize("submit"), func() {
		fmt.Println("Submitting leave request...")

		// Проверяем формат и валидность дат только при отправке формы
		isValid := true

		// Проверка начальной даты
		if startDate.Text == "" {
			startDate.SetText(localize("invalid_date_format"))
			isValid = false
		} else if len(startDate.Text) != 10 || startDate.Text[2] != '.' || startDate.Text[5] != '.' {
			startDate.SetText(localize("invalid_date_format"))
			isValid = false
		} else {
			_, err := time.Parse("02.01.2006", startDate.Text)
			if err != nil {
				startDate.SetText(localize("invalid_date"))
				isValid = false
			}
		}

		// Проверка конечной даты
		if endDate.Text == "" {
			endDate.SetText(localize("invalid_date_format"))
			isValid = false
		} else if len(endDate.Text) != 10 || endDate.Text[2] != '.' || endDate.Text[5] != '.' {
			endDate.SetText(localize("invalid_date_format"))
			isValid = false
		} else {
			_, err := time.Parse("02.01.2006", endDate.Text)
			if err != nil {
				endDate.SetText(localize("invalid_date"))
				isValid = false
			}
		}

		// Если есть ошибки в датах, прерываем отправку
		if !isValid {
			fmt.Println("Cannot submit with invalid dates")
			return
		}

		// Проверяем, что тип отпуска выбран
		if leaveTypeLabel.Text == localize("select_leave_type") {
			fmt.Println("Leave type not selected")
			return
		}

		// Парсим даты (они уже проверены выше)
		start, _ := time.Parse("02.01.2006", startDate.Text)
		end, _ := time.Parse("02.01.2006", endDate.Text)

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
