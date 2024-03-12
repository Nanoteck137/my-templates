package handlers

import (
	"github.com/labstack/echo/v4"
	"github.com/nanoteck137/{{ .ProjectName }}/database"
	"github.com/nanoteck137/{{ .ProjectName }}/types"
)

func (h *Handlers) GetTodos(c echo.Context) error {
	todos, err := h.db.GetAllTodos(c.Request().Context())
	if err != nil {
		return err
	}

	res := types.GetTodos{
		Todos: make([]types.Todo, len(todos)),
	}

	for i, todo := range todos {
		res.Todos[i] = types.Todo{
			Id:      todo.Id,
			Content: todo.Content,
			Done:    todo.Done,
		}
	}

	return c.JSON(200, types.NewApiSuccessResponse(res))
}

func (h *Handlers) PostTodos(c echo.Context) error {
	type Body struct {
		Content string `json:"content"`
	}

	var body Body
	err := c.Bind(&body)
	if err != nil {
		return err
	}

	todo, err := h.db.CreateTodo(c.Request().Context(), body.Content)

	return c.JSON(200, types.NewApiSuccessResponse(types.PostTodos{
		Id:      todo.Id,
		Content: todo.Content,
		Done:    todo.Done,
	}))
}

func (h *Handlers) GetTodoById(c echo.Context) error {
	id := c.Param("id")

	todo, err := h.db.GetTodoById(c.Request().Context(), id)
	if err != nil {
		return err
	}

	return c.JSON(200, types.GetTodoById{
		Id:      todo.Id,
		Content: todo.Content,
		Done:    todo.Done,
	})
}

func (h *Handlers) PutTodoById(c echo.Context) error {
	id := c.Param("id")

	type Body struct {
		Content *string `json:"content"`
		Done    *bool   `json:"done"`
	}

	var body Body
	err := c.Bind(&body)
	if err != nil {
		return err
	}

	changes := database.TodoChanges{}
	if body.Content != nil {
		changes.Content.Value = *body.Content
		changes.Content.Changed = true
	}
	if body.Done != nil {
		changes.Done.Value = *body.Done
		changes.Done.Changed = true
	}

	todo, err := h.db.UpdateTodo(c.Request().Context(), id, changes)
	if err != nil {
		return err
	}

	return c.JSON(200, types.NewApiSuccessResponse(types.PutTodoById{
		Id:      todo.Id,
		Content: todo.Content,
		Done:    todo.Done,
	}))
}

func (h *Handlers) InstallTodoHandlers(group *echo.Group) {
	group.GET("/todos", h.GetTodos)
	group.POST("/todos", h.PostTodos)
	group.GET("/todos/:id", h.GetTodoById)
	group.PUT("/todos/:id", h.PutTodoById)
}
