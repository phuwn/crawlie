package handler

type Config struct {
	ListeningPort int    `json:"listening_port"`
	JwtSecretKey  string `json:"jwt_secret_key"`
}
