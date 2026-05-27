package postgres

import (
	"context"

	"github.com/gofrs/uuid"
	"github.com/maryam-nokohan/secure-chat/internal/core/domain/chat"
	"github.com/maryam-nokohan/secure-chat/internal/core/ports"
	"gorm.io/gorm"
)

type ChatRoomDB struct {
	ID        uuid.UUID
	Name      string
	CreatorID uuid.UUID
	Users     []uuid.UUID
}

type ChatRepository struct {
	db gorm.DB
}

func NewChatRepository(db *gorm.DB) ports.ChatRepositoryI {
	return &ChatRepository{
		db: *db,
	}
}

func (r *ChatRepository) CreateRoom(ctx context.Context, room *chat.Room) error {

	dbModel := ChatRoomDB{
		ID:        room.ID,
		Name:      room.Name,
		CreatorID: room.CreatorID,
		Users:     room.Users,
	}
	return r.db.WithContext(ctx).Create(&dbModel).Error

}
func (r *ChatRepository) FindRoomByID(ctx context.Context, id uuid.UUID) (*chat.Room, error) {

	var dbModel ChatRoomDB
	result := r.db.WithContext(ctx).Where("id = ?", id).First(&dbModel)
	if result.Error != nil {
		return nil, result.Error
	}
	return &chat.Room{
		ID:        dbModel.ID,
		Name:      dbModel.Name,
		CreatorID: dbModel.CreatorID,
		Users:     dbModel.Users,
	}, nil
}
func (r *ChatRepository) FindRoomByName(ctx context.Context, name string) (*chat.Room, error) {

		var dbModel ChatRoomDB
	result := r.db.WithContext(ctx).Where("name = ?", name).First(&dbModel)
	if result.Error != nil {
		return nil, result.Error
	}
	return &chat.Room{
		ID:        dbModel.ID,
		Name:      dbModel.Name,
		CreatorID: dbModel.CreatorID,
		Users:     dbModel.Users,
	}, nil

}
func (r *ChatRepository) ListRooms(ctx context.Context) ([]*chat.Room, error) {

	var dbModels []ChatRoomDB
	err := r.db.WithContext(ctx).Find(&dbModels).Error
	if err != nil {
		return nil , err
	}
	rooms := make([]*chat.Room ,len(dbModels))
	for i , model := range dbModels {
		rooms[i] = &chat.Room{
			ID: model.ID,
			Name: model.Name,
			Users: model.Users,
			CreatorID: model.CreatorID,
		}
	}
	return  rooms , nil
}
func (r *ChatRepository) DeleteRoom(ctx context.Context, id uuid.UUID) error {

	return r.db.WithContext(ctx).Where("id = ?" , id).Delete(&ChatRoomDB{}).Error
}
func (r *ChatRepository) AddUserToRoom(ctx context.Context, roomID, userID uuid.UUID) error {
	// resuls := r.db.WithContext(ctx).Where("id = ?" , roomID)
	// if resuls.Error != nil {
	// 	return  resuls.Error
	// }

	return  nil
}
func (r *ChatRepository) RemoveUserFromRoom(ctx context.Context, roomID, userID uuid.UUID) error {
	return  nil
}
