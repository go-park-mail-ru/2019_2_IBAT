package users

import (
	. "2019_2_IBAT/pkg/pkg/models"

	"github.com/google/uuid"
)

type Repository interface {
	CreateEmployer(seekerInput Employer) bool
	CreateSeeker(seekerInput Seeker) bool
	CreateResume(resumeReg Resume) bool
	CreateVacancy(vacancyReg Vacancy) bool

	CreateRespond(respond Respond, userId uuid.UUID) bool
	CreateFavorite(favVac FavoriteVacancy) bool

	DeleteUser(id uuid.UUID) error
	DeleteResume(id uuid.UUID) error
	DeleteVacancy(id uuid.UUID) error
	DeleteFavoriteVacancy(id uuid.UUID, authInfo AuthStorageValue) error

	CheckUser(email string, password string) (uuid.UUID, string, bool)

	PutSeeker(seekerInput SeekerReg, id uuid.UUID) bool
	PutEmployer(employerInput EmployerReg, id uuid.UUID) bool
	PutResume(resume Resume, userId uuid.UUID, resumeId uuid.UUID) bool
	PutVacancy(vacavcy Vacancy, userId uuid.UUID, resumeId uuid.UUID) bool

	GetEmployers(params map[string]interface{}) ([]Employer, error)
	GetSeekers() ([]Seeker, error)
	GetResumes(authInfo AuthStorageValue, params map[string]interface{}) ([]Resume, error)
	GetVacancies(authInfo AuthStorageValue, params map[string]interface{}) ([]Vacancy, error)

	GetResponds(record AuthStorageValue, params map[string]string) ([]Respond, error)
	GetFavoriteVacancies(record AuthStorageValue) ([]Vacancy, error)

	GetSeeker(id uuid.UUID) (Seeker, error)
	GetEmployer(id uuid.UUID) (Employer, error)
	GetResume(id uuid.UUID) (Resume, error)
	GetVacancy(id uuid.UUID, userId uuid.UUID) (Vacancy, error)

	SetImage(id uuid.UUID, class string, imageName string) bool

	GetTags() ([]Tag, error)
	GetTagIDs(tags []Pair) ([]uuid.UUID, error)
	GetVacancyTagIDs(vacancyId uuid.UUID) ([]uuid.UUID, error)
}
