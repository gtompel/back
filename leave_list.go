package main

import (
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

// Функция отображения списка отпусков
func showLeaveRequestsImpl(w fyne.Window) {
	var results []struct {
		LastName   string
		FirstName  string
		MiddleName string
		StartDate  time.Time
		EndDate    time.Time
		Type       string
		Reason     string // Добавляем поле для причины
	}

	// Добавляем отчество и причину в запрос
	db.Table("leave_requests").
		Select("employees.last_name, employees.first_name, employees.middle_name, leave_requests.start_date, leave_requests.end_date, leave_requests.type, leave_requests.reason").
		Joins("join employees on leave_requests.employee_id = employees.id").
		Scan(&results)

	// Добавляем колонку для отчества и причины
	data := [][]string{{
		localize("last_name"),
		localize("first_name"),
		localize("middle_name"),
		localize("start_date"),
		localize("end_date"),
		localize("leave_type"),
		localize("reason"), // Добавляем заголовок для причины
	}}

	for _, r := range results {
		data = append(data, []string{
			r.LastName,
			r.FirstName,
			r.MiddleName,
			r.StartDate.Format("02.01.2006"),
			r.EndDate.Format("02.01.2006"),
			r.Type,
			r.Reason, // Добавляем причину
		})
	}

	// Создаем улучшенную таблицу с лучшей читаемостью
	table := widget.NewTable(
		func() (int, int) { return len(data), 7 }, // Увеличиваем количество колонок до 7
		func() fyne.CanvasObject {
			return container.NewHBox(widget.NewLabel(""))
		},
		func(id widget.TableCellID, o fyne.CanvasObject) {
			container := o.(*fyne.Container)
			label := container.Objects[0].(*widget.Label)

			// Устанавливаем текст
			label.SetText(data[id.Row][id.Col])

			// Стилизуем заголовки
			if id.Row == 0 {
				label.TextStyle = fyne.TextStyle{Bold: true}
			} else {
				label.TextStyle = fyne.TextStyle{Bold: false}
			}

			// Выравнивание текста
			label.Alignment = fyne.TextAlignCenter
		},
	)

	// Устанавливаем размеры колонок для лучшей читаемости
	table.SetColumnWidth(0, 120) // Фамилия
	table.SetColumnWidth(1, 120) // Имя
	table.SetColumnWidth(2, 120) // Отчество
	table.SetColumnWidth(3, 100) // Дата начала
	table.SetColumnWidth(4, 100) // Дата окончания
	table.SetColumnWidth(5, 120) // Тип отпуска
	table.SetColumnWidth(6, 200) // Причина - делаем шире для текста

	// Добавляем кнопку "+" для добавления нового отпуска
	addButton := widget.NewButtonWithIcon("", theme.ContentAddIcon(), func() {
		showAddLeaveForm(w)
	})

	// Добавляем заголовок
	title := canvas.NewText(localize("leave_requests_list"), theme.PrimaryColor())
	title.TextSize = 20
	title.TextStyle = fyne.TextStyle{Bold: true}
	title.Alignment = fyne.TextAlignCenter

	// Создаем контейнер с заголовком, таблицей и кнопкой
	content := container.NewBorder(
		container.NewVBox(
			container.NewPadded(title),
			widget.NewSeparator(),
		),
		container.NewHBox(
			layout.NewSpacer(),
			addButton,
		),
		nil,
		nil,
		container.NewScroll(container.NewPadded(table)), // Добавляем скроллинг для таблицы
	)

	w.SetContent(content)
}
