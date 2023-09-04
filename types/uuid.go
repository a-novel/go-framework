package types

import (
	"github.com/google/uuid"
	"strings"
)

// https://github.com/gin-gonic/gin/issues/1516#issuecomment-1269846541

type StringUUID string

func (s StringUUID) Value() uuid.UUID {
	parsed, err := uuid.Parse(string(s))
	if err != nil {
		return uuid.Nil
	}

	return parsed
}

type StringUUIDs string

func (s StringUUIDs) Value() []uuid.UUID {
	var uuids []uuid.UUID

	for _, id := range strings.Split(string(s), ",") {
		parsed, err := uuid.Parse(id)
		if err != nil {
			continue
		}

		uuids = append(uuids, parsed)
	}

	return uuids
}
