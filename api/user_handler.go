package api

import (
	"errors"

	"github.com/AyanokojiKiyotaka8/booking.git/db"
	"github.com/AyanokojiKiyotaka8/booking.git/types"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserHandler struct {
	userStore db.UserStore
}

func NewUserHandler(userStore db.UserStore) *UserHandler {
	return &UserHandler{
		userStore: userStore,
	}
}

func (h *UserHandler) HandleGetUser(c *fiber.Ctx) error {
	id := c.Params("id")
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	filter := bson.M{"_id": oid}
	user, err := h.userStore.GetUser(c.Context(), filter)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return c.JSON(map[string]string{"error:": "not found"})
		}
		return err
	}
	return c.JSON(user)
}

func (h *UserHandler) HandlePostUser(c *fiber.Ctx) error {
	var u types.CreateUserParams
	if err := c.BodyParser(&u); err != nil {
		return ErrBadRequest()
	}
	if errors := u.Validate(); len(errors) > 0 {
		return c.JSON(errors)
	}

	user, err := types.NewUserFromParams(&u)
	if err != nil {
		return err
	}

	insertedUser, err := h.userStore.InsertUser(c.Context(), user)
	if err != nil {
		return err
	}
	return c.JSON(insertedUser)
}

func (h *UserHandler) HandleGetUsers(c *fiber.Ctx) error {
	filter := bson.M{}
	users, err := h.userStore.GetUsers(c.Context(), filter)
	if err != nil {
		return ErrResourceNotFound("user")
	}
	return c.JSON(users)
}

func (h *UserHandler) HandleDeleteUser(c *fiber.Ctx) error {
	id := c.Params("id")
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	filter := bson.M{"_id": oid}
	if err := h.userStore.DeleteUser(c.Context(), filter); err != nil {
		return err
	}
	return c.JSON(map[string]string{"deleted:": id})
}

func (h *UserHandler) HandlePutUser(c *fiber.Ctx) error {
	id := c.Params("id")
	var params types.UpdateUserParams
	if err := c.BodyParser(&params); err != nil {
		return ErrBadRequest()
	}

	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return ErrInvalidID()
	}

	filter := bson.M{"_id": oid}
	update := bson.M{
		"$set": params.ToBSON(),
	}

	err = h.userStore.UpdateUser(c.Context(), filter, update)
	if err != nil {
		return err
	}
	return c.JSON(map[string]string{"updated:": id})
}
