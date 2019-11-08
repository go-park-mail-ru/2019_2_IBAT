package users

import (
	. "2019_2_IBAT/internal/pkg/interfaces"
	"encoding/json"
	"io"
	"io/ioutil"
	"log"

	"github.com/pkg/errors"
)

func (h *UserService) CreateFavorite(body io.ReadCloser, record AuthStorageValue) error { //should do this part by one r with if?
	if record.Role != SeekerStr {
		return errors.New("Invalid action")
	}

	bytes, err := ioutil.ReadAll(body)
	if err != nil {
		log.Printf("error while reading body: %s", err)
		err = errors.Wrap(err, "reading body error")
		return err
	}

	var favVac FavoriteVacancy
	err = json.Unmarshal(bytes, &favVac)
	if err != nil {
		log.Printf("Error while unmarshaling: %s", err)
		err = errors.Wrap(err, "unmarshaling error")
		return err
	}

	favVac.PersonID = record.ID
	ok := h.Storage.CreateFavorite(favVac)

	if !ok {
		log.Printf("Error while creating favorite_vacancy: %s", err)
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
