package users

import (
	mock_user_repo "2019_2_IBAT/pkg/app/users/service/mock_user_repo"
	. "2019_2_IBAT/pkg/pkg/models"
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

// func (h *UserService) GetTags() (map[string][]string, error) {
// 	tagMap := map[string][]string{}

// 	tags, err := h.Storage.GetTags()
// 	if err != nil {
// 		return tagMap, err
// 	}

// 	for _, item := range tags {
// 		if tagMap[item.ParentTag] == nil {
// 			tagMap[item.ParentTag] = []string{}
// 		}
// 		tagMap[item.ParentTag] = append(tagMap[item.ParentTag], item.ChildTag)
// 	}

// 	return tagMap, nil
// }

func TestUserService_GetTags(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockUserRepo := mock_user_repo.NewMockRepository(mockCtrl)

	h := UserService{
		Storage: mockUserRepo,
	}

	expTags := []Tag{
		{
			ParentTag: "232",
			ChildTag:  "23232",
		},
		{
			ParentTag: "heroku",
			ChildTag:  "23232",
		},
	}

	expTagMap := map[string][]string{}
	for _, item := range expTags {
		if expTagMap[item.ParentTag] == nil {
			expTagMap[item.ParentTag] = []string{}
		}
		expTagMap[item.ParentTag] = append(expTagMap[item.ParentTag], item.ChildTag)
	}

	tests := []struct {
		name             string
		wantFail         bool
		wantRecomms      bool
		wantErrorMessage string
	}{
		{
			name: "Test1",
		},
		{
			name:             "Test2",
			wantFail:         true,
			wantErrorMessage: InternalErrorMsg,
		},
		{
			name:             "Test2",
			wantFail:         true,
			wantRecomms:      true,
			wantErrorMessage: InternalErrorMsg,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if !tt.wantFail {
				mockUserRepo.
					EXPECT().
					GetTags().
					Return(expTags, nil)
			} else {
				mockUserRepo.
					EXPECT().
					GetTags().
					Return([]Tag{}, fmt.Errorf(InternalErrorMsg))
			}

			gotTags, err := h.GetTags()

			if !tt.wantFail {
				if err != nil {
					t.Error("Error is not nil\n")
				}
				require.Equal(t, expTagMap, gotTags, "The two values should be the same.")
			} else {
				require.Equal(t, tt.wantErrorMessage, err.Error(), "The two values should be the same.")
			}
		})
	}
}
