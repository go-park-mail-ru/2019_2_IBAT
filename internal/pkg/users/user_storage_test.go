package users

import (
	. "2019_2_IBAT/internal/pkg/interfaces"
	"reflect"
	"sync"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestMapUserStorage_CreateSeeker(t *testing.T) {
	type args struct {
		seekerInput SeekerReg
	}

	m := MapUserStorage{
		SekMu:  &sync.Mutex{},
		EmplMu: &sync.Mutex{},
		SeekerStorage: map[uuid.UUID]Seeker{
			uuid.MustParse("6ba7b810-9dad-11d1-80b1-00c04fd430c8"): {
				Email:      "some@mail.com",
				FirstName:  "Vova",
				SecondName: "Zyablikov",
				Password:   "1234",
				Resumes:    make([]uuid.UUID, 0),
			},
		},
		EmployerStorage: map[uuid.UUID]Employer{},
	}

	tests := []struct {
		name   string
		args   args
		wantOk bool
	}{
		{
			name: "Test1",
			args: args{
				SeekerReg{
					Email:      "some_another@mail.com",
					FirstName:  "Pasha",
					SecondName: "Zyablikov",
					Password:   "12345",
				},
			},
			wantOk: true,
		},
		{
			name: "Test2",
			args: args{
				SeekerReg{
					Email:      "some_another@mail.com",
					FirstName:  "Petya",
					SecondName: "Zyablikov",
					Password:   "12345",
				},
			},
			wantOk: false,
		},
		{
			name: "Test3",
			args: args{
				SeekerReg{
					Email:      "third@mail.com",
					FirstName:  "Petr",
					SecondName: "Zyablikov",
					Password:   "12345",
				},
			},
			wantOk: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			id, got1 := m.CreateSeeker(tt.args.seekerInput)

			if got1 != tt.wantOk {
				t.Errorf("MapUserStorage.CreateSeeker() got1 = %v, want %v", got1, tt.wantOk)
			}

			if got1 == tt.wantOk && got1 == false {
				return
			}

			got, class, ok := m.CheckUser(tt.args.seekerInput.Email, tt.args.seekerInput.Password)

			if !reflect.DeepEqual(got, id) {
				require.Equal(t, id, got, "The two values should be the same.")
			}

			if class != SeekerStr {
				t.Error(`class != SeekerStr`)
			}

			if !ok {
				t.Error("Check user failed")
			}
		})
	}
	// for i, item := range m.SeekerStorage {
	// 	fmt.Printf("uuid: %s  value: %s\n", i, item)
	// }
}

func TestMapUserStorage_CreateEmployer(t *testing.T) {
	m := MapUserStorage{
		SekMu:         &sync.Mutex{},
		EmplMu:        &sync.Mutex{},
		SeekerStorage: map[uuid.UUID]Seeker{},
		EmployerStorage: map[uuid.UUID]Employer{
			uuid.MustParse("6ba7b810-9dad-11d1-80b1-00c04fd430c8"): {
				CompanyName:      "Petushki",
				Site:             "petushki.com",
				Email:            "petushki@mail.com",
				FirstName:        "Vova",
				SecondName:       "Zyablikov",
				Password:         "1234",
				PhoneNumber:      "12345678911",
				ExtraPhoneNumber: "12345678910",
				City:             "Petushki",
				EmplNum:          "1488",

				Vacancies: make([]uuid.UUID, 0),
			},
		},
	}
	type args struct {
		employerInput EmployerReg
	}
	tests := []struct {
		name   string
		args   args
		wantOk bool
	}{
		{
			name: "Test1",
			args: args{
				employerInput: EmployerReg{
					CompanyName:      "hh",
					Site:             "hh.ua",
					Email:            "hh@mail.com",
					FirstName:        "HH",
					SecondName:       "Zyablikov",
					Password:         "1234",
					PhoneNumber:      "12345678911",
					ExtraPhoneNumber: "12345678910",
					City:             "hhland",
					EmplNum:          "322",
				},
			},
			wantOk: true,
		},
		{
			name: "Test2",
			args: args{
				employerInput: EmployerReg{
					CompanyName:      "hh",
					Site:             "hh.ua",
					Email:            "hh@mail.com",
					FirstName:        "HH",
					SecondName:       "Zyablikov",
					Password:         "1234",
					PhoneNumber:      "12345678911",
					ExtraPhoneNumber: "12345678910",
					City:             "hhland",
					EmplNum:          "322",
				},
			},
			wantOk: false,
		},
		{
			name: "Test3",
			args: args{
				employerInput: EmployerReg{
					CompanyName:      "BMSTU",
					Site:             "bmstu.ru",
					Email:            "bmstu@mail.com",
					FirstName:        "Tolya",
					SecondName:       "Alex",
					Password:         "1234",
					PhoneNumber:      "12345678911",
					ExtraPhoneNumber: "12345678910",
					City:             "Moscow",
					EmplNum:          "1830",
				},
			},
			wantOk: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			id, gotOk := m.CreateEmployer(tt.args.employerInput)

			if gotOk != tt.wantOk {
				t.Errorf("MapUserStorage.CreateEmployer() got1 = %v, want %v", gotOk, tt.wantOk)
			}

			if gotOk == tt.wantOk && gotOk == false {
				return
			}

			gotId, class, gotOk := m.CheckUser(tt.args.employerInput.Email, tt.args.employerInput.Password)

			if !reflect.DeepEqual(gotId, id) {
				require.Equal(t, id, gotId, "The two values should be the same.")
			}

			if class != EmployerStr {
				t.Error(`class != EmployerStr`)
			}

			if !gotOk {
				t.Error("Check employer failed")
			}
		})
	}
}

func TestMapUserStorage_DeleteEmployer(t *testing.T) {
	m := MapUserStorage{
		SekMu:         &sync.Mutex{},
		EmplMu:        &sync.Mutex{},
		SeekerStorage: map[uuid.UUID]Seeker{},
		EmployerStorage: map[uuid.UUID]Employer{
			uuid.MustParse("6ba7b810-9dad-11d1-80b1-00c04fd430c8"): {
				CompanyName:      "Petushki",
				Site:             "petushki.com",
				Email:            "petushki@mail.com",
				FirstName:        "Vova",
				SecondName:       "Zyablikov",
				Password:         "1234",
				PhoneNumber:      "12345678911",
				ExtraPhoneNumber: "12345678910",
				City:             "Petushki",
				EmplNum:          "1488",
				Vacancies:        make([]uuid.UUID, 0),
			},
			uuid.MustParse("6ba7b811-9dab-11d1-80b1-00c04fd430c8"): {
				CompanyName:      "BMSTU",
				Site:             "bmstu.ru",
				Email:            "bmstu@mail.com",
				FirstName:        "Tolya",
				SecondName:       "Alex",
				Password:         "1234",
				PhoneNumber:      "12345678911",
				ExtraPhoneNumber: "12345678910",
				City:             "Moscow",
				EmplNum:          "1830",
				Vacancies:        make([]uuid.UUID, 0),
			},
		},
	}
	type args struct {
		id uuid.UUID
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "Test1",
			args: args{
				id: uuid.MustParse("6ba7b810-9dad-11d1-80b1-00c04fd430c8"),
			},
		},
		{
			name: "Test2",
			args: args{
				id: uuid.MustParse("6ba7b811-9dab-11d1-80b1-00c04fd430c8"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m.DeleteEmployer(tt.args.id)
			_, ok := m.EmployerStorage[tt.args.id]
			if ok {
				t.Error(`Employer was not deleted`)
			}
		})
	}
}

func TestMapUserStorage_DeleteSeeker(t *testing.T) {
	m := MapUserStorage{
		SekMu:  &sync.Mutex{},
		EmplMu: &sync.Mutex{},
		SeekerStorage: map[uuid.UUID]Seeker{
			uuid.MustParse("6ba7b810-9dad-11d1-80b1-00c04fd430c8"): {
				Email:      "some@mail.com",
				FirstName:  "Vova",
				SecondName: "Zyablikov",
				Password:   "1234",
				Resumes:    make([]uuid.UUID, 0),
			},

			uuid.MustParse("6ba7b810-9bad-11d1-80b1-00c04fd430c8"): {
				Email:      "third@mail.com",
				FirstName:  "Petr",
				SecondName: "Zyablikov",
				Password:   "12345",
				Resumes:    make([]uuid.UUID, 0),
			},

			uuid.MustParse("6ba6b810-9bad-11d1-80b2-00c04fd430c8"): {
				Email:      "some_another@mail.com",
				FirstName:  "Petya",
				SecondName: "Zyablikov",
				Password:   "12345",
				Resumes:    make([]uuid.UUID, 0),
			},
		},
		EmployerStorage: map[uuid.UUID]Employer{},
	}

	type args struct {
		id uuid.UUID
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "Test1",
			args: args{
				id: uuid.MustParse("6ba6b810-9bad-11d1-80b2-00c04fd430c8"),
			},
		},
		{
			name: "Test2",
			args: args{
				id: uuid.MustParse("6ba7b810-9bad-11d1-80b1-00c04fd430c8"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m.DeleteSeeker(tt.args.id)
			_, ok := m.SeekerStorage[tt.args.id]
			if ok {
				t.Error(`User was not deleted`)
			}
		})
	}
}

func TestMapUserStorage_CheckUser(t *testing.T) {
	m := MapUserStorage{
		SekMu:  &sync.Mutex{},
		EmplMu: &sync.Mutex{},
		SeekerStorage: map[uuid.UUID]Seeker{
			uuid.MustParse("6ba7b810-9dad-11d1-80b1-00c04fd430c8"): {
				Email:      "some@mail.com",
				FirstName:  "Vova",
				SecondName: "Zyablikov",
				Password:   "1234",
				Resumes:    make([]uuid.UUID, 0),
			},

			uuid.MustParse("6ba7b810-9bad-11d1-80b1-00c04fd430c8"): {
				Email:      "third@mail.com",
				FirstName:  "Petr",
				SecondName: "Zyablikov",
				Password:   "12345",
				Resumes:    make([]uuid.UUID, 0),
			},

			uuid.MustParse("6ba6b810-9bad-11d1-80b2-00c04fd430c8"): {
				Email:      "some_another@mail.com",
				FirstName:  "Petya",
				SecondName: "Zyablikov",
				Password:   "12345",
				Resumes:    make([]uuid.UUID, 0),
			},
		},
		EmployerStorage: map[uuid.UUID]Employer{
			uuid.MustParse("6ba7b810-9dad-11d1-80b1-10c05fd430c8"): {
				CompanyName:      "Petushki",
				Site:             "petushki.com",
				Email:            "petushki@mail.com",
				FirstName:        "Vova",
				SecondName:       "Zyablikov",
				Password:         "1234",
				PhoneNumber:      "12345678911",
				ExtraPhoneNumber: "12345678910",
				City:             "Petushki",
				EmplNum:          "1488",
				Vacancies:        make([]uuid.UUID, 0),
			},
			uuid.MustParse("6ba7b811-9dab-11d1-80b1-00c04fd441c8"): {
				CompanyName:      "BMSTU",
				Site:             "bmstu.ru",
				Email:            "bmstu@mail.com",
				FirstName:        "Tolya",
				SecondName:       "Alex",
				Password:         "12",
				PhoneNumber:      "12345678911",
				ExtraPhoneNumber: "12345678910",
				City:             "Moscow",
				EmplNum:          "1830",
				Vacancies:        make([]uuid.UUID, 0),
			},
		},
	}

	type args struct {
		Email    string
		Password string
	}
	tests := []struct {
		name      string
		args      args
		wantId    uuid.UUID
		wantClass string
		wantOk    bool
	}{
		{
			name: "Test1",
			args: args{
				Email:    "some_another@mail.com",
				Password: "12345",
			},
			wantId:    uuid.MustParse("6ba6b810-9bad-11d1-80b2-00c04fd430c8"),
			wantClass: SeekerStr,
			wantOk:    true,
		},
		{
			name: "Test2",
			args: args{
				Email:    "third@mail.com",
				Password: "12345",
			},
			wantId:    uuid.MustParse("6ba7b810-9bad-11d1-80b1-00c04fd430c8"),
			wantClass: SeekerStr,
			wantOk:    true,
		},
		{
			name: "Test3",
			args: args{
				Email:    "petushki@mail.com",
				Password: "1234",
			},
			wantId:    uuid.MustParse("6ba7b810-9dad-11d1-80b1-10c05fd430c8"),
			wantClass: EmployerStr,
			wantOk:    true,
		},
		{
			name: "Test4",
			args: args{
				Email:    "bmstu@mail.com",
				Password: "12",
			},
			wantId:    uuid.MustParse("6ba7b811-9dab-11d1-80b1-00c04fd441c8"),
			wantClass: EmployerStr,
			wantOk:    true,
		},
		{
			name: "Test5",
			args: args{
				Email:    "bmstu@mail.com",
				Password: "12345",
			},
			wantId:    uuid.UUID{},
			wantClass: "",
			wantOk:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotId, gotClass, gotOk := m.CheckUser(tt.args.Email, tt.args.Password)
			if gotOk != tt.wantOk {
				t.Errorf("MapUserStorage.CheckUser() got2 = %v, want %v", gotOk, tt.wantOk)
			}

			if !reflect.DeepEqual(gotId, tt.wantId) {
				require.Equal(t, gotId, tt.wantId, "The two values should be the same.")
			}
			if gotClass != tt.wantClass {
				t.Errorf("MapUserStorage.CheckUser() got1 = %v, want %v", gotClass, tt.wantClass)
			}
		})
	}
}

func TestMapUserStorage_GetResume(t *testing.T) {
	m := MapUserStorage{
		SekMu:  &sync.Mutex{},
		EmplMu: &sync.Mutex{},
		ResMu:  &sync.Mutex{},

		SeekerStorage: map[uuid.UUID]Seeker{
			uuid.MustParse("6ba7b810-9dad-11d1-80b1-00c04fd430c8"): {
				Email:      "some@mail.com",
				FirstName:  "Vova",
				SecondName: "Zyablikov",
				Password:   "1234",
				Resumes: []uuid.UUID{
					uuid.MustParse("7ba7b810-9dad-12d1-80b1-00c04fd430c8"),
				},
			},

			uuid.MustParse("6ba7b810-9bad-11d1-80b1-00c04fd430c8"): {
				Email:      "third@mail.com",
				FirstName:  "Petr",
				SecondName: "Zyablikov",
				Password:   "12345",
				Resumes: []uuid.UUID{
					uuid.MustParse("7ba7b810-9dad-11d1-71b5-04c04fd430c8"),
				},
			},
			uuid.MustParse("6ba6b810-9bad-11d1-80b2-00c04fd430c8"): {
				Email:      "some_another@mail.com",
				FirstName:  "Petya",
				SecondName: "Zyablikov",
				Password:   "12345",
				Resumes:    make([]uuid.UUID, 0),
			},
		},
		EmployerStorage: map[uuid.UUID]Employer{},
		ResumeStorage: map[uuid.UUID]Resume{
			uuid.MustParse("7ba7b810-9dad-12d1-80b1-00c04fd430c8"): {
				OwnerID:     uuid.MustParse("6ba7b810-9dad-11d1-80b1-00c04fd430c8"),
				FirstName:   "Vova",
				SecondName:  "Zyablikov",
				City:        "Moscow",
				PhoneNumber: "12345678910",
				BirthDate:   "1994-21-08",
				Sex:         "male",
				Citizenship: "Russia",
				Experience:  "7 years",
				Profession:  "programmer",
				Position:    "middle",
				Wage:        "100500",
				Education:   "MSU",
				About:       "Hello employer",
			},
			uuid.MustParse("7aa7b810-9dad-11d1-72b5-04c04fd430c8"): {
				OwnerID:     uuid.MustParse("6ba7b810-9bad-11d1-80b1-00c04fd430c8"),
				FirstName:   "Petr",
				SecondName:  "Zyablikov",
				City:        "Moscow",
				PhoneNumber: "12345678910",
				BirthDate:   "1994-21-08",
				Sex:         "male",
				Citizenship: "Russia",
				Experience:  "8 years",
				Profession:  "programmer",
				Position:    "senior",
				Wage:        "100500",
				Education:   "MSU",
				About:       "Hello employer",
			},
		},
	}
	type args struct {
		id uuid.UUID
	}
	tests := []struct {
		name string
		args args
		want Resume
	}{
		{
			name: "Test1",
			args: args{
				uuid.MustParse("7ba7b810-9dad-12d1-80b1-00c04fd430c8"),
			},
			want: Resume{
				OwnerID:     uuid.MustParse("6ba7b810-9dad-11d1-80b1-00c04fd430c8"),
				FirstName:   "Vova",
				SecondName:  "Zyablikov",
				City:        "Moscow",
				PhoneNumber: "12345678910",
				BirthDate:   "1994-21-08",
				Sex:         "male",
				Citizenship: "Russia",
				Experience:  "7 years",
				Profession:  "programmer",
				Position:    "middle",
				Wage:        "100500",
				Education:   "MSU",
				About:       "Hello employer",
			},
		},
		{
			name: "Test2",
			args: args{
				uuid.MustParse("7aa7b810-9dad-11d1-72b5-04c04fd430c8"),
			},
			want: Resume{
				OwnerID:     uuid.MustParse("6ba7b810-9bad-11d1-80b1-00c04fd430c8"),
				FirstName:   "Petr",
				SecondName:  "Zyablikov",
				City:        "Moscow",
				PhoneNumber: "12345678910",
				BirthDate:   "1994-21-08",
				Sex:         "male",
				Citizenship: "Russia",
				Experience:  "8 years",
				Profession:  "programmer",
				Position:    "senior",
				Wage:        "100500",
				Education:   "MSU",
				About:       "Hello employer",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, _ := m.GetResume(tt.args.id)

			if !reflect.DeepEqual(got, tt.want) {
				require.Equal(t, tt.want, got, "The two values should be the same.")
			}
		})
	}
}

func TestMapUserStorage_CreateResume(t *testing.T) {
	m := MapUserStorage{
		SekMu:  &sync.Mutex{},
		EmplMu: &sync.Mutex{},
		ResMu:  &sync.Mutex{},

		SeekerStorage: map[uuid.UUID]Seeker{
			uuid.MustParse("6ba7b810-9dad-11d1-80b1-00c04fd430c8"): {
				Email:      "some@mail.com",
				FirstName:  "Vova",
				SecondName: "Zyablikov",
				Password:   "1234",
				Resumes:    make([]uuid.UUID, 0),
			},

			uuid.MustParse("6ba7b810-9bad-11d1-80b1-00c04fd430c8"): {
				Email:      "third@mail.com",
				FirstName:  "Petr",
				SecondName: "Zyablikov",
				Password:   "12345",
				Resumes:    make([]uuid.UUID, 0),
			},

			uuid.MustParse("6ba6b810-9bad-11d1-80b2-00c04fd430c8"): {
				Email:      "some_another@mail.com",
				FirstName:  "Petya",
				SecondName: "Zyablikov",
				Password:   "12345",
				Resumes:    make([]uuid.UUID, 0),
			},
		},
		EmployerStorage: map[uuid.UUID]Employer{},
		ResumeStorage:   map[uuid.UUID]Resume{},
	}

	type args struct {
		resumeReg Resume
		userId    uuid.UUID
	}
	tests := []struct {
		name       string
		args       args
		wantResume Resume
		wantOk     bool
	}{
		{
			name: "Test1",
			args: args{
				resumeReg: Resume{
					OwnerID:     uuid.UUID{},
					FirstName:   "Vova",
					SecondName:  "Zyablikov",
					City:        "Moscow",
					PhoneNumber: "12345678910",
					BirthDate:   "1994-21-08",
					Sex:         "male",
					Citizenship: "Russia",
					Experience:  "7 years",
					Profession:  "programmer",
					Position:    "middle",
					Wage:        "100500",
					Education:   "MSU",
					About:       "Hello employer",
				},
				userId: uuid.MustParse("6ba7b810-9dad-11d1-80b1-00c04fd430c8"),
			},
			wantResume: Resume{
				OwnerID:     uuid.MustParse("6ba7b810-9dad-11d1-80b1-00c04fd430c8"),
				FirstName:   "Vova",
				SecondName:  "Zyablikov",
				City:        "Moscow",
				PhoneNumber: "12345678910",
				BirthDate:   "1994-21-08",
				Sex:         "male",
				Citizenship: "Russia",
				Experience:  "7 years",
				Profession:  "programmer",
				Position:    "middle",
				Wage:        "100500",
				Education:   "MSU",
				About:       "Hello employer",
			},
			wantOk: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wantId, gotOk := m.CreateResume(tt.args.resumeReg, tt.args.userId)

			if gotOk != tt.wantOk {
				t.Errorf("MapUserStorage.CreateSeeker() got1 = %v, want %v", gotOk, tt.wantOk)
			}

			if gotOk == tt.wantOk && gotOk == false {
				return
			}

			gotResume, _ := m.GetResume(wantId)
			if !reflect.DeepEqual(gotResume, tt.wantResume) {
				require.Equal(t, tt.wantResume, gotResume, "The two values should be the same.")
			}
		})
	}
}
