package matchingvalidator

import (
	"fmt"

	"github.com/aghaghiamh/gocast/QAGame/dto"
	"github.com/aghaghiamh/gocast/QAGame/entity"
	"github.com/aghaghiamh/gocast/QAGame/pkg/errmsg"
	"github.com/aghaghiamh/gocast/QAGame/pkg/richerr"

	"github.com/go-ozzo/ozzo-validation/v4"
)

func (v MatchingValidator) ValidateAddToWaitingList(req dto.AddToWaitingListRequest) (map[string]string, error) {
	const op = "validator.ValidateLoginRequest"

	err := validation.ValidateStruct(&req,

		validation.Field(&req.UserID, validation.Required),

		validation.Field(&req.Category, validation.Required, validation.By(v.isCategoryValid)),
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

func (v MatchingValidator) isCategoryValid(value interface{}) error {
	const op = "MatchingValidator.IsCategoryValid"
	category := value.(entity.Category)

	if !category.IsValid() {
		return fmt.Errorf(errmsg.ErrMsgCategoryIsNotValid)
	}

	return nil
}
