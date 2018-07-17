package utils

import (
	"math/rand"
	"time"
)

var (
	letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890")
)

const USER_STATE_KEY = "state"

func init() {
	rand.Seed(time.Now().UTC().UnixNano())
}

func GenerateRandStr(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

// type CurrentUserManager struct {
// 	key int
// }

// func NewCurrentUserManager(key int) *CurrentUserManager {
// 	return &CurrentUserManager{key: key}
// }

// func (cu *CurrentUserManager) Set(r *http.Request, user *models.User) {
// 	*r = *r.WithContext(context.WithValue(r.Context(), cu.key, user))
// }

// func (cu *CurrentUserManager) Get(r *http.Request) *models.User {
// 	if val := r.Context().Value(cu.key); val != nil {
// 		return val.(*models.User)
// 	}
// 	return nil
// }
