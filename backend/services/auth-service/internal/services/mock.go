package services

type FakeAuthService struct {
	// Expected
	ResponseToReturn *LoginResponse
	ErrorToReturn    error
	// check what we called
	LoginCalledWith *LoginRequest
}

func (f *FakeAuthService) Login(req *LoginRequest) (*LoginResponse, error) {
	f.LoginCalledWith = req
	return f.ResponseToReturn, f.ErrorToReturn
}
