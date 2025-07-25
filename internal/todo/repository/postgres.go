package repository

import (
	"context"
	"database/sql"
	"time"
)

type Todo struct {
	ID        int64     `json:"id"`
	Title     string    `json:"title"`
	Completed bool      `json:"completed"`
	CreatedAt time.Time `json:"created_at"`
}

type Repository interface {
	Create(ctx context.Context, title string) (int64, error)
	List(ctx context.Context) ([]Todo, error)
	Toggle(ctx context.Context, id int64, completed bool) error
	Delete(ctx context.Context, id int64) error
}

type PostgresRepo struct {
	db *sql.DB
}

func NewPostgres(db *sql.DB) *PostgresRepo {
	return &PostgresRepo{db: db}
}

func (r *PostgresRepo) Create(ctx context.Context, title string) (int64, error) {
	var id int64
	err := r.db.QueryRowContext(ctx,
		`INSERT INTO todos (title) VALUES ($1) RETURNING id`, title).Scan(&id)
	return id, err
}

func (r *PostgresRepo) List(ctx context.Context) ([]Todo, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, title, completed, created_at FROM todos ORDER BY id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var todos []Todo
	for rows.Next() {
		var t Todo
		if err := rows.Scan(&t.ID, &t.Title, &t.Completed, &t.CreatedAt); err != nil {
			return nil, err
		}
		todos = append(todos, t)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return todos, nil
}

func (r *PostgresRepo) Toggle(ctx context.Context, id int64, completed bool) error {
	result, err := r.db.ExecContext(ctx, 
		`UPDATE todos SET completed = $1 WHERE id = $2`,
		completed, id)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (r *PostgresRepo) Delete(ctx context.Context, id int64) error {
	result, err := r.db.ExecContext(ctx, 
		`DELETE FROM todos WHERE id = $1`,
 		id)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return sql.ErrNoRows
	}

	return nil
}