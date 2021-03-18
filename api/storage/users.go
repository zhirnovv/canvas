package storage

import (
	"github.com/google/uuid"
	"hash"
)

// User represents a chat user.
type User struct {
	id   uuid.UUID
	name string
}

type UserStorage struct {
	users map[uuid.UUID]*User
}

func (storage *UserStorage) Create(name string) (User, error) {
	userUuid := uuid.New()
	newUser := &User{userUuid, name}

	storage.users[userUuid] = newUser

	return *newUser, nil
}

// func (storage *UserStorage) Read(id uuid.UUID) (User, error) {
// 	user, userExists := *storage.users[id]

// 	if userExists {
// 		return user, nil
// 	}

// 	return nil, 
// }
