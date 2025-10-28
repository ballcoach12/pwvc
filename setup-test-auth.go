package main

import (
	"crypto/sha256"
	"fmt"
	"log"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Attendee struct {
	ID            uint   `json:"id" gorm:"primaryKey"`
	ProjectID     uint   `json:"project_id" gorm:"not null"`
	Name          string `json:"name" gorm:"not null"`
	Role          string `json:"role"`
	IsFacilitator bool   `json:"is_facilitator" gorm:"default:false"`
	Email         string `json:"email"`
	PinHash       string `json:"-" gorm:"column:pin_hash"` // Don't include in JSON responses
	CreatedAt     int64  `json:"created_at" gorm:"autoCreateTime"`
}

func main() {
	// Open database
	db, err := gorm.Open(sqlite.Open("data/pairwise.db"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Auto migrate to ensure new columns exist
	err = db.AutoMigrate(&Attendee{})
	if err != nil {
		log.Fatal("Failed to migrate:", err)
	}

	// Set up test PINs for existing attendees
	testData := []struct {
		ID    uint
		Email string
		PIN   string
	}{
		{10, "usera@test.com", "1234"},
		{11, "userb@test.com", "5678"},
		{12, "testuser@test.com", "9999"},
	}

	for _, data := range testData {
		// Hash the PIN
		hash := sha256.Sum256([]byte(data.PIN))
		pinHash := fmt.Sprintf("%x", hash)

		// Update attendee
		result := db.Model(&Attendee{}).Where("id = ?", data.ID).Updates(map[string]interface{}{
			"email":    data.Email,
			"pin_hash": pinHash,
		})

		if result.Error != nil {
			log.Printf("Failed to update attendee %d: %v", data.ID, result.Error)
		} else if result.RowsAffected == 0 {
			log.Printf("No attendee found with ID %d", data.ID)
		} else {
			log.Printf("Updated attendee %d with email %s and PIN %s", data.ID, data.Email, data.PIN)
		}
	}

	fmt.Println("Test authentication setup complete!")
	fmt.Println("Test credentials:")
	for _, data := range testData {
		fmt.Printf("  Attendee ID %d: PIN %s\n", data.ID, data.PIN)
	}
}