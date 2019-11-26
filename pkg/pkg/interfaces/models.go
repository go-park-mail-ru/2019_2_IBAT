package interfaces

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
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
const AwaitSt = "awaits"
const RejectedSt = "rejected"
const Accepted = "accepted"

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
	Position         string    `json:"position"            db:"position"`
	Wage             string    `json:"wage"                db:"wage"`
	Education        string    `json:"education"           db:"education"`
	About            string    `json:"about"               db:"about"`
	Spheres          []Pair    `json:"spheres"`
}

type Pair struct {
	First  string `json:"first"`
	Second string `json:"second"`
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
	Position         string    `json:"position"            db:"position"`
	Tasks            string    `json:"tasks"               db:"tasks"`
	Requirements     string    `json:"requirements"        db:"requirements"`
	WageFrom         string    `json:"wage_from"           db:"wage_from"`
	TypeOfEmployment string    `json:"type_of_employment"  db:"type_of_employment"`
	WorkSchedule     string    `json:"work_schedule"       db:"work_schedule"`
	WageTo           string    `json:"wage_to"             db:"wage_to"`
	Conditions       string    `json:"conditions"          db:"conditions"`
	Favorite         bool      `json:"favorite"            db:"favorite"`
	IsResponded      bool      `json:"is_responded"        db:"is_responded"`
	About            string    `json:"about"               db:"about"`
	Spheres          []Pair    `json:"spheres"`
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

type TagAnswer struct {
	ParentTag   string   `json:"parent_tag"   db:"parent_tag"`
	ChildTagArr []string `json:"child_tag"    db:"child_tag"`
}

type Tag struct {
	ParentTag string `json:"parent_tag"   db:"parent_tag"`
	ChildTag  string `json:"child_tag"    db:"child_tag"`
}

type NotifStruct struct {
	VacancyId uuid.UUID
	TagIDs    []uuid.UUID
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

type Connect struct {
	Conn *websocket.Conn
	Ch   chan uuid.UUID
}

type ConnectsPerUser struct {
	Channels []chan uuid.UUID
	Mu       *sync.Mutex
	Ch       chan uuid.UUID
}

type WsConnects struct {
	ConsMu   *sync.Mutex
	Connects map[uuid.UUID]*ConnectsPerUser
}

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	maxMessageSize = 512
)

func (c *Connect) ReadPump() {
	// defer func() {
	// 	c.hub.unregister <- c
	// 	c.conn.Close()
	// }()
	c.Conn.SetReadLimit(maxMessageSize)
	c.Conn.SetReadDeadline(time.Now().Add(pongWait))
	c.Conn.SetPongHandler(func(string) error { c.Conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		if _, _, err := c.Conn.NextReader(); err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			// c.Close()
			break
		}
	}
}

// for {
// 	_, message, err := c.conn.ReadMessage()
// 	if err != nil {
// 		if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
// 			log.Printf("error: %v", err)
// 		}
// 		break
// 	}
// 	message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))
// 	c.hub.broadcast <- message
// }

func (c *Connect) WritePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.Conn.Close()
	}()
	for {
		select {
		case id := <-c.Ch: //ok
			// if ok {
			w, err := c.Conn.NextWriter(websocket.TextMessage)
			if err != nil {
				ticker.Stop()
				break
			}
			idJSON, _ := json.Marshal(Id{Id: id.String()})
			w.Write(idJSON)
			w.Close()
			fmt.Printf("id %s was sent user", id.String())
			// } else {
			// 	fmt.Println("Channel closed!")
			// 	//close channel
			// 	return
			// 	// break
			// }
		case <-ticker.C:
			if _, _, err := c.Conn.NextReader(); err != nil {
				c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
				if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
					return
				}
			}
		}
	}
}
