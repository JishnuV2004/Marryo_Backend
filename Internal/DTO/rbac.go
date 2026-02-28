package dto

type CreateRoleRequest struct {
	Name string `json:"name"`
}

type UpdateRolePermissionsRequest struct {
	RoleID        uint   `json:"role_id"`
	PermissionIDs []uint `json:"permission_ids"`
}