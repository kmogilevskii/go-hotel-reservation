package api

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/kmogilevskii/hotel-reservation/db"
	custom_errors "github.com/kmogilevskii/hotel-reservation/errors"
	"github.com/kmogilevskii/hotel-reservation/types"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserHandler struct {
	userStore db.UserStore
}

func NewUserHandler(userStore db.UserStore) *UserHandler {
	return &UserHandler{userStore: userStore}
}

func (h *UserHandler) HandleGetUsers(c *fiber.Ctx) error {
	var params db.UserQueryParams
	if err := c.QueryParser(&params); err != nil {
		return custom_errors.ErrBadRequest()
	}
	filter := db.Map{}
	if params.FirstName != "" {
		filter["firstName"] = params.FirstName
	}
	users, err := h.userStore.GetUsers(c.Context(), filter, &params.Pagination)

	if err != nil {
		return err
	}

	resp := db.ResourceResponse{
		Results: len(users),
		Data:    users,
		Page:    params.Page,
	}

	return c.JSON(resp)
}

func (h *UserHandler) HandleGetUser(c *fiber.Ctx) error {
	id := c.Params("id")

	user, err := h.userStore.GetUserByID(c.Context(), id)

	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return custom_errors.ErrResourceNotFound("user")
		}
		return err
	}

	return c.JSON(user)
}

func (h *UserHandler) HandlePostUser(c *fiber.Ctx) error {
	var params types.CreateUserParams
	if err := c.BodyParser(&params); err != nil {
		return custom_errors.ErrBadRequest()
	}
	if err := params.Validate(c.Context()); err != nil {
		return custom_errors.ErrBadRequest()
	}

	user, err := types.NewUserFromParams(params)
	if err != nil {
		return custom_errors.ErrBadRequest()
	}
	insertedUser, err := h.userStore.CreateUser(c.Context(), user)
	if err != nil {
		return err
	}
	return c.Status(http.StatusCreated).JSON(insertedUser)
}

func (h *UserHandler) HandleDeleteUser(c *fiber.Ctx) error {
	id := c.Params("id")
	if err := h.userStore.DeleteUser(c.Context(), id); err != nil {
		return custom_errors.NewError(http.StatusInternalServerError, "couldn't delete given user")
	}
	return c.JSON(map[string]string{"msg": fmt.Sprintf("deleted user with id %s", id)})
}

func (h *UserHandler) HandlePutUser(c *fiber.Ctx) error {
	var (
		params types.UpdateUserParams
		userID = c.Params("id")
	)
	if err := c.BodyParser(&params); err != nil {
		return custom_errors.ErrBadRequest()
	}
	if err := h.userStore.UpdateUser(c.Context(), userID, params); err != nil {
		return err
	}
	return c.JSON(map[string]string{"msg": fmt.Sprintf("updated user with id %s", userID)})
}
