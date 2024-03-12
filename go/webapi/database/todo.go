package database

import (
	"context"
	"errors"

	"github.com/doug-martin/goqu/v9"
	"github.com/jackc/pgx/v5"
	"github.com/nanoteck137/{{ .ProjectName }}/types"
	"github.com/nanoteck137/{{ .ProjectName }}/utils"
)

type Todo struct {
	Id string
	Content string
	Done bool
}

func (db *Database) GetAllTodos(ctx context.Context) ([]Todo, error) {
	ds := dialect.From("todos").Select("id", "content", "done")

	rows, err := db.Query(ctx, ds)
	if err != nil {
		return nil, err
	}

	items, err := pgx.CollectRows[Todo](rows, func(row pgx.CollectableRow) (Todo, error) {
		var item Todo
		err := row.Scan(&item.Id, &item.Content, &item.Done)
		return item, err
	})
	if err != nil {
		return nil, err
	}

	return items, nil
}

func (db *Database) GetTodoById(ctx context.Context, id string) (Todo, error) {
	ds := dialect.From("todos").
		Select("id", "content", "done").
		Where(goqu.I("id").Eq(id)).
		Prepared(true)

	row, err := db.QueryRow(ctx, ds)
	if err != nil {
		return Todo{}, err
	}

	var item Todo
	err = row.Scan(&item.Id, &item.Content, &item.Done)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return Todo{}, types.ErrNoTodo
		}
		return Todo{}, err
	}

	return item, nil
}

func (db *Database) CreateTodo(ctx context.Context, content string) (Todo, error) {
	ds := dialect.Insert("todos").Rows(goqu.Record{
		"id": utils.CreateId(),
		"content": content,
		"done": false,
	}).Returning("id", "content", "done")

	row, err := db.QueryRow(ctx, ds)
	if err != nil {
		return Todo{}, err
	}

	var item Todo
	err = row.Scan(&item.Id, &item.Content, &item.Done)
	if err != nil {
		return Todo{}, err
	}

	return item, nil
}

type TodoChanges struct {
	Content types.Change[string]
	Done    types.Change[bool]
}

func (db *Database) UpdateTodo(ctx context.Context, id string, changes TodoChanges) (Todo, error) {
	record := goqu.Record{}

	if changes.Content.Changed {
		record["content"] = changes.Content.Value
	}

	if changes.Done.Changed {
		record["done"] = changes.Done.Value
	}

	if len(record) == 0 {
		return Todo{}, types.ErrNoChanges
	}

	ds := dialect.Update("todos").
		Set(record).
		Where(goqu.I("id").Eq(id)).
		Returning("id", "content", "done").
		Prepared(true)

	row, err := db.QueryRow(ctx, ds)
	if err != nil {
		return Todo{}, err
	}

	var item Todo
	err = row.Scan(&item.Id, &item.Content, &item.Done)
	if err != nil {
		return Todo{}, err
	}

	return item, nil
}
