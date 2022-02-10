package resources

type AuthHeader struct {
	JWTToken string `header:"Authorization"`
}
