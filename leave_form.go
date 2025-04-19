package main

import (
	"fmt"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

// Функция отображения формы добавления отпуска в модальном окне
func showAddLeaveFormImpl(parentWindow fyne.Window) {
	// Создаем новое модальное окно
	modalWindow := fyne.CurrentApp().NewWindow(localize("add_leave_request"))
	modalWindow.Resize(fyne.NewSize(600, 700)) // Устанавливаем размер модального окна
	modalWindow.CenterOnScreen()

	// Создаем заголовок формы
	title := canvas.NewText(localize("add_leave_request"), theme.PrimaryColor())
	title.TextSize = 20
	title.TextStyle = fyne.TextStyle{Bold: true}
	title.Alignment = fyne.TextAlignCenter

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
		// Вычисляем позицию для PopUpMenu (под кнопкой)
		popUpPos := fyne.NewPos(buttonPos.X, buttonPos.Y+buttonSize.Height)
		widget.NewPopUpMenu(menu, modalWindow.Canvas()).ShowAtPosition(popUpPos)
	})

	// Создаем поля для ввода дат без немедленной валидации
	startDate := widget.NewEntry()
	startDate.SetPlaceHolder(localize("start_date_format"))

	endDate := widget.NewEntry()
	endDate.SetPlaceHolder(localize("end_date_format"))

	reason := widget.NewMultiLineEntry()
	reason.SetPlaceHolder(localize("reason"))
	reason.SetMinRowsVisible(3) // Увеличиваем высоту поля для причины

	submitButton := widget.NewButtonWithIcon(localize("submit"), theme.ConfirmIcon(), func() {
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
		modalWindow.Close()             // Закрываем модальное окно
		showLeaveRequests(parentWindow) // Обновляем список отпусков в основном окне
	})

	// Кнопка отмены для закрытия модального окна
	cancelButton := widget.NewButtonWithIcon(localize("cancel"), theme.CancelIcon(), func() {
		modalWindow.Close()
	})

	// Создаем форму с отступами для лучшей читаемости
	form := container.NewVBox(
		container.NewPadded(widget.NewLabel(localize("last_name")+":")),
		container.NewPadded(lastName),
		container.NewPadded(widget.NewLabel(localize("first_name")+":")),
		container.NewPadded(firstName),
		container.NewPadded(widget.NewLabel(localize("middle_name")+":")),
		container.NewPadded(middleName),
		container.NewPadded(widget.NewLabel(localize("leave_type")+":")),
		container.NewPadded(leaveTypeButton),
		container.NewPadded(widget.NewLabel(localize("start_date")+":")),
		container.NewPadded(startDate),
		container.NewPadded(widget.NewLabel(localize("end_date")+":")),
		container.NewPadded(endDate),
		container.NewPadded(widget.NewLabel(localize("reason")+":")),
		container.NewPadded(reason),
		container.NewPadded(container.NewHBox(
			layout.NewSpacer(),
			cancelButton,
			submitButton,
		)),
	)

	// Создаем скроллируемый контейнер для формы
	scrollContainer := container.NewScroll(form)

	// Создаем общий контейнер с заголовком и формой
	content := container.NewBorder(
		container.NewVBox(
			container.NewPadded(title),
			widget.NewSeparator(),
		),
		nil,
		nil,
		nil,
		scrollContainer,
	)

	modalWindow.SetContent(content)
	modalWindow.Show() // Показываем модальное окно
}
