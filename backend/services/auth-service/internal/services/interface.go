package services

type AuthServiceInterface interface {
	Login(req *LoginRequest) (*LoginResponse, error)
}
