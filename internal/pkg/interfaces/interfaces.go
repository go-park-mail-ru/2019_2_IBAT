package interfaces

import (
	"github.com/google/uuid"
)

const SeekerStr = "Seeker"
const EmployerStr = "Employer"

type UserStorage interface {
	CreateEmployer(seekerInput EmployerReg) (uuid.UUID, bool)
	CreateSeeker(seekerInput SeekerReg) (uuid.UUID, bool)
	CreateResume(resumeReg Resume, userId uuid.UUID) (uuid.UUID, bool)
	CreateVacancy(vacancyReg Vacancy, userId uuid.UUID) (uuid.UUID, bool)

	DeleteEmployer(id uuid.UUID)
	DeleteSeeker(id uuid.UUID)
	DeleteResume(id uuid.UUID) bool
	DeleteVacancy(id uuid.UUID) bool

	CheckUser(email string, password string) (uuid.UUID, string, bool)

	PutSeeker(seekerInput SeekerReg, id uuid.UUID) bool
	PutEmployer(employerInput EmployerReg, id uuid.UUID) bool
	PutResume(resume Resume, userId uuid.UUID, resumeId uuid.UUID) bool
	PutVacancy(vacavcy Vacancy, userId uuid.UUID, resumeId uuid.UUID) bool

	GetEmployers() map[uuid.UUID]Employer
	GetResumes() map[uuid.UUID]Resume
	GetVacancies() map[uuid.UUID]Vacancy
	GetSeekers() map[uuid.UUID]Seeker

	GetSeeker(id uuid.UUID) (Seeker, bool)
	GetEmployer(id uuid.UUID) (Employer, bool)
	GetResume(id uuid.UUID) (Resume, bool)
	GetVacancy(id uuid.UUID) (Vacancy, bool)
}

type AuthStorage interface {
	Get(cookie string) (AuthStorageValue, bool)
	Set(id uuid.UUID, class string) (AuthStorageValue, string)
	Delete(cookie string) bool
}
