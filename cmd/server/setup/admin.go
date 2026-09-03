package setup

import (
	"context"
	"fmt"
	"os"

	"github.com/gofrs/uuid"

	"github.com/maryam-nokohan/secure-chat/internal/core/domain/user"
	"github.com/maryam-nokohan/secure-chat/internal/core/ports"
	"github.com/maryam-nokohan/secure-chat/pkg"
)
func BootstrapAdmin(repo ports.UserRepository) error {
	username := os.Getenv("ADMIN_BOOTSTRAP_USERNAME")
	password := os.Getenv("ADMIN_BOOTSTRAP_PASSWORD")
	if username == "" || password == "" {
		return nil
	}
	if _, err := repo.FindUserByUsername(context.Background(), username); err == nil {
		return nil 
	}
	if err := pkg.ValidatePassword(password); err != nil {
		return fmt.Errorf("ADMIN_BOOTSTRAP_PASSWORD invalid: %w", err)
	}
	_, pub, err := pkg.GenerateRSAKeyPair(2048)
	if err != nil {
		return err
	}
	pubPEM, err := pkg.EncodePublicKey(pub)
	if err != nil {
		return err
	}
	hash, err := pkg.HashPassword(password)
	if err != nil {
		return err
	}
	id, err := uuid.NewV4()
	if err != nil {
		return err
	}
	u := user.User{ID: id, Username: username, PassHash: hash, PublicKey: pubPEM, Role: "admin"}
	if err := repo.CreateUser(context.Background(), u); err != nil {
		return err
	}
	pkg.LogInfo("Bootstrapped admin user: " + username)
	return nil
}