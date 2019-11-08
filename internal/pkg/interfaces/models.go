package interfaces

import (
	"context"
	"time"

	"github.com/google/uuid"
)

const UnauthorizedMsg = "Unauthorized"
const InternalErrorMsg = "Internal server error"
const ForbiddenMsg = "Forbidden"
const InvalidIdMsg = "Invalid ID"
const InvalidJSONMsg = "Invalid JSON"
const BadRequestMsg = "Bad request"

const EmailExistsMsg = "Email already exists"
const InvPassOrEmailMsg = "Invalid password or email"

//respond/offer status
const AwaitSt = "Await"
const RejectedSt = "RejectedSt"
const Accepted = "Accepted"

// type key string

// const AuthRec key = "AuthRecord" ///fix

var Loc *time.Location

func init() {
	Loc, _ = time.LoadLocation("Europe/Moscow")
}

const TimeFormat = time.RFC3339 //duplicate

type SeekerReg struct {
	Email      string `json:"email"`
	FirstName  string `json:"first_name"`
	SecondName string `json:"second_name"`
	Password   string `json:"password"`
}

type EmployerReg struct {
	CompanyName      string `json:"company_name"`
	Site             string `json:"site"`
	FirstName        string `json:"first_name"`
	SecondName       string `json:"second_name"`
	Email            string `json:"email"`
	PhoneNumber      string `json:"phone_number"`
	ExtraPhoneNumber string `json:"extra_phone_number"`
	Password         string `json:"password"`
	Region           string `json:"region"`
	EmplNum          string `json:"empl_num"`
	//SpheresOfWork    string      `json:"spheres_of_work"    db:"spheres_of_work"`
}

type Seeker struct {
	ID         uuid.UUID   `json:"id"                 db:"id"`
	FirstName  string      `json:"first_name"         db:"first_name"`
	SecondName string      `json:"second_name"        db:"second_name"`
	Email      string      `json:"email"              db:"email"`
	Password   string      `json:"password"           db:"-"`
	PathToImg  string      `json:"path_to_img"        db:"path_to_image"`
	Resumes    []uuid.UUID `json:"resumes"            db:"-"`
} //add extra fields

type Employer struct {
	ID               uuid.UUID   `json:"id"                 db:"id"`
	CompanyName      string      `json:"company_name"       db:"company_name"`
	Site             string      `json:"site"               db:"site"`
	FirstName        string      `json:"first_name"         db:"first_name"`
	SecondName       string      `json:"second_name"        db:"second_name"`
	Email            string      `json:"email"              db:"email"`
	PhoneNumber      string      `json:"phone_number"       db:"phone_number"`
	ExtraPhoneNumber string      `json:"extra_phone_number" db:"extra_phone_number"`
	SpheresOfWork    string      `json:"spheres_of_work"    db:"spheres_of_work"`
	Password         string      `json:"password"           db:"-"`
	Region           string      `json:"region"             db:"region"`
	EmplNum          string      `json:"empl_num"           db:"empl_num"`
	PathToImg        string      `json:"path_to_img"        db:"path_to_image"`
	Description      string      `json:"description"        db:"description"`
	Vacancies        []uuid.UUID `json:"vacancies"          db:"-"`
} //add extra fields

type Resume struct {
	ID               uuid.UUID `json:"id"                  db:"id"`
	OwnerID          uuid.UUID `json:"own_id"              db:"own_id"`
	FirstName        string    `json:"first_name"          db:"first_name"`
	SecondName       string    `json:"second_name"         db:"second_name"`
	Region           string    `json:"region"              db:"region"`
	Email            string    `json:"email"               db:"email"`
	PhoneNumber      string    `json:"phone_number"        db:"phone_number"`
	BirthDate        string    `json:"birth_date"          db:"birth_date"`
	Sex              string    `json:"sex"                 db:"sex"`
	TypeOfEmployment string    `json:"type_of_employment"  db:"type_of_employment"`
	WorkSchedule     string    `json:"work_schedule"       db:"work_schedule"`
	Citizenship      string    `json:"citizenship"         db:"citizenship"`
	Experience       string    `json:"experience"          db:"experience"`
	Profession       string    `json:"profession"          db:"profession"`
	Position         string    `json:"position"            db:"position"`
	Wage             string    `json:"wage"                db:"wage"`
	Education        string    `json:"education"           db:"education"`
	About            string    `json:"about"               db:"about"`
}

type Message struct {
	Body string `json:"message"`
}

type Id struct {
	Id string `json:"id"`
}

type Role struct {
	Role string `json:"role"`
}

type Vacancy struct {
	ID               uuid.UUID `json:"id"                  db:"id"`
	OwnerID          uuid.UUID `json:"owner_id"            db:"own_id"`
	Region           string    `json:"region"              db:"region"`
	CompanyName      string    `json:"company_name"        db:"company_name"`
	Experience       string    `json:"experience"          db:"experience"`
	Profession       string    `json:"profession"     db:"profession"`
	Position         string    `json:"position"            db:"position"`
	Tasks            string    `json:"tasks"               db:"tasks"`
	Requirements     string    `json:"requirements"        db:"requirements"`
	WageFrom         string    `json:"wage_from"           db:"wage_from"`
	TypeOfEmployment string    `json:"type_of_employment"  db:"type_of_employment"`
	WorkSchedule     string    `json:"work_schedule"       db:"work_schedule"`
	WageTo           string    `json:"wage_to"             db:"wage_to"`
	Conditions       string    `json:"conditions"          db:"conditions"`
	About            string    `json:"about"               db:"about"`
}

type Respond struct {
	Status    string
	ResumeID  uuid.UUID `json:"resume_id"        db:"resume_id"`
	VacancyID uuid.UUID `json:"vacancy_id"       db:"vacancy_id"`
}

// type FavoriteResume struct {
// 	PersonID  uuid.UUID `json:"resume_id"        db:"resume_id"`
// 	VacancyID uuid.UUID `json:"vacancy_id"       db:"vacancy_id"`
// }

type FavoriteVacancy struct {
	PersonID  uuid.UUID `json:"person_id"        db:"person_id"`
	VacancyID uuid.UUID `json:"vacancy_id"       db:"vacancy_id"`
}

type AuthStorageValue struct {
	ID      uuid.UUID
	Expires string
	Role    string
}

type UserAuthInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Error struct {
	Message string            `json:"error"`
	Params  map[string]string `json:"params"`
}

type key string

const AuthRec key = "AuthRecord" ///fix

func NewContext(ctx context.Context, authInfo AuthStorageValue) context.Context {
	return context.WithValue(ctx, AuthRec, authInfo)
}

func FromContext(ctx context.Context) (AuthStorageValue, bool) {
	authInfo, ok := ctx.Value(AuthRec).(AuthStorageValue)
	return authInfo, ok
}
