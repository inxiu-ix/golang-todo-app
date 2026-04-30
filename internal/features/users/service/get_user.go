package users_service

import (
	"context"
	"fmt"

	"github.com/inxiu-ix/golang-todo-app/internal/core/domain"
)

func (s *UsersService) GetUser(
	ctx context.Context,
	id int,
) (domain.User, error) {
	user, err := s.usersRepository.GetUser(ctx, id)

	if err != nil {
		return domain.User{}, fmt.Errorf("failed to get user: %w", err)
	}

	return user, nil

}
