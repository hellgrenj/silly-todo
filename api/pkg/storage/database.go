package storage

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/hellgrenj/silly-todo/pkg/domainerrors"
	"github.com/hellgrenj/silly-todo/pkg/todo"
	"github.com/jackc/pgconn"
	pgx "github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

// Database is the struct for the database
type Database struct {
	pool *pgxpool.Pool
}

// NewDatabase sets up connection pool (retries connection max 5 times over 20 seconds, depending on what you pass in as connectionAttempt)
func NewDatabase(connectionAttempt int) *Database {

	connectionString := os.Getenv("POSTGRES_CONNECTION_STRING")
	if len(connectionString) == 0 {
		connectionString = "postgresql://silly:silly@localhost/silly"
	}
	dbpool, err := pgxpool.Connect(context.Background(), connectionString)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		if connectionAttempt < 5 {
			connectionAttempt++
			fmt.Printf("Trying again in 4 seconds attempt %v of 5\n", connectionAttempt)
			time.Sleep(4 * time.Second)
			return NewDatabase(connectionAttempt)
		}
		os.Exit(1)

	}
	fmt.Println("successfully connected to database")
	return &Database{
		pool: dbpool,
	}
}

// AddList adds a new todolist in the database and returns the id for the created list or an error
func (d *Database) AddList(list todo.List) (int, error) {
	var id int

	err := d.pool.QueryRow(context.Background(), "INSERT INTO list (name) VALUES ($1) RETURNING id", list.Name).Scan(&id)
	if err != nil {
		pgErr := err.(*pgconn.PgError)
		switch pgErr.Code {
		case "23505":
			expl := fmt.Sprintf("A list with name %s already exists", list.Name)
			return 0, domainerrors.NewConflictError(expl, err)
		default:
			return 0, err
		}
	}
	return id, nil
}

// AddItem adds an item in the database (for the todo list with the provided listID) and returns the created item id or an error
func (d *Database) AddItem(item todo.Item, listID int) (int, error) {
	var id int
	err := d.pool.QueryRow(context.Background(), "INSERT INTO item (name, done, list_id) VALUES ($1, $2, $3) RETURNING id", item.Name, item.Done, listID).Scan(&id)
	if err != nil {
		pgErr := err.(*pgconn.PgError)
		switch pgErr.Code {
		case "23503":
			return 0, domainerrors.NewNotFoundError(fmt.Sprintf("list with id %v", listID), err)
		default:
			return 0, err
		}
	}
	return id, nil
}

// GetAllLists gets all todo lists (with items) from database
func (d *Database) GetAllLists() ([]todo.List, error) {
	rows, err := d.pool.Query(context.Background(), "SELECT list.id, list.name, COALESCE (item.id, 0), COALESCE (item.name, ''), COALESCE (item.done, false) FROM list LEFT JOIN item ON (item.list_id = list.id)")
	if err != nil {
		return nil, err
	}
	allLists, err := dbRowsToTodoLists(rows)
	if err != nil {
		return nil, err
	}
	return allLists, err
}

// GetListByID finds a list by id
func (d *Database) GetListByID(id int) (todo.List, error) {
	rows, err := d.pool.Query(context.Background(), "SELECT list.id, list.name, COALESCE (item.id, 0), COALESCE (item.name, ''), COALESCE (item.done, false) FROM list LEFT JOIN item ON (item.list_id = list.id) WHERE list.id = $1", id)
	if err != nil {
		return todo.List{}, err
	}
	lists, err := dbRowsToTodoLists(rows)
	if err != nil {
		return todo.List{}, err
	}
	if len(lists) == 0 {
		return todo.List{}, domainerrors.NewNotFoundError(fmt.Sprintf("list with id %v", id), errors.New("list not found"))
	}
	return lists[0], err
}

// DeleteListByID removes a list by id
func (d *Database) DeleteListByID(id int) error {
	rows, err := d.pool.Query(context.Background(), "DELETE FROM list WHERE list.id = $1 RETURNING id", id)
	defer rows.Close()
	if err != nil {
		pgErr := err.(*pgconn.PgError)
		switch pgErr.Code {
		case "23503":
			return domainerrors.NewNotFoundError(fmt.Sprintf("list with id %v", id), err)
		default:
			return err
		}
	}
	return nil
}

// DeleteItemByID removes an item by id
func (d *Database) DeleteItemByID(id int) error {
	rows, err := d.pool.Query(context.Background(), "DELETE FROM item WHERE item.id = $1", id)
	defer rows.Close()
	if err != nil {
		pgErr := err.(*pgconn.PgError)
		switch pgErr.Code {
		case "23503":
			return domainerrors.NewNotFoundError(fmt.Sprintf("item with id %v", id), err)
		default:
			return err
		}
	}
	return nil
}

// UpdateItemDone updates a single todolist items done field
func (d *Database) UpdateItemDone(id int, done bool) error {
	rows, err := d.pool.Query(context.Background(), "UPDATE item SET done = $2 WHERE item.id = $1", id, done)
	defer rows.Close()
	if err != nil {
		pgErr := err.(*pgconn.PgError)
		switch pgErr.Code {
		case "23503":
			return domainerrors.NewNotFoundError(fmt.Sprintf("item with id %v", id), err)
		default:
			return err
		}
	}
	return nil
}
func dbRowsToTodoLists(rows pgx.Rows) ([]todo.List, error) {
	var allLists []todo.List
	var dataset = make(map[int]*todo.List)

	for rows.Next() {
		l := todo.List{}
		i := todo.Item{}
		err := rows.Scan(&l.ID, &l.Name, &i.ID, &i.Name, &i.Done)
		if err != nil {
			return nil, err
		}

		if _, ok := dataset[l.ID]; ok {
			if itemNotNull(i) {
				dataset[l.ID].Items = append(dataset[l.ID].Items, i)
			}

		} else {
			if itemNotNull(i) {
				l.Items = append(l.Items, i)
			}
			dataset[l.ID] = &l
		}

	}
	// map to list
	for _, v := range dataset {
		allLists = append(allLists, *v)
	}
	return allLists, nil
}
func itemNotNull(i todo.Item) bool {
	return len(i.Name) > 0 && i.ID != 0
}
