package port

type IAuthenticationService interface {
	GenerateToken(userIdentifier string) (string, string, int64, error)
}
