package usecase

import (
	"go-rest-api/model"
	"go-rest-api/repository"
	"go-rest-api/validator"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

type IUserUsecase interface {
	SignUp(user model.User) (model.UserResponse, error)
	LogIn(user model.User) (string, error)
}

type userUsecase struct {
	ur repository.IUserRepository
	uv validator.IUserValidator
}

// コンストラクター
func NewUserUsecase(ur repository.IUserRepository, uv validator.IUserValidator) IUserUsecase {
	return &userUsecase{ur, uv}
}

// userUsecase型 の SignUpメソッド
func (uu *userUsecase) SignUp(user model.User) (model.UserResponse, error) {
	// バリデーション
	if err := uu.uv.UserValidate(user); err != nil {
		return model.UserResponse{}, err
	}
	// パスワードをハッシュ化
	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), 10)
	if err != nil {
		return model.UserResponse{}, err
	}
	// ハッシュ化したパスワードを user.Password に代入して newUser を作成
	newUser := model.User{Email: user.Email, Password: string(hash)}
	// ユーザーを作成
	if err := uu.ur.CreateUser(&newUser); err != nil {
		return model.UserResponse{}, err
	}
	// レスポンス用のユーザーを作成
	resUser := model.UserResponse{
		ID:    newUser.ID,
		Email: newUser.Email,
	}
	return resUser, nil
}

func (uu *userUsecase) LogIn(user model.User) (string, error) {
	// バリデーション
	if err := uu.uv.UserValidate(user); err != nil {
		return "", err
	}
	storedUser := model.User{}
	// ユーザーを取得
	if err := uu.ur.GetUserByEmail(&storedUser, user.Email); err != nil {
		return "", err
	}
	// 保存されてるパスワードと入力されたパスワードを比較
	err := bcrypt.CompareHashAndPassword([]byte(storedUser.Password), []byte(user.Password))
	if err != nil {
		return "", err
	}
	// NewWithClaims の第一引数にアルゴリズム、第二引数にペイロードの設定
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": storedUser.ID,
		"exp":     time.Now().Add(time.Hour * 12).Unix(), // 有効期限12時間
	})
	// JWTを作成
	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET")))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}
