package models

// User represents the stored user profile data.
type User struct {
	ID              int    `json:"id"`
	Username        string `json:"username"`
	ProfileImageURL string `json:"profile_image_url"`
}

// UpdateProfilePictureRequest is the DTO for updating the profile picture.
type UpdateProfilePictureRequest struct {
	// ImageURL is the source of the external image.
	// The 'url' binding tag provides initial syntactic validation.
	ImageURL string `json:"image_url" binding:"required,url"`
}