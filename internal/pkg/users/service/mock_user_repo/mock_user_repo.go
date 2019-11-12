// Automatically generated by MockGen. DO NOT EDIT!
// Source: repository.go

package mock_users

import (
	. "2019_2_IBAT/internal/pkg/interfaces"

	gomock "github.com/golang/mock/gomock"
	uuid "github.com/google/uuid"
)

// Mock of Repository interface
type MockRepository struct {
	ctrl     *gomock.Controller
	recorder *_MockRepositoryRecorder
}

// Recorder for MockRepository (not exported)
type _MockRepositoryRecorder struct {
	mock *MockRepository
}

func NewMockRepository(ctrl *gomock.Controller) *MockRepository {
	mock := &MockRepository{ctrl: ctrl}
	mock.recorder = &_MockRepositoryRecorder{mock}
	return mock
}

func (_m *MockRepository) EXPECT() *_MockRepositoryRecorder {
	return _m.recorder
}

func (_m *MockRepository) CreateEmployer(seekerInput Employer) bool {
	ret := _m.ctrl.Call(_m, "CreateEmployer", seekerInput)
	ret0, _ := ret[0].(bool)
	return ret0
}

func (_mr *_MockRepositoryRecorder) CreateEmployer(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "CreateEmployer", arg0)
}

func (_m *MockRepository) CreateSeeker(seekerInput Seeker) bool {
	ret := _m.ctrl.Call(_m, "CreateSeeker", seekerInput)
	ret0, _ := ret[0].(bool)
	return ret0
}

func (_mr *_MockRepositoryRecorder) CreateSeeker(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "CreateSeeker", arg0)
}

func (_m *MockRepository) CreateResume(resumeReg Resume) bool {
	ret := _m.ctrl.Call(_m, "CreateResume", resumeReg)
	ret0, _ := ret[0].(bool)
	return ret0
}

func (_mr *_MockRepositoryRecorder) CreateResume(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "CreateResume", arg0)
}

func (_m *MockRepository) CreateVacancy(vacancyReg Vacancy) bool {
	ret := _m.ctrl.Call(_m, "CreateVacancy", vacancyReg)
	ret0, _ := ret[0].(bool)
	return ret0
}

func (_mr *_MockRepositoryRecorder) CreateVacancy(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "CreateVacancy", arg0)
}

func (_m *MockRepository) CreateRespond(respond Respond, userId uuid.UUID) bool {
	ret := _m.ctrl.Call(_m, "CreateRespond", respond, userId)
	ret0, _ := ret[0].(bool)
	return ret0
}

func (_mr *_MockRepositoryRecorder) CreateRespond(arg0, arg1 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "CreateRespond", arg0, arg1)
}

func (_m *MockRepository) CreateFavorite(favVac FavoriteVacancy) bool {
	ret := _m.ctrl.Call(_m, "CreateFavorite", favVac)
	ret0, _ := ret[0].(bool)
	return ret0
}

func (_mr *_MockRepositoryRecorder) CreateFavorite(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "CreateFavorite", arg0)
}

func (_m *MockRepository) DeleteUser(id uuid.UUID) error {
	ret := _m.ctrl.Call(_m, "DeleteUser", id)
	ret0, _ := ret[0].(error)
	return ret0
}

func (_mr *_MockRepositoryRecorder) DeleteUser(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "DeleteUser", arg0)
}

func (_m *MockRepository) DeleteResume(id uuid.UUID) error {
	ret := _m.ctrl.Call(_m, "DeleteResume", id)
	ret0, _ := ret[0].(error)
	return ret0
}

func (_mr *_MockRepositoryRecorder) DeleteResume(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "DeleteResume", arg0)
}

func (_m *MockRepository) DeleteVacancy(id uuid.UUID) error {
	ret := _m.ctrl.Call(_m, "DeleteVacancy", id)
	ret0, _ := ret[0].(error)
	return ret0
}

func (_mr *_MockRepositoryRecorder) DeleteVacancy(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "DeleteVacancy", arg0)
}

func (_m *MockRepository) CheckUser(email string, password string) (uuid.UUID, string, bool) {
	ret := _m.ctrl.Call(_m, "CheckUser", email, password)
	ret0, _ := ret[0].(uuid.UUID)
	ret1, _ := ret[1].(string)
	ret2, _ := ret[2].(bool)
	return ret0, ret1, ret2
}

func (_mr *_MockRepositoryRecorder) CheckUser(arg0, arg1 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "CheckUser", arg0, arg1)
}

func (_m *MockRepository) PutSeeker(seekerInput SeekerReg, id uuid.UUID) bool {
	ret := _m.ctrl.Call(_m, "PutSeeker", seekerInput, id)
	ret0, _ := ret[0].(bool)
	return ret0
}

func (_mr *_MockRepositoryRecorder) PutSeeker(arg0, arg1 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "PutSeeker", arg0, arg1)
}

func (_m *MockRepository) PutEmployer(employerInput EmployerReg, id uuid.UUID) bool {
	ret := _m.ctrl.Call(_m, "PutEmployer", employerInput, id)
	ret0, _ := ret[0].(bool)
	return ret0
}

func (_mr *_MockRepositoryRecorder) PutEmployer(arg0, arg1 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "PutEmployer", arg0, arg1)
}

func (_m *MockRepository) PutResume(resume Resume, userId uuid.UUID, resumeId uuid.UUID) bool {
	ret := _m.ctrl.Call(_m, "PutResume", resume, userId, resumeId)
	ret0, _ := ret[0].(bool)
	return ret0
}

func (_mr *_MockRepositoryRecorder) PutResume(arg0, arg1, arg2 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "PutResume", arg0, arg1, arg2)
}

func (_m *MockRepository) PutVacancy(vacavcy Vacancy, userId uuid.UUID, resumeId uuid.UUID) bool {
	ret := _m.ctrl.Call(_m, "PutVacancy", vacavcy, userId, resumeId)
	ret0, _ := ret[0].(bool)
	return ret0
}

func (_mr *_MockRepositoryRecorder) PutVacancy(arg0, arg1, arg2 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "PutVacancy", arg0, arg1, arg2)
}

func (_m *MockRepository) GetEmployers(params map[string]interface{}) ([]Employer, error) {
	ret := _m.ctrl.Call(_m, "GetEmployers", params)
	ret0, _ := ret[0].([]Employer)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockRepositoryRecorder) GetEmployers(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "GetEmployers", arg0)
}

func (_m *MockRepository) GetSeekers() ([]Seeker, error) {
	ret := _m.ctrl.Call(_m, "GetSeekers")
	ret0, _ := ret[0].([]Seeker)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockRepositoryRecorder) GetSeekers() *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "GetSeekers")
}

func (_m *MockRepository) GetResumes(params map[string]interface{}) ([]Resume, error) {
	ret := _m.ctrl.Call(_m, "GetResumes", params)
	ret0, _ := ret[0].([]Resume)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockRepositoryRecorder) GetResumes(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "GetResumes", arg0)
}

func (_m *MockRepository) GetVacancies(params map[string]interface{}) ([]Vacancy, error) {
	ret := _m.ctrl.Call(_m, "GetVacancies", params)
	ret0, _ := ret[0].([]Vacancy)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockRepositoryRecorder) GetVacancies(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "GetVacancies", arg0)
}

func (_m *MockRepository) GetResponds(record AuthStorageValue, params map[string]string) ([]Respond, error) {
	ret := _m.ctrl.Call(_m, "GetResponds", record, params)
	ret0, _ := ret[0].([]Respond)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockRepositoryRecorder) GetResponds(arg0, arg1 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "GetResponds", arg0, arg1)
}

func (_m *MockRepository) GetFavoriteVacancies(record AuthStorageValue) ([]Vacancy, error) {
	ret := _m.ctrl.Call(_m, "GetFavoriteVacancies", record)
	ret0, _ := ret[0].([]Vacancy)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockRepositoryRecorder) GetFavoriteVacancies(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "GetFavoriteVacancies", arg0)
}

func (_m *MockRepository) GetSeeker(id uuid.UUID) (Seeker, error) {
	ret := _m.ctrl.Call(_m, "GetSeeker", id)
	ret0, _ := ret[0].(Seeker)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockRepositoryRecorder) GetSeeker(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "GetSeeker", arg0)
}

func (_m *MockRepository) GetEmployer(id uuid.UUID) (Employer, error) {
	ret := _m.ctrl.Call(_m, "GetEmployer", id)
	ret0, _ := ret[0].(Employer)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockRepositoryRecorder) GetEmployer(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "GetEmployer", arg0)
}

func (_m *MockRepository) GetResume(id uuid.UUID) (Resume, error) {
	ret := _m.ctrl.Call(_m, "GetResume", id)
	ret0, _ := ret[0].(Resume)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockRepositoryRecorder) GetResume(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "GetResume", arg0)
}

func (_m *MockRepository) GetVacancy(id uuid.UUID) (Vacancy, error) {
	ret := _m.ctrl.Call(_m, "GetVacancy", id)
	ret0, _ := ret[0].(Vacancy)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockRepositoryRecorder) GetVacancy(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "GetVacancy", arg0)
}

func (_m *MockRepository) SetImage(id uuid.UUID, class string, imageName string) bool {
	ret := _m.ctrl.Call(_m, "SetImage", id, class, imageName)
	ret0, _ := ret[0].(bool)
	return ret0
}

func (_mr *_MockRepositoryRecorder) SetImage(arg0, arg1, arg2 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "SetImage", arg0, arg1, arg2)
}
