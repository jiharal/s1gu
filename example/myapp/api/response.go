package api

import (
	"context"
	"encoding/json"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

type (
	BaseResponse struct {
		Errors []string `json:"errors,omitempty"`
	}
	Response struct {
		Status       string `json:"status"`
		BaseResponse `json:"errors"`
		Data         interface{} `json:"result"`
	}
)

var (
	MessageGeneralError  = "Ada kesalahan, Silahakan coba beberapa saat lagi."
	MessageUnauthorized  = "Silahkan login terlebih dahulu atau login ulang."
	MessageInvalidLogin  = "Email atau password anda salah."
	MessageAccountExists = "Akun anda sudah terdaftar, silahkan login."
)

// RespondError writes / respond with JSON-formatted request of given message & http status.
func RespondError(w http.ResponseWriter, message string, status int) {
	resp := Response{
		Status: http.StatusText(status),
		Data:   nil,
		BaseResponse: BaseResponse{
			Errors: []string{
				message,
			},
		},
	}

	w.WriteHeader(status)
	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(&resp); err != nil {
		logger.Err.WithError(err).Println("Encode response error.")
		return
	}
}

func ComparePassword(ctx context.Context, hash string, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}
