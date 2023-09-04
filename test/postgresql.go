package test

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/uptrace/bun"
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

func RunTransactionalTest[Fixtures any](db bun.IDB, fixtures []Fixtures, call func(ctx context.Context, tx bun.Tx)) error {
	tx, err := db.BeginTx(context.TODO(), &sql.TxOptions{})
	if err != nil {
		return err
	}
	defer tx.Rollback()

	for _, data := range fixtures {
		if _, err := tx.NewInsert().Model(data).Exec(context.TODO()); err != nil {
			mrsh, _ := json.Marshal(data)
			return fmt.Errorf("failed to insert data %s: %w", string(mrsh), err)
		}
	}

	call(context.TODO(), tx)
	return nil
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
