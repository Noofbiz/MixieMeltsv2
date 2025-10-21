package database

import (
	"context"
	"log"

	"com.MixieMelts.users/internal/models"
	"golang.org/x/crypto/bcrypt"
)

func (db *DB) Seed(ctx context.Context) {
	// Check if there are any users in the database
	count, err := db.getUsersCount(ctx)
	if err != nil {
		log.Printf("failed to get users count: %v", err)
		return
	}

	if count > 0 {
		log.Println("Database already seeded.")
		return
	}

	log.Println("Seeding database...")

	// Create some sample users
	users := []models.User{
		{
			Username: "admin",
			Email:    "admin@example.com",
			Password: "password",
			IsAdmin:  true,
		},
		{
			Username: "user",
			Email:    "user@example.com",
			Password: "password",
		},
	}

	for _, user := range users {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		if err != nil {
			log.Printf("failed to hash password for user %s: %v", user.Username, err)
			continue
		}
		user.Password = string(hashedPassword)

		_, err = db.CreateUser(ctx, &user)
		if err != nil {
			log.Printf("failed to create user %s: %v", user.Username, err)
		}
	}

	log.Println("Database seeded successfully.")
}

func (db *DB) getUsersCount(ctx context.Context) (int, error) {
	var count int
	query := "SELECT COUNT(*) FROM users"
	err := db.QueryRow(ctx, query).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}
