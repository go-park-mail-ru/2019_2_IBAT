package users

import (
	"context"
	"encoding/json"

	"io"
	"io/ioutil"
	"log"

	"2019_2_IBAT/pkg/app/recommends/recomsproto"
	. "2019_2_IBAT/pkg/pkg/interfaces"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

func (h *UserService) CreateVacancy(body io.ReadCloser, authInfo AuthStorageValue) (uuid.UUID, error) {

	if authInfo.Role != EmployerStr {
		// log.Printf("Invalid action: %s", err)
		return uuid.UUID{}, errors.New(ForbiddenMsg)
	}

	bytes, err := ioutil.ReadAll(body)
	if err != nil {
		log.Printf("error while reading body: %s", err)
		err = errors.Wrap(err, "reading body error")
		return uuid.UUID{}, errors.New(BadRequestMsg)
	}

	var vacancyReg Vacancy
	// id := uuid.New()
	// vacancyReg.ID = id
	// vacancyReg.OwnerID = authInfo.ID

	err = json.Unmarshal(bytes, &vacancyReg)
	if err != nil {
		log.Printf("Error while unmarshaling: %s", err)
		err = errors.Wrap(err, "unmarshaling error")
		return uuid.UUID{}, errors.New(InvalidJSONMsg)
	}

	id := uuid.New()
	vacancyReg.ID = id
	vacancyReg.OwnerID = authInfo.ID

	ok := h.Storage.CreateVacancy(vacancyReg)

	if !ok {
		log.Printf("Error while creating vacancy: %s", err)
		return uuid.UUID{}, errors.New(BadRequestMsg)
	}

	// tagIDs := h.Storage.GetTagIDs(vacancyReg.Spheres)

	// h.NotifChan <- NotifStruct{
	// 	VacancyId: id,
	// 	TagIDs:    tagIDs,
	// }

	return id, nil
}

func (h *UserService) GetVacancy(vacancyId uuid.UUID, authInfo AuthStorageValue) (Vacancy, error) {
	// wg := &sync.WaitGroup{}
	// wg.Add(1)

	tagIDs, err := h.Storage.GetVacancyTagIDs(vacancyId)

	ctx := context.Background()
	h.RecomService.SetTagIDs(
		ctx,
		&recomsproto.SetTagIDsMessage{
			ID:      authInfo.ID.String(),
			Role:    authInfo.Role,
			Expires: authInfo.Expires,
			IDs:     UuidsToStrings(tagIDs),
		},
	)

	vacancy, err := h.Storage.GetVacancy(vacancyId, authInfo.ID)

	if err != nil { //error wrap
		return vacancy, errors.New(InvalidIdMsg)
	}

	return vacancy, nil
}

func (h *UserService) DeleteVacancy(vacancyId uuid.UUID, authInfo AuthStorageValue) error {
	if authInfo.Role != EmployerStr {
		return errors.New(ForbiddenMsg)
	}

	vacancy, err := h.Storage.GetVacancy(vacancyId, authInfo.ID)

	if err != nil {
		return errors.New(InvalidIdMsg)
	}

	if vacancy.OwnerID != authInfo.ID { //error wrap
		return errors.New(ForbiddenMsg)
	}

	err = h.Storage.DeleteVacancy(vacancyId)

	if err != nil {
		return errors.New(InternalErrorMsg)
	}

	return nil
}

func (h *UserService) PutVacancy(vacancyId uuid.UUID, body io.ReadCloser, authInfo AuthStorageValue) error {
	if authInfo.Role != EmployerStr {
		// log.Printf("Invalid action: %s", err)
		return errors.New(ForbiddenMsg)
	}

	bytes, err := ioutil.ReadAll(body)
	if err != nil {
		return errors.New(BadRequestMsg)
	}

	var vacancy Vacancy
	err = json.Unmarshal(bytes, &vacancy)
	if err != nil {
		log.Printf("Error while unmarshaling: %s", err)
		err = errors.New(InvalidJSONMsg)
		return err
	}

	ok := h.Storage.PutVacancy(vacancy, authInfo.ID, vacancyId)

	if !ok {
		log.Printf("Error while changing vacancy")
		return errors.New(BadRequestMsg)
	}

	return nil
}

func (h *UserService) GetVacancies(authInfo AuthStorageValue, params map[string]interface{},
	tagParams map[string]interface{}) ([]Vacancy, error) {
	var tagIDs []uuid.UUID
	var err error

	if params["recommended"] != nil {
		ctx := context.Background()
		tagIDsMsg, err := h.RecomService.GetTagIDs(
			ctx,
			&recomsproto.GetTagIDsMessage{
				ID:      authInfo.ID.String(),
				Role:    authInfo.Role,
				Expires: authInfo.Expires,
			})

		tagIDs = StringsToUuids(tagIDsMsg.IDs)
		if err != nil {
			log.Printf("Failed to get recommendations")
		}
	} else {
		var tags []Pair
		for keyTag, item := range tagParams {
			arr := item.([]string)
			for _, tag := range arr {
				tags = append(tags, Pair{
					First:  keyTag,
					Second: tag,
				})
			}
		}
		tagIDs = h.Storage.GetTagIDs(tags)
	}

	params["tag_ids"] = tagIDs

	vacancies, err := h.Storage.GetVacancies(authInfo, params)

	if err != nil {
		return vacancies, err
	}

	if params["recommended"] == nil {
		ctx := context.Background()
		_, err = h.RecomService.SetTagIDs(
			ctx,
			&recomsproto.SetTagIDsMessage{
				ID:      authInfo.ID.String(),
				Role:    authInfo.Role,
				Expires: authInfo.Expires,
				IDs:     UuidsToStrings(tagIDs),
			})
		if err != nil {
			log.Printf("Failed to set recommendations")
		}
	}

	return vacancies, err
}
