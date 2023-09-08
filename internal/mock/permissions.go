package mock

import (
	"greenlight.jattueyi.com/internal/data"
)

var mockPermission data.Permissions = []string{"fooo"}

type PermissionModel struct{}

func (m PermissionModel) AddForUser(userID int64, codes ...string) error {
	return nil
}

func (m PermissionModel) GetAllForUser(userID int64) (data.Permissions, error) {
	return mockPermission, nil
}
