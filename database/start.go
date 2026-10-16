package database

import (
	"errors"
	"fmt"
	"os"
)

func InitiateDatabase() error {
	file, err := os.OpenFile("./database/database.db", os.O_RDONLY, 0644)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return err
	}

	if file == nil {
		_, err := os.Create("./database/database.db")
		if err != nil {
			return err
		}
		fmt.Println("database.db created")
		return nil

	}
	fmt.Println("database.db already exists")
	return nil
}
