package users_transport_http

import "github.com/inxiu-ix/golang-todo-app/internal/core/domain"

type UserDTOResponse struct {
	ID          int     `json:"id"`
	Version     int     `json:"version"`
	FullName    string  `json:"full_name"`
	PhoneNumber *string `json:"phone_number"`
}

func userDTOFromDomain(domainUser domain.User) UserDTOResponse {
	return UserDTOResponse{
		ID:          domainUser.ID,
		Version:     domainUser.Version,
		FullName:    domainUser.FullName,
		PhoneNumber: domainUser.PhoneNumber,
	}
}

func usersDTOFromDomains(users []domain.User) []UserDTOResponse {
	usersDTO := make([]UserDTOResponse, len(users))

	for i, user := range users {
		usersDTO[i] = userDTOFromDomain(user)
	}

	return usersDTO
}
