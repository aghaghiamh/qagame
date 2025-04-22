package user

import (
	"database/sql"

	"github.com/aghaghiamh/gocast/QAGame/entity"
	"github.com/aghaghiamh/gocast/QAGame/pkg/errmsg"
	"github.com/aghaghiamh/gocast/QAGame/pkg/richerr"
)

type User struct {
	ID             uint
	Name           string
	PhoneNumber    string
	HashedPassword sql.NullString
	CreatedAt      []uint8
	Role           string
}

func userScanner(row *sql.Row, fetchedUser *User) error {
	return row.Scan(&fetchedUser.ID, &fetchedUser.Name, &fetchedUser.PhoneNumber,
		&fetchedUser.HashedPassword, &fetchedUser.CreatedAt, &fetchedUser.Role)
}

func (mysql Storage) IsAlreadyExist(phoneNumber string) (bool, error) {
	const op = "mysql.IsAlreadyExist"
	var fetchedUser User

	query := `SELECT * FROM users WHERE phone_number = ?`
	row := mysql.db.QueryRow(query, phoneNumber)

	sErr := userScanner(row, &fetchedUser)
	if sErr != nil {
		if sErr == sql.ErrNoRows {
			return false, nil
		}

		return false, richerr.New(op).
			WithError(sErr).
			WithCode(richerr.ErrServer).
			WithMessage(errmsg.ErrMsgCantScanQueryResult)
	}

	return true, nil
}

func (s Storage) Register(user entity.User) (entity.User, error) {
	const op = "mysql.Register"

	var query string
	var res sql.Result
	var err error

	// if password is not provided, it should be saved as null in db
	if len(user.HashedPassword) > 0 {
		query = `INSERT INTO users(name, phone_number, hashed_password) VALUES (?, ?, ?)`
		res, err = s.db.Exec(query, user.Name, user.PhoneNumber, user.HashedPassword)
	} else {
		query = `INSERT INTO users(name, phone_number) VALUES (?, ?)`
		res, err = s.db.Exec(query, user.Name, user.PhoneNumber)
	}

	if err != nil {

		return entity.User{}, richerr.New(op).
			WithError(err).
			WithCode(richerr.ErrServer).
			WithMessage(errmsg.ErrMsgCantScanQueryResult)
	}

	lastID, _ := res.LastInsertId()
	user.ID = uint(lastID)

	return user, nil
}

func (s Storage) GetUserByPhoneNumber(phoneNumber string) (entity.User, error) {
	const op = "mysql.GetUserByPhoneNumber"
	var fetchedUser User

	query := `SELECT * FROM users WHERE phone_number = ?`
	row := s.db.QueryRow(query, phoneNumber)

	sErr := userScanner(row, &fetchedUser)
	if sErr != nil {
		richErr := richerr.New(op).WithError(sErr)

		if sErr == sql.ErrNoRows {
			return entity.User{}, richErr.
				WithCode(richerr.ErrEntityNotFound).
				WithMessage(errmsg.ErrMsgNotFound)
		}

		return entity.User{}, richErr.
			WithCode(richerr.ErrServer).
			WithMessage(errmsg.ErrMsgCantScanQueryResult)
	}

	user := entity.User{
		ID:             fetchedUser.ID,
		Name:           fetchedUser.Name,
		PhoneNumber:    fetchedUser.PhoneNumber,
		HashedPassword: fetchedUser.HashedPassword.String,
		Role:           entity.MapToEntityRole(fetchedUser.Role),
	}

	return user, nil
}

func (s Storage) GetUserByID(user_id uint) (entity.User, error) {
	const op = "mysql.GetUserByID"
	var fetchedUser User

	query := `SELECT * FROM users WHERE id = ?`
	row := s.db.QueryRow(query, user_id)

	sErr := userScanner(row, &fetchedUser)
	if sErr != nil {
		richErr := richerr.New(op).WithError(sErr)

		if sErr == sql.ErrNoRows {
			return entity.User{}, richErr.
				WithCode(richerr.ErrEntityNotFound).
				WithMessage(errmsg.ErrMsgNotFound)
		}

		return entity.User{}, richErr.
			WithCode(richerr.ErrServer).
			WithMessage(errmsg.ErrMsgCantScanQueryResult)
	}

	user := entity.User{
		ID:             fetchedUser.ID,
		Name:           fetchedUser.Name,
		PhoneNumber:    fetchedUser.PhoneNumber,
		HashedPassword: fetchedUser.HashedPassword.String,
	}

	return user, nil
}
