package users

import (
	. "hh_workspace/2019_2_IBAT/internal/pkg/interfaces"
	"sync"

	"github.com/google/uuid"
)

type MapUserStorage struct {
	SekMu  *sync.Mutex
	EmplMu *sync.Mutex
	ResMu  *sync.Mutex

	SeekerStorage   map[uuid.UUID]Seeker
	EmployerStorage map[uuid.UUID]Employer
	ResumeStorage   map[uuid.UUID]Resume
}

func (m *MapUserStorage) CreateSeeker(seekerInput SeekerReg) (uuid.UUID, bool) {
	var err_flag bool
	m.SekMu.Lock()
	for _, user := range m.SeekerStorage {
		if user.Email == seekerInput.Email {
			err_flag = true
			break
		}
	}
	m.SekMu.Unlock()
	if err_flag {
		return uuid.UUID{}, false
	}

	m.EmplMu.Lock()
	for _, user := range m.EmployerStorage {
		if user.Email == seekerInput.Email {
			err_flag = true
			break
		}
	}
	m.EmplMu.Unlock() //hard duplication
	if err_flag {
		return uuid.UUID{}, false
	}

	id := uuid.New()
	newSeeker := Seeker{
		Email:      seekerInput.Email,
		FirstName:  seekerInput.FirstName,
		SecondName: seekerInput.SecondName,
		Password:   seekerInput.Password,
		Resumes:    make([]uuid.UUID, 0),
	}

	m.SekMu.Lock()
	m.SeekerStorage[id] = newSeeker
	// for i, item := range m.SeekerStorage { //to remove
	// 	fmt.Printf("uuid: %s  value: %s\n", i, item)
	// }
	// fmt.Println()
	m.SekMu.Unlock()

	return id, true
}

func (m MapUserStorage) CreateEmployer(employerInput EmployerReg) (uuid.UUID, bool) {
	var err_flag bool
	m.EmplMu.Lock()
	for _, user := range m.EmployerStorage {
		if user.Email == employerInput.Email || user.CompanyName == employerInput.CompanyName {
			err_flag = true
			break
		}
	}
	m.EmplMu.Unlock()

	if err_flag {
		return uuid.UUID{}, false
	}

	m.SekMu.Lock()
	for _, user := range m.SeekerStorage {
		if user.Email == employerInput.Email {
			err_flag = true
			break
		}
	}
	m.SekMu.Unlock() //hard duplication

	if err_flag {
		return uuid.UUID{}, false
	}

	id := uuid.New()
	newEmployer := Employer{
		CompanyName: employerInput.CompanyName,
		Site:        employerInput.Site,
		FirstName:   employerInput.FirstName,
		SecondName:  employerInput.SecondName,
		Email:       employerInput.Email,
		Number:      employerInput.Number,
		ExtraNumber: employerInput.ExtraNumber,
		Password:    employerInput.Password,
		City:        employerInput.City,
		EmplNum:     employerInput.EmplNum,
		Vacancies:   make([]uuid.UUID, 0),
	}

	m.EmplMu.Lock()
	m.EmployerStorage[id] = newEmployer
	m.EmplMu.Unlock()

	return id, true
}

func (m MapUserStorage) DeleteEmployer(id uuid.UUID) {
	m.EmplMu.Lock()
	delete(m.EmployerStorage, id)

	// fmt.Print("Deleted\nEmployers:\n")
	// for i, item := range m.EmployerStorage { //to remove
	// 	fmt.Printf("uuid: %s  value: %s\n", i, item)
	// }
	// fmt.Println()

	m.EmplMu.Unlock()

}

func (m MapUserStorage) DeleteSeeker(id uuid.UUID) {
	m.SekMu.Lock()
	delete(m.SeekerStorage, id)

	// fmt.Print("Deleted\nUsers:\n")
	// for i, item := range m.SeekerStorage { //to remove
	// 	fmt.Printf("uuid: %s  value: %s\n", i, item)
	// }
	// fmt.Println()

	m.SekMu.Unlock()
}

func (m MapUserStorage) GetSeekers() []Seeker {
	res := make([]Seeker, 0)

	m.SekMu.Lock()
	for _, seeker := range m.SeekerStorage {
		res = append(res, seeker)
	}
	m.SekMu.Unlock()

	return res
}

func (m MapUserStorage) CheckUser(email string, password string) (uuid.UUID, string, bool) {
	resId := uuid.UUID{}
	var class string
	var flag bool

	//can be parallel
	m.SekMu.Lock()
	for i, seeker := range m.SeekerStorage {
		if seeker.Email == email && seeker.Password == password {
			resId = i
			class = SeekerStr
			flag = true
			break
		}
	}
	m.SekMu.Unlock()
	if flag {
		return resId, class, flag
	}

	m.EmplMu.Lock()
	for i, employer := range m.EmployerStorage {
		if employer.Email == email && employer.Password == password {
			resId = i
			class = EmployerStr
			flag = true
			break
		}
	}
	m.EmplMu.Unlock()

	return resId, class, flag
}

func (m MapUserStorage) CreateResume(resumeReg Resume, userId uuid.UUID) (uuid.UUID, bool) {

	id := uuid.New()
	resumeReg.OwnerID = userId

	m.ResMu.Lock()
	m.ResumeStorage[id] = resumeReg
	m.ResMu.Unlock()

	//what if user were deleted on this line?
	//should use two locks in the same time?

	flag := true
	m.SekMu.Lock()
	newSeeker, ok := m.SeekerStorage[userId]

	if !ok {
		flag = false
		m.ResMu.Lock()
		delete(m.ResumeStorage, id)
		m.ResMu.Unlock()
	} else {
		newSeeker.Resumes = append(newSeeker.Resumes, id)
		m.SeekerStorage[userId] = newSeeker
	}

	m.SekMu.Unlock()

	if !flag {
		return id, false
	}

	// m.ResMu.Lock()
	// for i, item := range m.ResumeStorage { //to remove
	// 	fmt.Printf("uuid: %s\nvalue: %s\n", i, item)
	// }
	// fmt.Println()
	// m.ResMu.Unlock()

	return id, true
}

func (m MapUserStorage) GetResume(id uuid.UUID) (Resume, bool) {
	m.ResMu.Lock()
	result, ok := m.ResumeStorage[id]
	m.ResMu.Unlock()

	if !ok {
		return Resume{}, false
	}
	return result, true
}

func (m MapUserStorage) DeleteResume(id uuid.UUID) bool {
	m.ResMu.Lock()
	// fmt.Println("Resume deleted")

	delete(m.ResumeStorage, id)

	// for i, item := range m.ResumeStorage { //to remove
	// 	fmt.Printf("uuid: %s\nvalue: %s\n", i, item)
	// }
	// fmt.Println()

	m.ResMu.Unlock()
	return true //make false case
}

func (m MapUserStorage) GetSeeker(id uuid.UUID) Seeker {
	m.SekMu.Lock()
	res := m.SeekerStorage[id]
	m.SekMu.Unlock()

	return res
}

func (m MapUserStorage) GetEmployer(id uuid.UUID) Employer {
	m.EmplMu.Lock()
	res := m.EmployerStorage[id]
	m.EmplMu.Unlock()

	return res
}

func (m MapUserStorage) PutSeeker(seekerInput SeekerReg, id uuid.UUID) bool {
	var err_flag bool

	m.SekMu.Lock()

	for i, user := range m.SeekerStorage {
		if user.Email == seekerInput.Email {
			if i == id {
				continue
			}
			err_flag = true
			break
		}
	}
	m.SekMu.Unlock()
	if err_flag {
		return false
	}

	m.EmplMu.Lock()
	for _, user := range m.EmployerStorage {
		if user.Email == seekerInput.Email {
			err_flag = true
			break
		}
	}
	m.EmplMu.Unlock() //hard duplication

	if err_flag {
		return false
	}

	m.SekMu.Lock()
	resumes := m.SeekerStorage[id].Resumes

	m.SeekerStorage[id] = Seeker{
		Email:      seekerInput.Email,
		FirstName:  seekerInput.FirstName,
		SecondName: seekerInput.SecondName,
		Password:   seekerInput.Password,
		Resumes:    resumes,
	}

	m.SekMu.Unlock()

	return true
}

func (m MapUserStorage) PutEmployer(employerInput EmployerReg, id uuid.UUID) bool {
	var err_flag bool

	m.EmplMu.Lock()
	for i, user := range m.EmployerStorage {
		if user.Email == employerInput.Email || user.CompanyName == employerInput.CompanyName {
			if i == id {
				continue
			}
			err_flag = true
			break
		}
	}
	m.EmplMu.Unlock()

	if err_flag {
		return false
	}

	m.SekMu.Lock()
	for _, user := range m.SeekerStorage {
		if user.Email == employerInput.Email {
			err_flag = true
			break
		}
	}
	m.SekMu.Unlock() //hard duplication

	if err_flag {
		return false
	}

	m.EmplMu.Lock()
	vacancies := m.EmployerStorage[id].Vacancies

	m.EmployerStorage[id] = Employer{
		CompanyName: employerInput.CompanyName,
		Site:        employerInput.Site,
		FirstName:   employerInput.FirstName,
		SecondName:  employerInput.SecondName,
		Email:       employerInput.Email,
		Number:      employerInput.Number,
		ExtraNumber: employerInput.ExtraNumber,
		Password:    employerInput.Password,
		City:        employerInput.City,
		EmplNum:     employerInput.EmplNum,
		Vacancies:   vacancies,
	}

	m.EmplMu.Unlock()

	return true
}

func (m MapUserStorage) PutResume(resume Resume, userId uuid.UUID, resumeId uuid.UUID) bool {
	flag := false

	m.SekMu.Lock()
	user := m.SeekerStorage[userId]
	for _, res := range user.Resumes {
		if res == resumeId {
			flag = true
			break
		}
	}
	m.SekMu.Unlock()

	if !flag {
		return false
	}

	resume.OwnerID = userId
	m.ResMu.Lock()
	m.ResumeStorage[resumeId] = resume
	m.ResMu.Unlock()

	return true
}
