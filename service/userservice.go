package userservice

import (
	"fmt"
	"regexp"
	"time"

	"github.com/aghaghiamh/gocast/QAGame/entity"
	utils "github.com/aghaghiamh/gocast/QAGame/pkg"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type UserRepo interface {
	IsAlreadyExist(phoneNumber string) (bool, error)
	Register(user entity.User) (entity.User, error)
	GetUser(phoneNumber string) (entity.User, error)
}

type UserService struct {
	repo UserRepo
}

func New(repo UserRepo) *UserService {
	return &UserService{
		repo: repo,
	}
}

type RegisterRequest struct {
	Name        string `json:"name"`
	PhoneNumber string `json:"phone_number"`
	Password    string `json:"password"`
}

type RegisterResponse struct {
	User entity.User
}

type LoginRequest struct {
	PhoneNumber string `json:"phone_number"`
	Password    string `json:"password"`
}

type LoginResponse struct {
	UserID int
	Token  string
}

func (s *UserService) Login(req LoginRequest) (LoginResponse, error) {
	if valid, err := isValidPhoneNumber(req.PhoneNumber); !valid || err != nil {
		return LoginResponse{}, err
	}

	var loginResponse LoginResponse

	if len(req.Password) == 0 {
		// TODO: OTP Login
	} else if len(req.Password) < 8 {
		return LoginResponse{}, &utils.RichErr{
			Code:    utils.GeneralBusinessLogicError,
			Message: "Password should contains at least 8 characters",
		}
	} else { // Normal Authentication Path
		user, err := s.repo.GetUser(req.PhoneNumber)
		if err != nil {
			return LoginResponse{}, err
		}

		if len(user.HashedPassword) == 0 {
			// TODO: OTP Login
		}

		// verify password
		verErr := bcrypt.CompareHashAndPassword([]byte(user.HashedPassword), []byte(req.Password))
		if verErr != nil {
			return LoginResponse{}, utils.RichErr{
				Code:    utils.GeneralBusinessLogicError,
				Message: "Phone number or password is incorrect.",
			}
		}

		token, tErr := createToken(req.PhoneNumber, "secret-key")
		if tErr != nil {
			return LoginResponse{}, tErr
		}
		loginResponse = LoginResponse{
			UserID: int(user.ID),
			Token: token,
		}
	}
	return loginResponse, nil
}

func (s *UserService) Register(req RegisterRequest) (RegisterResponse, error) {
	if valid, err := isValidPhoneNumber(req.PhoneNumber); !valid || err != nil {
		return RegisterResponse{}, err
	}

	alreadyExist, err := s.repo.IsAlreadyExist(req.PhoneNumber)
	if alreadyExist {
		return RegisterResponse{}, fmt.Errorf("duplicated user with %s phone number", req.PhoneNumber)
	}
	if err != nil {
		return RegisterResponse{}, fmt.Errorf("database error: %w", err)
	}

	// validate name
	if len(req.Name) < 3 {
		return RegisterResponse{}, fmt.Errorf("user's name must contains at least 3 character")
	}

	// password validation
	var hashedPassword string
	if len(req.Password) == 0 {
		// TODO: OTP path of Registeration
	} else if len(req.Password) < 8 {
		return RegisterResponse{}, fmt.Errorf("user's password must contains at least 8 character")
	} else {
		var passErr error
		hashedPassword, passErr = hashPassword(req.Password)
		if passErr != nil {
			return RegisterResponse{}, passErr
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
		return RegisterResponse{}, regErr
	}

	return RegisterResponse{User: user}, nil
}

func createToken(phoneNumber string, secretKey string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, 
        jwt.MapClaims{ 
        "phone_number": phoneNumber, 
        "exp": time.Now().Add(time.Minute * 20).Unix(), 
        })

    tokenString, err := token.SignedString([]byte(secretKey))
    if err != nil {
    return "", utils.RichErr{
		Code: utils.GeneralServerErr,
		Message: fmt.Sprintf("Couldn't sign JWT token: %s", err),
		}
    }

 return tokenString, nil
}

// Validate phone number using (+98) 09xxxxxxxxx pattern which x is a digit
func isValidPhoneNumber(phoneNumber string) (bool, error) {
	r, err := regexp.Compile(`^(\(?\+98\)?)?[-\s]?(09)(\d{9})$`)
	if err != nil {
		return false, &utils.RichErr{
			Code:    utils.GeneralBusinessLogicError,
			Message: "regexp pattern format is incorrect",
		}
	}

	if !r.MatchString(phoneNumber) {
		return false, &utils.RichErr{
			Code:    utils.GeneralBusinessLogicError,
			Message: `phone number must start with "09", and all characters must be digits`,
		}
	}

	return true, nil
}

func hashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", &utils.RichErr{
			Code:    utils.GeneralServerErr,
			Message: fmt.Sprintf("Couldn't hash the password: %s", password),
		}
	}
	return string(hashedPassword), nil
}
