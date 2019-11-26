package users

import (
	. "2019_2_IBAT/pkg/pkg/interfaces"
	"io"

	"github.com/google/uuid"
)

type Service interface {
	CreateSeeker(body io.ReadCloser) (uuid.UUID, error)
	CreateEmployer(body io.ReadCloser) (uuid.UUID, error)
	DeleteUser(authInfo AuthStorageValue) error
	PutSeeker(body io.ReadCloser, id uuid.UUID) error
	PutEmployer(body io.ReadCloser, id uuid.UUID) error

	GetSeeker(id uuid.UUID) (Seeker, error)
	GetEmployer(id uuid.UUID) (Employer, error)

	CreateVacancy(body io.ReadCloser, authInfo AuthStorageValue) (uuid.UUID, error)
	DeleteVacancy(vacancyId uuid.UUID, authInfo AuthStorageValue) error
	GetVacancy(vacancyId uuid.UUID, authInfo AuthStorageValue) (Vacancy, error)
	PutVacancy(vacancyId uuid.UUID, body io.ReadCloser, authInfo AuthStorageValue) error

	CreateResume(body io.ReadCloser, authInfo AuthStorageValue) (uuid.UUID, error)
	DeleteResume(resumeId uuid.UUID, authInfo AuthStorageValue) error
	GetResume(resumeId uuid.UUID) (Resume, error)
	PutResume(resumeId uuid.UUID, body io.ReadCloser, authInfo AuthStorageValue) error

	CreateRespond(body io.ReadCloser, authInfo AuthStorageValue) error
	GetResponds(authInfo AuthStorageValue, params map[string]string) ([]Respond, error)

	CreateFavorite(vacancyId uuid.UUID, authInfo AuthStorageValue) error
	DeleteFavoriteVacancy(vacancyId uuid.UUID, authInfo AuthStorageValue) error

	GetFavoriteVacancies(authInfo AuthStorageValue) ([]Vacancy, error)

	GetEmployers(params map[string]interface{}) ([]Employer, error)
	GetSeekers() ([]Seeker, error)
	GetResumes(authInfo AuthStorageValue, params map[string]interface{}) ([]Resume, error)
	GetVacancies(authInfo AuthStorageValue, params map[string]interface{},
		tagParams map[string]interface{}) ([]Vacancy, error)
	GetTags() (map[string][]string, error)

	SetImage(id uuid.UUID, class string, imageName string) bool

	CheckUser(email string, password string) (uuid.UUID, string, bool)

	Notifications(connects *WsConnects)
}
