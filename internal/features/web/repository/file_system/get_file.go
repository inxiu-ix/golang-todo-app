package web_fs_repository

import (
	"fmt"
	"os"

	core_errors "github.com/inxiu-ix/golang-todo-app/internal/core/errors"
)

func (r *WebFSRepository) GetFile(filePath string) ([]byte, error) {
	file, err := os.ReadFile(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("file not found: %w", core_errors.ErrNotFound)
		}
		return nil, fmt.Errorf("read file: %w", err)
	} 

	return file, nil
}
