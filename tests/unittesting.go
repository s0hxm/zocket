func TestCreateTask(t *testing.T) {
    mockRepo := new(MockTaskRepository)
    service := NewTaskService(mockRepo)

    task := &Task{Title: "Test Task", Description: "Test Description"}
    mockRepo.On("Create", task).Return(nil)

    err := service.CreateTask(task)
    assert.NoError(t, err)
    mockRepo.AssertExpectations(t)
}
