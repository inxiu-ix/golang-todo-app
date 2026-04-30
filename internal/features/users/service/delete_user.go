package users_service

import (
	"context"
	"fmt"
)

func (s *UsersService) DeleteUser(ctx context.Context, id int) error {
	err := s.usersRepository.DeleteUser(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	return nil
}
