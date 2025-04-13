package dto

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
	UserInfo   UserInfo
	AuthTokens AuthTokens
}

type UserProfileRequest struct {
	UserID uint `json:"user_id"`
}

type UserProfileResponse struct {
	UserInfo UserInfo
}
