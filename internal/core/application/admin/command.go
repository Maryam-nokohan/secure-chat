package admin

import (
	"context"
	"fmt"

	"github.com/gofrs/uuid"
	"github.com/maryam-nokohan/secure-chat/internal/core/domain/user"
	"github.com/maryam-nokohan/secure-chat/internal/core/ports"
)
type Command interface {
	Execute(ctx context.Context) error
	Undo(ctx context.Context) error
	Description() string
}
type DeleteUserCommand struct {
	Repo   ports.UserRepository
	UserID uuid.UUID

	snapshot user.User
}

func (c *DeleteUserCommand) Execute(ctx context.Context) error {
	u, err := c.Repo.FindUserByID(ctx, c.UserID)
	if err != nil {
		return fmt.Errorf("user not found: %w", err)
	}
	c.snapshot = *u
	return c.Repo.DeleteUser(ctx, *u)
}

func (c *DeleteUserCommand) Undo(ctx context.Context) error {
	return c.Repo.RestoreUser(ctx, c.UserID)
}

func (c *DeleteUserCommand) Description() string {
	return fmt.Sprintf("delete user %q", c.snapshot.Username)
}

type EditUserCommand struct {
	Repo        ports.UserRepository
	UserID      uuid.UUID
	NewUsername string
	NewBio      string

	before user.User
}

func (c *EditUserCommand) Execute(ctx context.Context) error {
	u, err := c.Repo.FindUserByID(ctx, c.UserID)
	if err != nil {
		return fmt.Errorf("user not found: %w", err)
	}
	c.before = *u

	updated := *u
	if c.NewUsername != "" {
		updated.Username = c.NewUsername
	}
	updated.Bio = c.NewBio
	return c.Repo.EditUser(ctx, updated)
}

func (c *EditUserCommand) Undo(ctx context.Context) error {
	return c.Repo.EditUser(ctx, c.before)
}

func (c *EditUserCommand) Description() string {
	return fmt.Sprintf("edit user %q", c.before.Username)
}

type SetRoleCommand struct {
	Repo    ports.UserRepository
	UserID  uuid.UUID
	NewRole string

	before user.User
}

func (c *SetRoleCommand) Execute(ctx context.Context) error {
	u, err := c.Repo.FindUserByID(ctx, c.UserID)
	if err != nil {
		return fmt.Errorf("user not found: %w", err)
	}
	c.before = *u
	updated := *u
	updated.Role = c.NewRole
	return c.Repo.EditUser(ctx, updated)
}

func (c *SetRoleCommand) Undo(ctx context.Context) error {
	return c.Repo.EditUser(ctx, c.before)
}

func (c *SetRoleCommand) Description() string {
	return fmt.Sprintf("set role of %q to %q (was %q)", c.before.Username, c.NewRole, c.before.Role)
}