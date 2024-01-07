package repository

import (
	"go-rest-api/model"

	"gorm.io/gorm"
)

// Go では Interface はメソッドの一覧
/*
インターフェース (interface)
インターフェースは、特定の「約束」みたいなものです。Go言語におけるインターフェースは、ある特定のメソッドを持っているという「約束」を表します。
例えば、IUserRepositoryインターフェースは、「GetUserByEmailとCreateUserという2つのメソッドを持っている」という約束です。
*/
type IUserRepository interface {
	GetUserByEmail(user *model.User, email string) error
	CreateUser(user *model.User) error
}

// 構造体は、いくつかの異なるデータ（変数）をひとまとめにしたもの
type userRepository struct {
	db *gorm.DB
}

/*
コンストラクターは、新しいオブジェクト（構造体のインスタンス）を作るための特別な関数です。
NewUserRepository関数は、新しいuserRepositoryオブジェクトを作り、それを返します。
これにより、データベースへの接続を持つuserRepositoryを使って作業ができるようになります。
*/
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
