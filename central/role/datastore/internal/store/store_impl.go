package store

import (
	"github.com/stackrox/rox/generated/storage"
	"github.com/stackrox/rox/pkg/bolthelper/crud/proto"
)

type storeImpl struct {
	roleCrud proto.MessageCrud
}

// AddRole adds a role to the store.
// Returns an error if the role already exists
func (s *storeImpl) AddRole(role *storage.Role) error {
	return s.roleCrud.Create(role)
}

// UpdateRole udpates a role to the store.
// Returns an error if the role does not already exist, or if the role is a pre-loaded role.
// Pre-loaded roles cannot be updated./
func (s *storeImpl) UpdateRole(role *storage.Role) error {
	_, _, err := s.roleCrud.Update(role)
	return err
}

// RemoveRole removes a role from the store.
// Pre-loaded roles cannot be removed.
func (s *storeImpl) RemoveRole(name string) error {
	_, _, err := s.roleCrud.Delete(name)
	return err
}

// GetRole returns a role from the store by name.
// Returns nil without an error if the requested role does not exist.
func (s *storeImpl) GetRole(name string) (*storage.Role, error) {
	msg, err := s.roleCrud.Read(name)
	if msg == nil {
		return nil, err
	}
	return msg.(*storage.Role), err
}

// GetAllRoles returns all of the roles in the store.
// Returns nil without an error if no roles exist in the store (default roles cannot be deleted, so never)
func (s *storeImpl) GetAllRoles() ([]*storage.Role, error) {
	msgs, err := s.roleCrud.ReadAll()
	if err != nil {
		return nil, err
	}
	if len(msgs) == 0 {
		return nil, err
	}
	// Cast to a list of roles.
	Roles := make([]*storage.Role, 0, len(msgs))
	for _, msg := range msgs {
		Roles = append(Roles, msg.(*storage.Role))
	}
	return Roles, nil
}