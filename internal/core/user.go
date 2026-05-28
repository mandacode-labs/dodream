package core

import "time"

type UserID string

func (id UserID) String() string {
	return string(id)
}

type User struct {
	id         UserID
	nickname   string
	providerID string
	createdAt  time.Time
	updatedAt  time.Time
}

func NewUser(id UserID, nickname string, providerID string) *User {
	now := time.Now()
	return &User{
		id:         id,
		nickname:   nickname,
		providerID: providerID,
		createdAt:  now,
		updatedAt:  now,
	}
}

func (u *User) ID() UserID {
	return u.id
}

func (u *User) Nickname() string {
	return u.nickname
}

func (u *User) ProviderID() string {
	return u.providerID
}

func (u *User) CreatedAt() time.Time {
	return u.createdAt
}

func (u *User) UpdatedAt() time.Time {
	return u.updatedAt
}

func (u *User) SetNickname(nickname string) {
	u.nickname = nickname
	u.updatedAt = time.Now()
}

func (u *User) SetProviderID(providerID string) {
	u.providerID = providerID
	u.updatedAt = time.Now()
}
