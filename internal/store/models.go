package store

import (
	"time"
)

// Profile représente la table 'profiles' de Supabase
type Profile struct {
	ID           string  `db:"id" json:"id"`
	Username     string  `db:"username" json:"username"`
	DisplayName  *string `db:"display_name" json:"display_name"` // * car peut être NULL
	AvatarURL    *string `db:"avatar_url" json:"avatar_url"`
	BannerURL    *string `db:"banner_url" json:"banner_url"`
	Bio          *string `db:"bio" json:"bio"`
	ProfileColor string  `db:"profile_color" json:"profile_color"`

	// Options Flexx
	IsFlexx        bool    `db:"is_flexx" json:"is_flexx"`
	FlexxGlowColor *string `db:"flexx_glow_color" json:"flexx_glow_color"`

	// Présence (Mapped sur l'Enum SQL)
	Presence    string  `db:"presence" json:"presence"`
	StatusText  *string `db:"status_text" json:"status_text"`
	StatusEmoji *string `db:"status_emoji" json:"status_emoji"`

	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}

// Relationship représente le lien entre deux utilisateurs (Amis/Bloqués)
type Relationship struct {
	UserID    string    `db:"user_id" json:"user_id"`
	TargetID  string    `db:"target_id" json:"target_id"`
	Status    string    `db:"status" json:"status"` // 'pending', 'friend', 'blocked'
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}
