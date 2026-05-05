package users_transport_http

import (
	"net/http"

	"github.com/inxiu-ix/golang-todo-app/internal/core/domain"
	core_logger "github.com/inxiu-ix/golang-todo-app/internal/core/logger"
	core_http_request "github.com/inxiu-ix/golang-todo-app/internal/core/transport/http/request"
	core_http_response "github.com/inxiu-ix/golang-todo-app/internal/core/transport/http/response"
)

type CreateUserRequest struct {
	FullName    string  `json:"full_name" validate:"required,min=3,max=100"    example:"Ivan Ivanov"`
	PhoneNumber *string `json:"phone_number" validate:"omitempty,e164"         example:"+79991234567"`
}

type CreateUserResponse UserDTOResponse

// CreateUser godoc
// @Summary Create a new user
// @Description Create a new user with the given full name and phone number
// @Tags users
// @Accept json
// @Produce json
// @Param request body CreateUserRequest true "Create user request"
// @Success 201 {object} CreateUserResponse
// @Failure 400 {object} core_http_response.ErrorResponse "Bad Request"
// @Failure 500 {object} core_http_response.ErrorResponse "Internal Server Error"
// @Router /users [post]
func (h *UserHTTPHandler) CreateUser(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)

	responseHandler := core_http_response.NewHTTPResponseHandler(log, rw)

	log.Debug("CreateUser request received")

	var request CreateUserRequest

	if err := core_http_request.DecodeAndValidateRequest(r, &request); err != nil {
		responseHandler.ErrorResponse(err, "failed to decode and validate HTTP request")

		return
	}

	userDomain := domainFromDTO(request)

	userDomain, err := h.usersService.CreateUser(ctx, userDomain)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to create user")

		return
	}

	response := CreateUserResponse(userDTOFromDomain(userDomain))

	responseHandler.JSONResponse(response, http.StatusCreated)
}

func domainFromDTO(dto CreateUserRequest) domain.User {
	return domain.NewUserUninitialized(dto.FullName, dto.PhoneNumber)
}
