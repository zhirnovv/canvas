package user

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
)

// User represents a chat user.
type User struct {
	ID   uuid.UUID `json:"id"`   // The unique identifier of a user.
	Name string    `json:"name"` // The name of a user.
}

// UserStorage represents a storage for all users.
type UserStorage struct {
	users map[uuid.UUID]*User
}

func createUserDoesNotExistError(id uuid.UUID) error {
	return errors.New(fmt.Sprintf("User with id %s does not exist", id.String()))
}

func NewUserStorage() *UserStorage {
	return &UserStorage{users: make(map[uuid.UUID]*User)}
}

// Create() creates a new user with a specified name.
func (storage *UserStorage) Create(name string) (*User, error) {
	userUuid := uuid.New()
	newUser := &User{userUuid, name}

	storage.users[userUuid] = newUser

	return newUser, nil
}

// Read() returns a user by his id. If user with specified id does not exist an error is thrown.
func (storage *UserStorage) Read(id uuid.UUID) (*User, error) {
	user, userExists := storage.users[id]

	if userExists {
		return user, nil
	}

	return nil, createUserDoesNotExistError(id)
}

// Delete() deletes a user by his id. If user with specified id does not exist an error is thrown.
func (storage *UserStorage) Delete(id uuid.UUID) error {
	if _, userDoesNotExist := storage.Read(id); userDoesNotExist != nil {
		return createUserDoesNotExistError(id)
	}

	delete(storage.users, id)

	return nil
}
