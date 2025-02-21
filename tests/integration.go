func TestTaskCreationFlow(t *testing.T) {
    db, err := setupTestDB()
    assert.NoError(t, err)
    defer db.Close()

    repo := NewTaskRepository(db)
    service := NewTaskService(repo)
    handler := NewTaskHandler(service)

    router := setupRouter(handler)

    task := &Task{Title: "Integration Test Task"}
    body, _ := json.Marshal(task)
    req, _ := http.NewRequest("POST", "/tasks", bytes.NewBuffer(body))
    resp := httptest.NewRecorder()

    router.ServeHTTP(resp, req)

    assert.Equal(t, http.StatusCreated, resp.Code)

    var createdTask Task
    json.Unmarshal(resp.Body.Bytes(), &createdTask)
    assert.Equal(t, task.Title, createdTask.Title)
}
