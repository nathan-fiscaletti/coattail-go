package permission

import (
	"context"
	"errors"

	"github.com/nathan-fiscaletti/coattail-go/internal/keys"
)

var (
	ErrPermissionsNotFound = errors.New("permissions not found in context")
)

type Permission int32

// permissions are stored in a 16 bit bitmask
const (
	ReadActions Permission = 1 << iota
	ReadReceivers
	ReadPeers

	All = ReadActions | ReadReceivers | ReadPeers
)

type Permissions interface {
	Has(permission Permission) bool
	HasOneOf(permission ...Permission) bool
	HasAnyOf(permission ...Permission) bool
	HasAllOf(permission ...Permission) bool
	Permitted() int32
}

func PermissionsFromContext(ctx context.Context) (Permissions, error) {
	permissions, ok := ctx.Value(keys.PermissionsKey).(Permissions)
	if !ok {
		return nil, ErrPermissionsNotFound
	}

	return permissions, nil
}

func ContextWithPermissions(ctx context.Context, permissions Permissions) context.Context {
	return context.WithValue(ctx, keys.PermissionsKey, permissions)
}

func GetPermissions(permitted int32) Permissions {
	return &permissions{
		permitted: permitted,
	}
}

func PermissionMask(permissions ...Permission) int32 {
	permitted := int32(0)
	for _, p := range permissions {
		permitted |= int32(p)
	}
	return permitted
}

type permissions struct {
	permitted int32
}

func (s *permissions) Has(permission Permission) bool {
	return s.permitted&int32(permission) != 0
}

func (s *permissions) HasOneOf(permission ...Permission) bool {
	for _, p := range permission {
		if s.Has(p) {
			return true
		}
	}
	return false
}

func (s *permissions) HasAnyOf(permission ...Permission) bool {
	for _, p := range permission {
		if s.Has(p) {
			return true
		}
	}
	return false
}

func (s *permissions) HasAllOf(permission ...Permission) bool {
	for _, p := range permission {
		if !s.Has(p) {
			return false
		}
	}
	return true
}

func (s *permissions) Permitted() int32 {
	return s.permitted
}
