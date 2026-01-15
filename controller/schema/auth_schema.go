package schema

type RegisterUser struct {
	Username        string `json:"username" form:"username"`
	Password        string `json:"password" form:"password"`
	ConfirmPassword string `json:"confirmPassword" form:"confirmPassword"`
}

type LoginUser struct {
	Username string `json:"username" form:"username"`
	Password string `json:"password" form:"password"`
}
