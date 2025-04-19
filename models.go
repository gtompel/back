package main

import (
	"time"
)

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
