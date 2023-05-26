package grpc

import (
	"context"
	"errors"

	"github.com/chopper-c2-framework/c2-chopper/grpc/proto"
	"github.com/google/uuid"

	"github.com/chopper-c2-framework/c2-chopper/core/domain/entity"
	services "github.com/chopper-c2-framework/c2-chopper/core/services"
)

type TaskService struct {
	proto.UnimplementedTaskServiceServer
	TaskService  services.ITaskService
	AgentService services.IAgentService
}

func (s *TaskService) GetTask(ctx context.Context, in *proto.GetTaskRequest) (*proto.GetTaskResponse, error) {
	if len(in.GetTaskId()) == 0 {
		return &proto.GetTaskResponse{}, errors.New("Task id required")
	}

	task, err := s.TaskService.FindTaskOrError(in.GetTaskId())
	if err != nil {
		return &proto.GetTaskResponse{}, errors.New("Task not found")
	}
	return &proto.GetTaskResponse{
		Task: ConvertTaskToProto(task),
	}, nil
}

func (s *TaskService) DeleteTask(ctx context.Context, in *proto.DeleteTaskRequest) (*proto.DeleteTaskResponse, error) {
	if len(in.GetTaskId()) == 0 {
		return &proto.DeleteTaskResponse{}, errors.New("Task id required")
	}

	task, err := s.TaskService.FindTaskOrError(in.GetTaskId())
	if err != nil {
		return &proto.DeleteTaskResponse{}, errors.New("Task not found")
	}

	s.TaskService.DeleteTask(task)
	return &proto.DeleteTaskResponse{}, nil
}

func (s *TaskService) CreateTask(ctx context.Context, in *proto.CreateTaskRequest) (*proto.CreateTaskResponse, error) {
	if len(in.GetAgentId()) == 0 {
		return &proto.CreateTaskResponse{}, errors.New("Agent id required")
	}
	agentId, err := uuid.Parse(in.GetAgentId())
	if err != nil {
		return &proto.CreateTaskResponse{}, errors.New("Invalid agent id")
	}

	taskProto := in.GetTask()
	err = ValidateTaskProto(taskProto)
	if err != nil {
		return &proto.CreateTaskResponse{}, err
	}

	// TODO: Add user id
	var task = entity.TaskModel{
		Name:    taskProto.GetName(),
		Args:    taskProto.GetArgs(),
		Type:    entity.TaskType(taskProto.GetType().String()),
		AgentId: agentId,
		// CreatorId: ,
	}

	err = s.TaskService.CreateTask(&task)
	if err != nil {
		return &proto.CreateTaskResponse{}, err
	}

	return &proto.CreateTaskResponse{}, nil
}

func (s *TaskService) GetAgentTasks(ctx context.Context, in *proto.GetAgentTasksRequest) (*proto.GetAgentTasksResponse, error) {
	agentId := in.GetAgentId()
	if len(agentId) == 0 {
		return &proto.GetAgentTasksResponse{}, errors.New("Agent id required")
	}

	agent, err := s.AgentService.FindAgentOrError(agentId)
	if err != nil {
		return &proto.GetAgentTasksResponse{}, errors.New("Agent not found")
	}

	tasks, err := s.TaskService.FindUnexecutedTasksForAgent(agentId)
	if err != nil {
		return &proto.GetAgentTasksResponse{}, err
	}

	protoList := make([]*proto.Task, len(tasks))
	for i, task := range tasks {
		protoList[i] = ConvertTaskToProto(task)
	}

	return &proto.GetAgentTasksResponse{
		Tasks:     protoList,
		SleepTime: agent.SleepTime,
	}, nil
}

func (s *TaskService) CreateTaskResult(ctx context.Context, in *proto.CreateTaskResultRequest) (*proto.CreateTaskResultResponse, error) {
	taskResProto := in.GetTaskResult()
	err := ValidateTaskResultProto(taskResProto)
	if err != nil {
		return &proto.CreateTaskResultResponse{}, err
	}

	taskUUID := uuid.MustParse(taskResProto.GetTaskId())
	taskResult := &entity.TaskResultModel{
		Status: taskResProto.GetStatus(),
		Output: taskResProto.GetOutput(),
		TaskID: taskUUID,
	}

	err = s.TaskService.CreateTaskResult(taskResult)
	if err != nil {
		return &proto.CreateTaskResultResponse{}, err
	}

	return &proto.CreateTaskResultResponse{}, nil
}

func (s *TaskService) GetTaskResults(ctx context.Context, in *proto.GetTaskResultsRequest) (*proto.GetTaskResultsResponse, error) {
	if len(in.GetTaskId()) == 0 {
		return &proto.GetTaskResultsResponse{}, errors.New("Task id required")
	}

	taskResults, err := s.TaskService.FindTaskResults(in.GetTaskId())
	if err != nil {
		return &proto.GetTaskResultsResponse{}, err
	}

	protoList := make([]*proto.TaskResult, len(taskResults))
	for i, taskRes := range taskResults {
		protoList[i] = ConvertTaskResultToProto(taskRes)
	}

	return &proto.GetTaskResultsResponse{
		Results: protoList,
	}, nil
}

func (s *TaskService) SetTaskResultsSeen(ctx context.Context, in *proto.SetTaskResultsSeenRequest) (*proto.SetTaskResultsSeenResponse, error) {
	resultIds := in.GetResultIds()
	if resultIds == nil {
		return &proto.SetTaskResultsSeenResponse{}, errors.New("At least 1 id is required in result_ids")
	}

	for _, id := range resultIds {
		err := s.TaskService.MarkTaskResultSeen(id)
		if err != nil {
			return &proto.SetTaskResultsSeenResponse{}, err
		}
	}
	return &proto.SetTaskResultsSeenResponse{}, nil
}
