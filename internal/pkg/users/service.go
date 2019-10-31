package users

import (
	"2019_2_IBAT/internal/pkg/auth"
	. "2019_2_IBAT/internal/pkg/interfaces"
	"io"

	"github.com/google/uuid"
)

type Service interface {
	CreateSeeker(body io.ReadCloser) (uuid.UUID, error)
	CreateEmployer(body io.ReadCloser) (uuid.UUID, error)
	DeleteUser(cookie string, authStor auth.Service) error
	PutSeeker(body io.ReadCloser, id uuid.UUID) error
	PutEmployer(body io.ReadCloser, id uuid.UUID) error

	GetSeeker(id uuid.UUID) (Seeker, error)
	GetEmployer(id uuid.UUID) (Employer, error)

	CreateVacancy(body io.ReadCloser, cookie string, authStor auth.Service) (uuid.UUID, error)
	GetVacancy(vacancyId uuid.UUID) (Vacancy, error)
	DeleteVacancy(vacancyId uuid.UUID, cookie string, authStor auth.Service) error
	PutVacancy(vacancyId uuid.UUID, body io.ReadCloser, cookie string, authStor auth.Service) error

	CreateResume(body io.ReadCloser, cookie string, authStor auth.Service) (uuid.UUID, error)
	DeleteResume(resumeId uuid.UUID, cookie string, authStor auth.Service) error
	GetResume(resumeId uuid.UUID) (Resume, error)
	PutResume(resumeId uuid.UUID, body io.ReadCloser, cookie string, authStor auth.Service) error

	CreateRespond(body io.ReadCloser, cookie string, authStor auth.Service) (uuid.UUID, error)
	GetResponds(cookie string, params map[string]string, authStor auth.Service) ([]Respond, error)

	GetEmployers() ([]Employer, error)
	GetSeekers() ([]Seeker, error)
	GetResumes() ([]Resume, error)
	GetVacancies() ([]Vacancy, error)

	SetImage(id uuid.UUID, class string, imageName string) bool

	CheckUser(email string, password string) (uuid.UUID, string, bool)
}
