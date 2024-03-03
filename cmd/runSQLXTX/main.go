package main

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/IcaroSilvaFK/go_sqlc/internal/db"
	"github.com/google/uuid"

	// mysql driver
	_ "github.com/go-sql-driver/mysql"
)

type CourseDB struct {
	dbConn *sql.DB
	*db.Queries
}

type CourseParams struct {
	ID          string
	Name        string
	Description sql.NullString
	Price       float64
}

type CategoryParams struct {
	ID          string
	Name        string
	Description sql.NullString
}

func (c *CourseDB) CreateCourseAndCategory(ctx context.Context, arg CourseParams, arg2 CategoryParams) error {

	err := c.callTx(ctx, func(q *db.Queries) error {
		catId := uuid.NewString()
		err := q.CreateCategory(ctx, db.CreateCategoryParams{
			ID:          catId,
			Name:        arg2.Name,
			Description: arg2.Description,
		})

		if err != nil {
			return err
		}

		err = q.CreateCourse(ctx, db.CreateCourseParams{
			ID:          arg.ID,
			Name:        arg.Name,
			Description: arg.Description,
			Price:       arg.Price,
			CategoryID:  catId,
		})

		return err
	})

	return err
}

func main() {
	ctx := context.Background()
	dbConn, err := sql.Open("mysql", "root:root@tcp(localhost:3306)/courses")

	if err != nil {
		panic(err)
	}

	defer dbConn.Close()

	// tx := NewCourseDB(dbConn)

	queries := db.New(dbConn)

	c, err := queries.ListCourses(ctx)

	if err != nil {
		panic(err)
	}

	for _, cat := range c {
		fmt.Println(cat)
	}

	// c := CourseParams{
	// 	ID:   uuid.NewString(),
	// 	Name: "Go",
	// 	Description: sql.NullString{
	// 		String: "Is the best language",
	// 		Valid:  true,
	// 	},
	// 	Price: 10.0,
	// }

	// cat := CategoryParams{
	// 	ID:   uuid.NewString(),
	// 	Name: "Go",
	// 	Description: sql.NullString{
	// 		String: "Is the best language",
	// 		Valid:  true,
	// 	},
	// }

	// err = tx.CreateCourseAndCategory(ctx, c, cat)

	// if err != nil {
	// 	panic(err)
	// }

}

func NewCourseDB(dbConn *sql.DB) *CourseDB {
	return &CourseDB{
		dbConn:  dbConn,
		Queries: db.New(dbConn),
	}
}

func (c *CourseDB) callTx(ctx context.Context, fn func(*db.Queries) error) error {

	tx, err := c.dbConn.BeginTx(ctx, nil)

	if err != nil {
		return err
	}

	q := db.New(tx)

	if err = fn(q); err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx err: %v, rb err: %v", err, rbErr)
		}

		return err
	}

	return tx.Commit()
}
