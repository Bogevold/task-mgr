package store

import (
	"fmt"

	"database/sql"

	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/bogevold/task-mgr/internal/task"
)

type PostgresStore struct {
	db *sql.DB
}

func NewPostgresStore(connStr string) (*PostgresStore, error) {
	// åpne tilkobling med sql.Open
	db, err := sql.Open("pgx", connStr)
	if err != nil {
		return nil, fmt.Errorf("Kunne ikke åpne database: %w\n", err)
	}
	// ping databasen for å verifisere tilkoblingen
	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("Kunne ikke pinge database: %w\n", err)
	}
	return &PostgresStore{db: db}, nil
}

func (p *PostgresStore) Save(inTask task.Task) (task.Task, error) {
	var newTask task.Task
	qry := "INSERT INTO tasks (title, done, created_at) VALUES ($1, $2, $3) RETURNING id, title, done, created_at"
	row := p.db.QueryRow(qry, inTask.Title, inTask.Done, inTask.CreatedAt)
	err := row.Scan(&newTask.ID, &newTask.Title, &newTask.Done, &newTask.CreatedAt)
	if err != nil {
		return task.Task{}, fmt.Errorf("Databaselagring feilet: %w\n", err)
	}
	return newTask, nil
}

func (p *PostgresStore) Get(i uint) (task.Task, error) {
	var foundTask task.Task
	row := p.db.QueryRow("select id, title, done, created_at from tasks where id = $1", i)
	err := row.Scan(&foundTask.ID, &foundTask.Title, &foundTask.Done, &foundTask.CreatedAt)
	if err != nil {
		return task.Task{}, fmt.Errorf("task med id %d finnes ikke i databasen", i)
	}
	return foundTask, nil
}

func (p *PostgresStore) GetAll() ([]task.Task, error) {
	var foundTasks []task.Task
	rows, err := p.db.Query("select id, title, done, created_at from tasks")
	if err != nil {
		return nil, fmt.Errorf("kunne ikke hente tasks: %w", err)
	}

	defer rows.Close()

	for rows.Next() {
		var t task.Task
		if err := rows.Scan(&t.ID, &t.Title, &t.Done, &t.CreatedAt); err != nil {
			return nil, fmt.Errorf("kunne ikke lese rad: %w", err)
		}
		foundTasks = append(foundTasks, t)
	}
	return foundTasks, nil
}

func (p *PostgresStore) Delete(i uint) error {
	_, err := p.db.Exec("delete from tasks where id = $1", i)
	if err != nil {
		return fmt.Errorf("Kunne ikke slette id: %d\nError: %w", i, err)
	}
	return nil
}

func (p *PostgresStore) Update(inTask task.Task) (task.Task, error) {
	qry := "update tasks set title = $1, done = $2, created_at = $3 where id = $4"
	_, err := p.db.Exec(qry, inTask.Title, inTask.Done, inTask.CreatedAt, inTask.ID)
	if err != nil {
		return task.Task{}, fmt.Errorf("Kunne ikke oppdatere task med id %d\nError: %w", inTask.ID, err)
	}
	return inTask, nil
}
