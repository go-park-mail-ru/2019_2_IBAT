package users

import (
	. "2019_2_IBAT/internal/pkg/interfaces"
	"log"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

func (h *UserService) CreateFavorite(vacancyId uuid.UUID, record AuthStorageValue) error { //should do this part by one r with if?
	if record.Role != SeekerStr {
		log.Println("Invalid action")
		return errors.New("Invalid action")
	}

	var favVac FavoriteVacancy
	favVac.PersonID = record.ID
	favVac.VacancyID = vacancyId
	ok := h.Storage.CreateFavorite(favVac)

	if !ok {
		log.Println("Error while creating favorite_vacancy")
		return errors.New("Error while creating favorite_vacancy")
	}

	return nil
}

func (h *UserService) GetFavoriteVacancies(authInfo AuthStorageValue) ([]Vacancy, error) {

	vacancies, err := h.Storage.GetFavoriteVacancies(authInfo)
	if err != nil {
		return vacancies, errors.New("Invalid action") ///
	}

	return vacancies, nil
}

func (h *UserService) DeleteFavoriteVacancy(vacancyId uuid.UUID, authInfo AuthStorageValue) error {

	err := h.Storage.DeleteFavoriteVacancy(vacancyId, authInfo)

	if err != nil {
		return errors.New(InternalErrorMsg)
	}

	return nil
}
