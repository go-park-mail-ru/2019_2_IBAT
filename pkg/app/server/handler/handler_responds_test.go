package handler

import (
	mock_users "2019_2_IBAT/pkg/app/server/handler/mock_users"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	. "2019_2_IBAT/pkg/pkg/models"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/require"
)

// func (h *Handler) CreateRespond(w http.ResponseWriter, r *http.Request) { //+
// 	w.Header().Set("Content-Type", "application/json; charset=UTF-8")s

// 	authInfo, ok := FromContext(r.Context())
// 	if !ok {
// 		SetError(w, http.StatusUnauthorized, UnauthorizedMsg)
// 		return
// 	}

// 	err := h.UserService.CreateRespond(r.Body, authInfo)
// 	if err != nil {
// 		SetError(w, http.StatusForbidden, ForbiddenMsg)
// 		return
// 	}
// }

func TestHandler_CreateRespond(t *testing.T) {
	mockCtrl1 := gomock.NewController(t)
	defer mockCtrl1.Finish()

	mockUserService := mock_users.NewMockService(mockCtrl1)

	h := Handler{
		UserService: mockUserService,
	}

	tests := []struct {
		name             string
		wantFail         bool
		wantStatusCode   int
		wantErrorMessage string
		wantForbidden    bool
		wantUnauth       bool
		ctx              context.Context
		respond          Respond
		record           AuthStorageValue
	}{
		{
			name: "Test1",
			ctx:  context.Background(),
			respond: Respond{
				Status:    AwaitSt,
				ResumeID:  uuid.New(),
				VacancyID: uuid.New(),
			},
			record: AuthStorageValue{
				ID:   uuid.New(),
				Role: SeekerStr,
			},
		},
		{
			name:             "Test2",
			wantFail:         true,
			ctx:              context.Background(),
			wantUnauth:       true,
			wantStatusCode:   http.StatusUnauthorized,
			wantErrorMessage: UnauthorizedMsg,
		},
		{
			name:             "Test3",
			wantForbidden:    true,
			wantFail:         true,
			wantStatusCode:   http.StatusForbidden,
			wantErrorMessage: ForbiddenMsg,
			respond: Respond{
				Status:    AwaitSt,
				ResumeID:  uuid.New(),
				VacancyID: uuid.New(),
			},
			record: AuthStorageValue{
				ID:   uuid.New(),
				Role: SeekerStr,
			},
		},

		// {
		// 	name:             "Test4",
		// 	role:             EmployerStr,
		// 	userId:           uuid.UUID{},
		// 	wantFail:         true,
		// 	wantForbidden:    true,
		// 	wantStatusCode:   http.StatusForbidden,
		// 	wantErrorMessage: ForbiddenMsg,
		// 	ctx:              context.Background(),
		// },
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var str string
			var wantJSON []byte
			wantJSON, _ = json.Marshal(tc.respond)
			str = string(wantJSON)

			reader := strings.NewReader(string(str))

			req, err := http.NewRequest("POST", "/respond", reader)
			if err != nil {
				t.Fatal(err)
			}
			req.Header.Set("Content-Type", "application/json")

			if !tc.wantFail && !tc.wantUnauth {
				mockUserService.
					EXPECT().
					CreateRespond(gomock.Any(), tc.record).
					Return(nil)
			} else if tc.wantForbidden {
				mockUserService.
					EXPECT().
					CreateRespond(gomock.Any(), tc.record).
					Return(fmt.Errorf(ForbiddenMsg))
			}

			rr := httptest.NewRecorder()

			ctx := context.TODO()
			if !tc.wantUnauth {
				ctx = NewContext(req.Context(), tc.record)
			}

			router := mux.NewRouter()
			router.HandleFunc("/respond", h.CreateRespond)
			router.ServeHTTP(rr, req.WithContext(ctx))

			if !tc.wantFail {
				if rr.Code != http.StatusOK {
					t.Error("status is not ok")
				}
			} else {
				bytes, _ := ioutil.ReadAll(rr.Body)
				var gotError Error
				_ = gotError.UnmarshalJSON(bytes)

				require.Equal(t, tc.wantStatusCode, rr.Code, "The two values should be the same.")
				require.Equal(t, tc.wantErrorMessage, gotError.Message, "The two values should be the same.")
			}
		})
	}
}
