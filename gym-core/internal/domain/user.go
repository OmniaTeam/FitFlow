package domain

type User struct {
	ID int `json:"id"`
	//FIO   string `json:"fio"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	//Password *string `json:"-"`
	GoogleID *string `json:"-"`
	VkID     *string `json:"-"`
	//IsSSO    *string  `json:"is_sso"`
	Roles []string `json:"roles"`
}
