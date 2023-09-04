package test

import (
	"fmt"
	"github.com/google/uuid"
)

func NumberUUID[Source string | int](nbr Source) uuid.UUID {
	src, ok := any(nbr).(string)
	if !ok {
		src = fmt.Sprintf("%v", nbr)
	}

	switch len(src) {
	case 1:
		return uuid.MustParse(fmt.Sprintf("0%[1]s0%[1]s0%[1]s0%[1]s-0%[1]s0%[1]s-0%[1]s0%[1]s-0%[1]s0%[1]s-0%[1]s0%[1]s0%[1]s0%[1]s0%[1]s0%[1]s", src))
	case 2:
		return uuid.MustParse(fmt.Sprintf("%[1]s%[1]s%[1]s%[1]s-%[1]s%[1]s-%[1]s%[1]s-%[1]s%[1]s-%[1]s%[1]s%[1]s%[1]s%[1]s%[1]s", src))
	case 3:
		return uuid.MustParse(fmt.Sprintf("0%[1]s0%[1]s-0%[1]s-0%[1]s-0%[1]s-0%[1]s0%[1]s0%[1]s", src))
	case 4:
		return uuid.MustParse(fmt.Sprintf("%[1]s%[1]s-%[1]s-%[1]s-%[1]s-%[1]s%[1]s%[1]s", src))
	default:
		panic("uuid number must be between 1 and 4 characters long")
	}
}
