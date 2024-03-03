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

func main() {
	ctx := context.Background()
	dbConn, err := sql.Open("mysql", "root:root@tcp(localhost:3306)/courses")

	if err != nil {
		panic(err)
	}

	defer dbConn.Close()

	queries := db.New(dbConn)

	err = queries.CreateCategory(ctx, db.CreateCategoryParams{
		ID:   uuid.NewString(),
		Name: "Go",
		Description: sql.NullString{
			String: "Is the best language",
			Valid:  true,
		},
	})

	if err != nil {
		panic(err)
	}

	err = queries.CreateCategory(ctx, db.CreateCategoryParams{
		ID:   uuid.NewString(),
		Name: "Java",
		Description: sql.NullString{
			String: "NullString",
			Valid:  true,
		},
	})

	if err != nil {
		panic(err)
	}

	queries.CreateCategory(ctx, db.CreateCategoryParams{
		ID:   uuid.NewString(),
		Name: "C#",
		Description: sql.NullString{
			String: "NullString",
			Valid:  true,
		},
	})

	all, _ := queries.ListCategories(ctx)

	for _, category := range all {
		fmt.Println(category)
	}

}
