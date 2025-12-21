package dto

// AuthRequest представляет запрос на регистрацию или аутентификацию.
type AuthRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

// IsValid проверяет корректность данных запроса.
func (r *AuthRequest) IsValid() bool {
	return r.Login != "" && r.Password != ""
}
