package auth

type UserInfo struct {
	ProviderID    string
	Email         string
	EmailVerified bool
	Name          string
}