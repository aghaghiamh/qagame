package uservalidator

import (
	"regexp"

	"github.com/aghaghiamh/gocast/QAGame/dto"
	"github.com/aghaghiamh/gocast/QAGame/pkg/richerr"

	"github.com/go-ozzo/ozzo-validation/v4"
)

func (v UserValidator) ValidateLoginRequest(req dto.LoginRequest) (map[string]string, error) {
	const op = "validator.ValidateLoginRequest"

	err := validation.ValidateStruct(&req,

		validation.Field(
			&req.PhoneNumber,
			validation.Required,
			validation.Match(regexp.MustCompile(PhoneNumberRegex)).Error("Phone number does not satisfy the valid pattern of `(+98) 09xxxxxxxxx`."),
			validation.By(v.doesPhoneNumberExist)),

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

func (v UserValidator) doesPhoneNumberExist(value interface{}) error {
	phoneNumber := value.(string)

	_, err := v.repo.GetUserByPhoneNumber(phoneNumber)
	if err != nil {
		return err
	}

	return nil
}
