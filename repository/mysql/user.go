package mysql

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
}

func (mysql MysqlDB) IsAlreadyExist(phoneNumber string) (bool, error) {
	const op = "mysql.IsAlreadyExist"
	var fetchedUser User

	query := `SELECT * FROM users WHERE phone_number = ?`
	row := mysql.db.QueryRow(query, phoneNumber)

	sErr := row.Scan(&fetchedUser.ID, &fetchedUser.Name, &fetchedUser.PhoneNumber, &fetchedUser.HashedPassword, &fetchedUser.CreatedAt)
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

func (mysql MysqlDB) Register(user entity.User) (entity.User, error) {
	const op = "mysql.Register"

	var query string
	var res sql.Result
	var err error

	// if password is not provided, it should be saved as null in db
	if len(user.HashedPassword) > 0 {
		query = `INSERT INTO users(name, phone_number, hashed_password) VALUES (?, ?, ?)`
		res, err = mysql.db.Exec(query, user.Name, user.PhoneNumber, user.HashedPassword)
	} else {
		query = `INSERT INTO users(name, phone_number) VALUES (?, ?)`
		res, err = mysql.db.Exec(query, user.Name, user.PhoneNumber)
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

func (mysql MysqlDB) GetUserByPhoneNumber(phoneNumber string) (entity.User, error) {
	const op = "mysql.GetUserByPhoneNumber"
	var fetchedUser User

	query := `SELECT * FROM users WHERE phone_number = ?`
	row := mysql.db.QueryRow(query, phoneNumber)

	sErr := row.Scan(&fetchedUser.ID, &fetchedUser.Name, &fetchedUser.PhoneNumber, &fetchedUser.HashedPassword, &fetchedUser.CreatedAt)
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

func (mysql MysqlDB) GetUserByID(user_id uint) (entity.User, error) {
	const op = "mysql.GetUserByID"
	var fetchedUser User

	query := `SELECT * FROM users WHERE id = ?`
	row := mysql.db.QueryRow(query, user_id)

	sErr := row.Scan(&fetchedUser.ID, &fetchedUser.Name, &fetchedUser.PhoneNumber, &fetchedUser.HashedPassword, &fetchedUser.CreatedAt)
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
