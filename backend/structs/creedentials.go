package structs

type Credentials struct {
	Password string `json:"password", db:"password"`
	Username string `json:"username", db:"username"`
	Role     string `json:"role", db:"role"`
	Email    string `json:"email", db:"email"`
}
