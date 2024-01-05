package repository

import (
	"go-rest-api/model"

	"gorm.io/gorm"
)

// Go では Interface はメソッドの一覧
type IUserRepository interface {
	GetUserByEmail(user *model.User, email string) error
	CreateUser(user *model.User) error
}

type userRepository struct {
	db *gorm.DB
}

// コンストラクター
func NewUserRepository(db *gorm.DB) IUserRepository {
	return &userRepository{db}
}

func (ur *userRepository) GetUserByEmail(user *model.User, email string) error {
	if err := ur.db.Where("email=?", email).First(user).Error; err != nil {
		/*
			データベースのユーザーテーブルkら email フィールドが引数 email と一致するものを取得
			最初に見つかったレコードを user 変数に格納
			GORMのクエリメソッドは、最終的に*gorm.DB型のオブジェクトを返します。このオブジェクトには、クエリの実行結果を表すErrorフィールドが含まれています。
		*/
		return err
	}
	return nil
}

func (ur *userRepository) CreateUser(user *model.User) error {
	if err := ur.db.Create(user).Error; err != nil {
		return err
	}
	return nil
}
