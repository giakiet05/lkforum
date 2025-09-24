package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Username    string             `bson:"username" json:"username"`
	Email       string             `bson:"email,omitempty" json:"email,omitempty"`
	Password    string             `bson:"password" json:"password"`
	Role        string             `bson:"role" json:"role"`
	RoleContent RoleContent        `bson:"role_content,omitempty" json:"role_content,omitempty"`
	CreateAt    time.Time          `bson:"create_at,omitempty" json:"create_at,omitempty"`
}

type RoleContent struct {
	User  *UserRoleContent  `bson:"user,omitempty" json:"user,omitempty"`
	Admin *AdminRoleContent `bson:"admin,omitempty" json:"admin,omitempty"`
}

type UserRoleContent struct {
	Avatar   string `bson:"avatar,omitempty" json:"avatar,omitempty"`
	Cover    string `bson:"cover,omitempty" json:"cover,omitempty"`
	BanStart string `bson:"ban_start,omitempty" json:"ban_start,omitempty"`
	BanEnd   string `bson:"ban_end,omitempty" json:"ban_end,omitempty"`
}

type AdminRoleContent struct {
	Name        string             `bson:"name,omitempty" json:"name,omitempty"`
	Permissions []string           `bson:"permissions,omitempty" json:"permissions,omitempty"`
	CreateAt    *time.Time         `bson:"update_at,omitempty" json:"update_at,omitempty"`
	CreateBy    primitive.ObjectID `bson:"create_by,omitempty" json:"create_by,omitempty"`
}

type UserStat struct {
}
