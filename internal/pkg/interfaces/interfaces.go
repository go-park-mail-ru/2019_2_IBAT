package interfaces

const SeekerStr = "seeker"
const EmployerStr = "employer"

// type UserStorage interface {
// 	CreateEmployer(seekerInput EmployerReg) (uuid.UUID, bool)
// 	CreateSeeker(seekerInput SeekerReg) (uuid.UUID, bool)
// 	CreateResume(resumeReg Resume, userId uuid.UUID) (uuid.UUID, bool)
// 	CreateVacancy(vacancyReg Vacancy, userId uuid.UUID) (uuid.UUID, bool)

// 	CreateRespond(respond Respond, userId uuid.UUID) (uuid.UUID, bool)

// 	DeleteUser(id uuid.UUID) error
// 	DeleteResume(id uuid.UUID) error
// 	DeleteVacancy(id uuid.UUID) error

// 	CheckUser(email string, password string) (uuid.UUID, string, bool)

// 	PutSeeker(seekerInput SeekerReg, id uuid.UUID) bool
// 	PutEmployer(employerInput EmployerReg, id uuid.UUID) bool
// 	PutResume(resume Resume, userId uuid.UUID, resumeId uuid.UUID) bool
// 	PutVacancy(vacavcy Vacancy, userId uuid.UUID, resumeId uuid.UUID) bool

// 	GetEmployers() ([]Employer, error)
// 	GetSeekers() ([]Seeker, error)
// 	GetResumes() ([]Resume, error)
// 	GetVacancies() ([]Vacancy, error)

// 	GetResponds(record AuthStorageValue, params map[string]string) ([]Respond, error)

// 	GetSeeker(id uuid.UUID) (Seeker, error)
// 	GetEmployer(id uuid.UUID) (Employer, error)
// 	GetResume(id uuid.UUID) (Resume, error)
// 	GetVacancy(id uuid.UUID) (Vacancy, error)

// 	SetImage(id uuid.UUID, class string, imageName string) bool
// }

// type AuthStorage interface {
// 	Get(cookie string) (AuthStorageValue, bool)
// 	Set(id uuid.UUID, class string) (AuthStorageValue, string, error) //bool
// 	Delete(cookie string) bool
// }
