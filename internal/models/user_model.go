package models

import "time"

type CreateUserRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type SignInRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type SignInResponse struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserResponse struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type User struct {
	ID        string  `json:"id"`
	Name      string  `json:"name"`
	Email     string  `json:"email"`
	Password  string  `json:"-"`
	AvatarURL *string `json:"avatar_url,omitempty"`
}

type RefreshToken struct {
	ID         string    `bson:"id,omitempty"`
	UserID     string    `bson:"user_id"`
	TokenHash  string    `bson:"token_hash"`
	ExpiresAt  time.Time `bson:"expires_at"`
	CreatedAt  time.Time `bson:"created_at"`
	LastUsedAt time.Time `bson:"last_used_at"`
	Revoked    bool      `bson:"revoked"`
}

type UpdateProfileRequest struct {
	Name      *string `json:"name,omitempty"`
	AvatarURL *string `json:"avatar_url,omitempty"`
}

type UpdateProfileResponse struct {
	ID        string  `json:"id"`
	Name      string  `json:"name"`
	AvatarURL *string `json:"avatar_url,omitempty"`
}

type ChangePasswordRequest struct {
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
}
