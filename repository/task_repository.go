package repository

import (
	"fmt"
	"go-rest-api/model"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ITaskRepository interface {
	GetAllTasks(tasks *[]model.Task, userId uint) error
	GetTaskById(task *model.Task, userId uint, taskId uint) error
	CreateTask(task *model.Task) error
	UpdateTask(task *model.Task, userId uint, taskId uint) error
	DeleteTask(userId uint, taskId uint) error
}

type taskRepository struct {
	db *gorm.DB
}

// コンストラクター
func NewTaskRepository(db *gorm.DB) ITaskRepository {
	return &taskRepository{db}
}

func (tr *taskRepository) GetAllTasks(tasks *[]model.Task, userId uint) error {
	/*
		taskの一覧の中で、user_idが引数のuserIdと一致するものを取得
		Order("created_at") で作成日時の昇順で取得
		Findメソッドは、検索結果をtasks変数にマッピングします。
		この変数が指すスライスには、取得したタスクのデータが順に格納されます。
	*/
	if err := tr.db.Joins("User").Where("user_id=?", userId).Order("created_at").Find(tasks).Error; err != nil {
		return err
	}
	return nil
}

func (tr *taskRepository) GetTaskById(task *model.Task, userId uint, taskId uint) error {
	/*
		taskの一覧の中で、user_idが引数のuserIdと一致するものを取得
		task_idが引数のtaskIdと一致するものを取得
		First(task, taskId)メソッドは、検索結果をtask変数にマッピング（書き込み）します。
		ここでtaskはmodel.Task型のポインタであるため、データベースから取得したタスク情報がtask変数の指すメモリ領域に直接書き込まれます。
	*/
	if err := tr.db.Joins("User").Where("user_id=?", userId).First(task, taskId).Error; err != nil {
		return err
	}
	return nil
}

func (tr *taskRepository) CreateTask(task *model.Task) error {
	if err := tr.db.Create(task).Error; err != nil {
		return err
	}
	return nil
}

func (tr *taskRepository) UpdateTask(task *model.Task, userId uint, taskId uint) error {
	/*
		tr.db.Model(task)は、操作対象となるモデル（この場合はtask）を指定します。
		Clauses(clause.Returning{})は、更新後のレコードを取得するためのオプションです。
		Where("id=? AND user_id=?", taskId, userId)は、更新対象のタスクを指定します。ここでは、タスクのIDがtaskIdで、ユーザーIDがuserIdであるものを条件としています。
		Update("title", task.Title)は、タスクのタイトルを更新します。第一引数には更新対象のフィールド名、第二引数には更新後の値を指定します。
	*/
	result := tr.db.Model(task).Clauses(clause.Returning{}).Where("id=? AND user_id=?", taskId, userId).Update("title", task.Title)
	if result.Error != nil {
		return result.Error
	}
	// 更新しようとしたオブジェクトが存在しなくてもエラーにならないため
	// RowsAffected フィールドで更新したレコード数を確認
	if result.RowsAffected < 1 {
		return fmt.Errorf("object dose not exist")
	}
	return nil
}

func (tr *taskRepository) DeleteTask(userId uint, taskId uint) error {
	// Delete(&model.Task{})は、タスクを削除します。引数には、削除対象のモデルを指定します。
	result := tr.db.Where("id=? AND user_id=?", taskId, userId).Delete(&model.Task{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected < 1 {
		return fmt.Errorf("object dose not exist")
	}
	return nil
}
