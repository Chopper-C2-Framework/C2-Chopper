package services

import (
	"fmt"

	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"

	orm "github.com/chopper-c2-framework/c2-chopper/core/domain"
	entity "github.com/chopper-c2-framework/c2-chopper/core/domain/entity"
)

type TaskService struct {
	ORMConnection *orm.ORMConnection
	repo          entity.TransactionRepository
}

func NewTaskService(db *orm.ORMConnection) *TaskService {
	logger := log.New()

	repo := entity.NewGormRepository(db.Db, logger)
	return &TaskService{
		repo: repo,
	}
}

func (s *TaskService) CreateTask(task *entity.TaskModel) error {
	err := s.repo.Create(task)
	if err != nil {
		log.Debugf("[-] failed to create task")
		return err
	}

	return nil
}

func (s *TaskService) DeleteTask(task *entity.TaskModel) error {
	err := s.repo.Delete(task)
	if err != nil {
		log.Debugf("[-] failed to delete task")
		return err
	}

	return nil
}

func (s *TaskService) FindTaskOrError(taskId string) (*entity.TaskModel, error) {
	var task entity.TaskModel
	err := s.repo.GetOneByID(&task, taskId)

	if err != nil {
		return nil, err
	}
	return &task, nil
}

func (s *TaskService) FindTasksForAgent(agentId string) ([]entity.TaskModel, error) {
	var tasks []entity.TaskModel

	// TODO: fix this
	fmt.Println(agentId)

	x, err := uuid.Parse(agentId)
	if err != nil {
		log.Debugf("[-] failed to parse uuid")
		return nil, err
	}
	err = s.repo.DB().Find(&tasks).Where("agent_id = ?", x).Error

	if err != nil {
		log.Debugf("[-] failed to find task by agentid")
		return nil, err
	}
	fmt.Println(tasks)

	return tasks, nil
}
