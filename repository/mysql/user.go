package mysql

import (
	"database/sql"
	"fmt"

	"github.com/aghaghiamh/gocast/QAGame/entity"
)

type User struct {
	ID          uint
	Name        string
	PhoneNumber string
	CreatedAt   []uint8
}

func (mysql MysqlDB) IsAlreadyExist(phoneNumber string) (bool, error) {
	var fetchedUser User

	query := `SELECT * FROM users WHERE phone_number = ?`
	row := mysql.db.QueryRow(query, phoneNumber)

	sErr := row.Scan(&fetchedUser.ID, &fetchedUser.Name, &fetchedUser.PhoneNumber, &fetchedUser.CreatedAt)
	if sErr != nil {
		if sErr == sql.ErrNoRows {
			return false, nil
		}

		return false, fmt.Errorf("database error: %w", sErr)
	}

	return true, nil
}

func (mysql MysqlDB) Register(user entity.User) (entity.User, error) {
	query := `INSERT INTO users(name, phone_number) VALUES (?, ?) `
	res, err := mysql.db.Exec(query, user.Name, user.PhoneNumber)
	if err != nil {
		return entity.User{}, fmt.Errorf("database error for %s query: %w", query, err)
	}

	lastID, _ := res.LastInsertId()
	user.ID = uint(lastID)

	return user, nil
}
