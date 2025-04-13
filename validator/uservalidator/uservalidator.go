package uservalidator

import (
	"fmt"
	"regexp"

	"github.com/aghaghiamh/gocast/QAGame/dto"
	"github.com/aghaghiamh/gocast/QAGame/pkg/errmsg"
	"github.com/aghaghiamh/gocast/QAGame/pkg/richerr"

	"github.com/go-ozzo/ozzo-validation/v4"
)

type UserRepo interface {
	IsAlreadyExist(phoneNumber string) (bool, error)
}

type UserValidator struct {
	repo UserRepo
}

func New(repo UserRepo) UserValidator {
	return UserValidator{
		repo: repo,
	}
}

func (v UserValidator) ValidateRegisterRequest(req dto.RegisterRequest) (map[string]string, error) {
	const op = "validator.ValidateRegisterRequest"

	err := validation.ValidateStruct(&req,
		// TODO: better to get these limitation from Config file, in order to not rebuild the image again.
		validation.Field(&req.Name, validation.Required, validation.Length(3, 20)),

		validation.Field(
			&req.PhoneNumber,
			validation.Required,
			validation.Match(regexp.MustCompile(`^(\(?\+98\)?)?[-\s]?(09)(\d{9})$`)).Error("Phone number does not satisfy the valid pattern of `(+98) 09xxxxxxxxx`."),
			validation.By(v.isUniquePhonenumber)),

		validation.Field(&req.Password, validation.Required,
			validation.Match(regexp.MustCompile((`^[A-Za-z0-9!@#%^&*]{8,}$`)))),
	)

	if err != nil {
		fieldErrors := map[string]string{}
		if vErr, ok := err.(validation.Errors); ok {
			for key, val := range vErr {
				if val != nil {
					fieldErrors[key] = val.Error()
				}
			}
		}
		return fieldErrors, richerr.New(op).WithError(err).WithCode(richerr.ErrInvalidInput)
	}

	return nil, nil
}

func (v UserValidator) isUniquePhonenumber(value interface{}) error {
	phoneNumber := value.(string)

	alreadyExist, err := v.repo.IsAlreadyExist(phoneNumber)
	if err != nil {
		return err
	}
	if alreadyExist {
		return fmt.Errorf(errmsg.ErrMsgPhoneNumberAlreadyExist)
	}

	return nil
}
