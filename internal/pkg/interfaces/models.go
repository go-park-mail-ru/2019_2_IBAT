package interfaces

import (
	"github.com/google/uuid"
)

type SeekerReg struct {
	Email      string `json:"email"`
	FirstName  string `json:"first_name"`
	SecondName string `json:"second_name"`
	Password   string `json:"password"`
}

type EmployerReg struct {
	CompanyName string `json:"company_name"`
	Site        string `json:"site"`
	FirstName   string `json:"first_name"`
	SecondName  string `json:"second_name"`
	Email       string `json:"email"`
	Number      string `json:"number"`
	ExtraNumber string `json:"extra_number"`
	Password    string `json:"password"`
	City        string `json:"city"`
	EmplNum     int    `json:"empl_num"`
}

type Seeker struct {
	Email      string      `json:"email"`
	FirstName  string      `json:"first_name"`
	SecondName string      `json:"second_name"`
	Password   string      `json:"password"`
	Resumes    []uuid.UUID `json:"resumes"`
} //add extra fields

type Employer struct {
	CompanyName string      `json:"company_name"`
	Site        string      `json:"site"`
	FirstName   string      `json:"first_name"`
	SecondName  string      `json:"second_name"`
	Email       string      `json:"email"`
	Number      string      `json:"number"`
	ExtraNumber string      `json:"extra_number"`
	Password    string      `json:"password"`
	City        string      `json:"city"`
	EmplNum     int         `json:"empl_num"`
	Vacancies   []uuid.UUID `json:"-"`
} //add extra fields

type Resume struct {
	OwnerID     uuid.UUID `json:"-"`
	FirstName   string    `json:"first_name"`
	SecondName  string    `json:"second_name"`
	City        string    `json:"city"`
	Number      string    `json:"number"`
	BirthDate   string    `json:"birth_date"`
	Sex         string    `json:"sex"`
	Citizenship string    `json:"citizenship"`
	Experience  string    `json:"experience"`
	Profession  string    `json:"profession"`
	Position    string    `json:"position"`
	Wage        string    `json:"wage"`
	Education   string    `json:"education"`
	About       string    `json:"About"`
}

type Message struct {
	Message string `json:"message"`
}

type Error struct {
	Error string `json:"error"`
}

// type Vacancy struct {
// 	OwnerID uuid.UUID `json:"-"`
// 	// ID uuid.UUID
// 	FirstName   string `json:"first_name"`
// 	SecondName  string `json:"second_name"`
// 	City        string `json:"city"`
// 	Number      string `json:"number"`
// 	BirthDate   string `json:"birth_date"`
// 	Sex         string `json:"sex"`
// 	Citizenship string `json:"citizenship"`
// 	Experience  string `json:"experience"`
// 	Profession  string `json:"profession"`
// 	Position    string `json:"position"`
// 	Wage        string `json:"wage"`
// 	Education   string `json:"education"`
// 	About       string `json:"About"`
// }

type AuthStorageValue struct {
	ID      uuid.UUID
	Expires string
	Class   string
}

type UserAuthInput struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

// type Handler struct {
// 	Storage AuthStorage
// 	Mu      *sync.Mutex
// }
