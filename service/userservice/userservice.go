package userservice

import (
	"fmt"
	"regexp"

	"github.com/aghaghiamh/gocast/QAGame/entity"
	"github.com/aghaghiamh/gocast/QAGame/pkg/richerr"
	"golang.org/x/crypto/bcrypt"
)

type UserRepo interface {
	IsAlreadyExist(phoneNumber string) (bool, error)
	Register(user entity.User) (entity.User, error)
	GetUserByPhoneNumber(phoneNumber string) (entity.User, error)
	GetUserByID(user_id uint) (entity.User, error)
}

type AuthGenerator interface {
	CreateAccessToken(userID uint) (string, error)
	CreateRefreshToken(userID uint) (string, error)
}

type Service struct {
	auth AuthGenerator
	repo UserRepo
}

func New(repo UserRepo, auth AuthGenerator) Service {
	return Service{
		auth: auth,
		repo: repo,
	}
}

type RegisterRequest struct {
	Name        string `json:"name"`
	PhoneNumber string `json:"phone_number"`
	Password    string `json:"password"`
}

type UserInfo struct {
	UserID      uint   `json:"user_id"`
	Name        string `json:"name"`
	PhoneNumber string `json:"phone_number"`
}

type RegisterResponse struct {
	UserInfo
}

// TODO: All Authorization tasks must handle in the authservice and
// be checked before comming into any other service as a middle-ware.
type AuthTokens struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type LoginRequest struct {
	PhoneNumber string `json:"phone_number"`
	Password    string `json:"password"`
}

type LoginResponse struct {
	UserInfo
	AuthTokens
}

type UserProfileRequest struct {
	UserID uint `json:"user_id"`
}

type UserProfileResponse struct {
	UserInfo
}

func (s *Service) Login(req LoginRequest) (LoginResponse, error) {
	const op = "userservice.Login"
	if valid, err := isValidPhoneNumber(req.PhoneNumber); !valid || err != nil {
		return LoginResponse{}, err
	}

	var loginResponse LoginResponse

	if len(req.Password) == 0 {
		// TODO: OTP Login
	} else if len(req.Password) < 8 {

		return LoginResponse{}, richerr.New(op).
			WithCode(richerr.ErrInvalidInput).
			WithMessage("Password should contains at least 8 characters")
	} else { // Normal Authentication Path
		user, dbErr := s.repo.GetUserByPhoneNumber(req.PhoneNumber)
		if dbErr != nil {
			return LoginResponse{}, richerr.New(op).
				WithError(dbErr).
				WithMessage(fmt.Sprintf(
					"User with %s phone number does not exist. Please register first", req.PhoneNumber))
		}

		if len(user.HashedPassword) == 0 {
			// TODO: OTP Login
		}

		// verify password
		verErr := bcrypt.CompareHashAndPassword([]byte(user.HashedPassword), []byte(req.Password))
		if verErr != nil {
			return LoginResponse{}, richerr.New(op).
				WithError(verErr).
				WithCode(richerr.ErrInvalidInput).
				// TODO: Should I hide the actual errror in this layer?
				WithMessage("Phone number or password is incorrect.")
		}

		accessToken, aErr := s.auth.CreateAccessToken(user.ID)
		if aErr != nil {
			return LoginResponse{}, richerr.New(op).WithError(aErr)
		}

		refreshToken, rErr := s.auth.CreateRefreshToken(user.ID)
		if rErr != nil {
			return LoginResponse{}, richerr.New(op).WithError(rErr)
		}

		loginResponse = LoginResponse{
			UserInfo{
				UserID:      user.ID,
				Name:        user.Name,
				PhoneNumber: user.PhoneNumber,
			},
			AuthTokens{
				AccessToken:  accessToken,
				RefreshToken: refreshToken,
			},
		}
	}
	return loginResponse, nil
}

func (s *Service) Register(req RegisterRequest) (RegisterResponse, error) {
	const op = "userservice.Register"

	if valid, err := isValidPhoneNumber(req.PhoneNumber); !valid || err != nil {

		return RegisterResponse{}, richerr.New(op).WithError(err)
	}

	alreadyExist, err := s.repo.IsAlreadyExist(req.PhoneNumber)
	if err != nil {

		return RegisterResponse{}, richerr.New(op).WithError(err)
	}
	if alreadyExist {

		return RegisterResponse{}, richerr.New(op).
			WithCode(richerr.ErrEntityDuplicate).
			WithMessage(fmt.Sprintf("Duplicated user with %s phone number", req.PhoneNumber)).
			WithMetadata(map[string]interface{}{"phone-number": req.PhoneNumber})
	}

	// validate name
	if len(req.Name) < 3 {

		return RegisterResponse{}, richerr.New(op).
			WithCode(richerr.ErrInvalidInput).
			WithMessage("user's name must contains at least 3 character")
	}

	// password validation
	var hashedPassword string
	if len(req.Password) == 0 {
		// TODO: OTP path of Registeration
	} else if len(req.Password) < 8 {

		return RegisterResponse{}, richerr.New(op).
			WithCode(richerr.ErrInvalidInput).
			WithMessage("user's password must contains at least 8 character")
	} else {
		var passErr error
		hashedPassword, passErr = hashPassword(req.Password)
		if passErr != nil {

			return RegisterResponse{}, richerr.New(op).WithError(passErr)
		}
	}

	user := entity.User{
		PhoneNumber:    req.PhoneNumber,
		Name:           req.Name,
		HashedPassword: hashedPassword,
	}
	var regErr error
	user, regErr = s.repo.Register(user)
	if regErr != nil {
		return RegisterResponse{}, richerr.New(op).WithError(regErr)
	}

	return RegisterResponse{
		UserInfo{
			UserID:      user.ID,
			Name:        user.Name,
			PhoneNumber: user.PhoneNumber,
		},
	}, nil
}

// Validate phone number using (+98) 09xxxxxxxxx pattern which x is a digit
func isValidPhoneNumber(phoneNumber string) (bool, error) {
	const op = "userservice.isValidPhoneNumber"

	r, err := regexp.Compile(`^(\(?\+98\)?)?[-\s]?(09)(\d{9})$`)
	if err != nil {

		return false, richerr.New(op).
			WithCode(richerr.ErrInvalidInput).
			WithMessage("Phone number does not satisfy the valid pattern of `(+98) 09xxxxxxxxx`.")
	}

	if !r.MatchString(phoneNumber) {

		return false, richerr.New(op).
			WithCode(richerr.ErrInvalidInput).
			WithMessage(`phone number must start with "09", and all characters must be digits`)
	}

	return true, nil
}

func hashPassword(password string) (string, error) {
	const op = "userservice.hashPassword"

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {

		return "", richerr.New(op).
			WithCode(richerr.ErrServer).
			WithMessage(fmt.Sprintf("Couldn't hash the password: %s", password))
	}

	return string(hashedPassword), nil
}

func (s *Service) GetProfile(req UserProfileRequest) (UserProfileResponse, error) {
	const op = "userservice.GetUserProfile"
	user, gErr := s.repo.GetUserByID(req.UserID)
	if gErr != nil {

		return UserProfileResponse{}, richerr.New(op).
			WithError(gErr).
			WithMessage(fmt.Sprintf("User with %d id does not exist. Please register first.", req.UserID))
	}

	return UserProfileResponse{
		UserInfo{
			UserID:      user.ID,
			Name:        user.Name,
			PhoneNumber: user.PhoneNumber,
		},
	}, nil
}
