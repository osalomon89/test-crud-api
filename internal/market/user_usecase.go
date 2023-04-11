package market

import (
	"context"
	"fmt"

	"github.com/osalomon89/test-crud-api/internal/market/item"
	"github.com/osalomon89/test-crud-api/internal/market/user"
	"github.com/osalomon89/test-crud-api/internal/market/user/token"
	"github.com/osalomon89/test-crud-api/internal/platform/crypto"
)

type UserUseCase interface {
	Login(ctx context.Context, email, password string) (string, error)
	Register(ctx context.Context, email, password string) error
	//MarkItemAsFavorite(ctx context.Context, itemID, userID uint) error
	//GetItemsByUser(ctx context.Context, userID uint) (*user.Items, error)
}

type userUseCase struct {
	userRepository user.Repository
	itemRepository item.Repository
	tokenService   token.Service
}

func NewUserUsecase(userRepository user.Repository,
	itemRepository item.Repository, tokenService token.Service) UserUseCase {
	return &userUseCase{
		userRepository: userRepository,
		tokenService:   tokenService,
		itemRepository: itemRepository,
	}
}

func (svc *userUseCase) Login(ctx context.Context,
	email, password string) (string, error) {
	userData, err := svc.userRepository.GetByEmail(ctx, email)
	if err != nil {
		return "", fmt.Errorf("error getting user: %w", err)
	}

	if !crypto.ComparePasswords(password, userData.Password) {
		return "", user.UserError{
			Message: "invalid credentials",
		}
	}

	token, err := svc.tokenService.Get(ctx, userData.ID, userData.Email)
	if err != nil {
		return "", fmt.Errorf("error getting token: %w", err)
	}

	return token, nil
}

func (svc *userUseCase) Register(ctx context.Context, email, password string) error {
	passwordHash, err := crypto.HashPassword(password)
	if err != nil {
		return fmt.Errorf("error hashing password: %w", err)
	}

	user := user.User{
		Email:    email,
		Password: passwordHash,
	}

	if err := svc.userRepository.Save(ctx, &user); err != nil {
		return fmt.Errorf("error saving user: %w", err)
	}

	return nil
}

func (svc *userUseCase) MarkItemAsFavorite(ctx context.Context, itemID, userID uint) error {
	_, err := svc.itemRepository.GetItemByID(ctx, itemID)
	if err != nil {
		return fmt.Errorf("error getting property: %w", err)
	}

	/*err = svc.userRepository.SaveUserItem(ctx, itemID, userID)
	if err != nil {
		return fmt.Errorf("error saving property: %w", err)
	}*/

	return nil
}

/*func (svc *userUseCase) GetItemsByUser(ctx context.Context, userID uint) (*user.Items, error) {
	items, err := svc.userRepository.GetItemsByUser(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("error getting property: %w", err)
	}

	return items, nil
}*/
