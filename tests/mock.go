type MockTaskRepository struct {
    mock.Mock
}

func (m *MockTaskRepository) Create(task *Task) error {
    args := m.Called(task)
    return args.Error(0)
}
