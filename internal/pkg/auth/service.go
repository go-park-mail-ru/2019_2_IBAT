package auth

import (
	. "2019_2_IBAT/internal/pkg/interfaces"

	// auth_serv "2019_2_IBAT/internal/pkg/users/service"

	"net/http"

	"github.com/google/uuid"
)

const CookieName = "session-id"

// Usecase represent the article's usecases
type Service interface {
	CreateSession(id uuid.UUID, class string) (http.Cookie, string, error)
	DeleteSession(cookie *http.Cookie) bool

	GetSession(cookie string) (AuthStorageValue, bool)
	SetRecord(id uuid.UUID, class string) (AuthStorageValue, string, error)

	AuthMiddleware(h http.Handler) http.Handler
}
