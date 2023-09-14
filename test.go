package goframework

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"os"
	"path"
	"testing"
	"time"
)

type FileFixture struct {
	Name    string
	Content []byte
	Date    time.Time
}

func RunFileTransactionalTest(t *testing.T, fixtures []FileFixture, call func(ctx context.Context, basePath string)) error {
	dir := t.TempDir()
	for _, file := range fixtures {
		fullPath := path.Join(dir, file.Name)
		if err := os.WriteFile(fullPath, file.Content, 0644); err != nil {
			return fmt.Errorf("failed to write file %s: %w", fullPath, err)
		}
		if err := os.Chtimes(fullPath, file.Date, file.Date); err != nil {
			return fmt.Errorf("failed to update access time for file %s: %w", fullPath, err)
		}
	}

	call(context.TODO(), dir)
	return nil
}

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
