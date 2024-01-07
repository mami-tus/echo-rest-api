package usecase

import (
	"go-rest-api/model"
	"go-rest-api/repository"
	"go-rest-api/validator"
)

type ITaskUsecase interface {
	GetAllTasks(userId uint) ([]model.TaskResponse, error)
	GetTaskById(userId uint, taskId uint) (model.TaskResponse, error)
	CreateTask(task model.Task) (model.TaskResponse, error)
	UpdateTask(task model.Task, userId uint, taskId uint) (model.TaskResponse, error)
	DeleteTask(userId uint, taskId uint) error
}

type taskUsecase struct {
	tr repository.ITaskRepository
	tv validator.ITaskValidator
}

// コンストラクター
func NewTaskUsecase(tr repository.ITaskRepository, tv validator.ITaskValidator) ITaskUsecase {
	return &taskUsecase{tr, tv}
}

func (tu *taskUsecase) GetAllTasks(userId uint) ([]model.TaskResponse, error) {
	// ユーザーのタスクを取得
	tasks := []model.Task{}
	if err := tu.tr.GetAllTasks(&tasks, userId); err != nil {
		// ライスのゼロ値（初期値）はnilのため、エラーが発生した場合、nil（スライスのゼロ値）とエラー情報を返す
		return nil, err
	}
	// レスポンス用のタスクを作成
	resTasks := []model.TaskResponse{}
	for _, v := range tasks {
		t := model.TaskResponse{
			ID:        v.ID,
			Title:     v.Title,
			CreatedAt: v.CreatedAt,
			UpdatedAt: v.UpdatedAt,
		}
		resTasks = append(resTasks, t)
	}
	return resTasks, nil
}

func (tu *taskUsecase) GetTaskById(userId uint, taskId uint) (model.TaskResponse, error) {
	// ユーザーのタスクを取得
	task := model.Task{}
	if err := tu.tr.GetTaskById(&task, userId, taskId); err != nil {
		// エラーが発生した場合、model.TaskResponse{}（構造体のゼロ値）とエラー情報を返す
		return model.TaskResponse{}, err
	}
	// レスポンス用のタスクを作成
	resTask := model.TaskResponse{
		ID:        task.ID,
		Title:     task.Title,
		CreatedAt: task.CreatedAt,
		UpdatedAt: task.UpdatedAt,
	}
	return resTask, nil
}

func (tu *taskUsecase) CreateTask(task model.Task) (model.TaskResponse, error) {
	/*
		メソッドの引数 - task model.Task:
		taskはmodel.Task型の変数で、新しいタスクの情報を保持しています。この変数は値渡し（コピー）されてCreateTaskメソッドに渡されます。

		&task - ポインタの使用:
		&taskは、task変数のアドレス、つまりポインタを取得する操作です。これにより、task変数自体（そのメモリ上の場所）への参照が得られます。
		tu.tr.CreateTask(&task)では、taskのポインタをCreateTaskメソッドに渡しています。これにより、CreateTaskメソッド内でtask変数を直接変更できるようになります。

		データベースへの書き込み:
		CreateTaskメソッド内で、データベースにtaskオブジェクトを保存します。ここでポインタを使用することで、データベース操作による変更（例えば、新しいIDの割り当てや作成日時の記録など）が元のtaskオブジェクトに直接反映されます。

		戻り値の作成:
		データベース操作が成功した後、新しいmodel.TaskResponseオブジェクトが作成され、データベースから得られたデータ（taskの内容）で初期化されています。
		これは、クライアントに返すためのレスポンスとして構成されます。
	*/
	if err := tu.tv.TaskValidate(task); err != nil {
		return model.TaskResponse{}, err
	}
	// タスクを作成
	if err := tu.tr.CreateTask(&task); err != nil {
		// エラーが発生した場合、model.TaskResponse{}（構造体のゼロ値）とエラー情報を返す
		return model.TaskResponse{}, err
	}
	// レスポンス用のタスクを作成
	resTask := model.TaskResponse{
		ID:        task.ID,
		Title:     task.Title,
		CreatedAt: task.CreatedAt,
		UpdatedAt: task.UpdatedAt,
	}
	return resTask, nil
}

func (tu *taskUsecase) UpdateTask(task model.Task, userId uint, taskId uint) (model.TaskResponse, error) {
	if err := tu.tv.TaskValidate(task); err != nil {
		return model.TaskResponse{}, err
	}
	// タスクを更新
	if err := tu.tr.UpdateTask(&task, userId, taskId); err != nil {
		return model.TaskResponse{}, err
	}
	// レスポンス用のタスクを作成
	resTask := model.TaskResponse{
		ID:        task.ID,
		Title:     task.Title,
		CreatedAt: task.CreatedAt,
		UpdatedAt: task.UpdatedAt,
	}
	return resTask, nil
}

func (tu *taskUsecase) DeleteTask(userId uint, taskId uint) error {
	// タスクを削除
	if err := tu.tr.DeleteTask(userId, taskId); err != nil {
		return err
	}
	return nil
}
