package interfaces

import "github.com/google/uuid"

const SeekerStr = "seeker"
const EmployerStr = "employer"
const DefaultImg = "default.jpg"

func UuidsToStrings(ids []uuid.UUID) []string {
	var strIDs []string
	for _, id := range ids {
		strIDs = append(strIDs, id.String())
	}
	return strIDs
}

func StringsToUuids(strIDs []string) []uuid.UUID {
	var ids []uuid.UUID
	for _, id := range strIDs {
		ids = append(ids, uuid.MustParse(id))
	}
	return ids
}
