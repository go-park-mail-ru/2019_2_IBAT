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

	DeleteEmployer(id uuid.UUID)
	DeleteSeeker(id uuid.UUID)

	CheckUser(email string, password string) (uuid.UUID, string, bool)

	PutSeeker(seekerInput SeekerReg, id uuid.UUID) bool
	PutEmployer(employerInput EmployerReg, id uuid.UUID) bool
	// h.Storage.PutResume(resume, user.ID, resumeId)
	PutResume(resume Resume, userId uuid.UUID, resumeId uuid.UUID) bool

	GetSeekers() []Seeker
	GetSeeker(id uuid.UUID) Seeker
	GetEmployer(id uuid.UUID) Employer
	GetResume(id uuid.UUID) (Resume, bool)

	DeleteResume(id uuid.UUID) bool
}

type AuthStorage interface {
	Get(cookie string) (AuthStorageValue, bool)
	Set(id uuid.UUID, class string) string
	Delete(cookie string) string
}
