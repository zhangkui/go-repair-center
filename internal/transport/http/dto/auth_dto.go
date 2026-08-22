package dto

type RegisterRequest struct { Username string `json:"username" validate:"required,min=3,max=64"`; Password string `json:"password" validate:"required,min=8,max=72"`; DisplayName string `json:"display_name" validate:"required,max=100"` }
type LoginRequest struct { Username string `json:"username" validate:"required"`; Password string `json:"password" validate:"required"` }
type RefreshRequest struct { RefreshToken string `json:"refresh_token" validate:"required"` }
type ChangePasswordRequest struct { OldPassword string `json:"old_password" validate:"required"`; NewPassword string `json:"new_password" validate:"required,min=8,max=72"` }
type ResetPasswordRequest struct { NewPassword string `json:"new_password" validate:"required,min=8,max=72"` }