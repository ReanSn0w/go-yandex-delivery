package api

type Error struct {
	Error   bool   `json:"error"`
	Message string `json:"message"`
}
