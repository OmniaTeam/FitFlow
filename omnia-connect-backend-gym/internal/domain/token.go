package domain

type UserTokens struct {
	AccessToken  string
	RefreshToken string
}

type CookieData struct {
	Domain string
}
