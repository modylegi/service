package auth

type UserDetails struct {
	Username    string
	Authorities []int
}

const (
	_ = iota
	UserRole
	AdminRole
)

const (
	_ = iota
	AccessToken
	RefreshToken
)
