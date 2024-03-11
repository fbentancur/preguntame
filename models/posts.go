package models

import (
	"database/sql"
	"time"

	"github.com/preguntame/preguntame-backend/databases"
)

type PostID = string

type Post struct {
	Id           PostID
	OwnerId      UserID
	Title        string
	Content      string
	CreationDate time.Time
	DeletionDate sql.NullTime
}

func InsertPost(post Post) error {
	stmt := "INSERT INTO Posts(id, content, title, owner_Id, creation_date, deletion_date) VALUES ($1, $2, $3, $4, $5, $6)"
	_, err := databases.DbPool.Exec(stmt, post.Id, post.Content, post.Title, post.OwnerId, post.CreationDate, post.DeletionDate)
	return err
}

//En sql null no es comparable con ningun otro valor por lo tanto el operador = no es aplicable, en su lugar se utiliza
//el operador IS.

func UpdatePost(ownerId UserID, postID PostID, content string, title string) (bool, error) {
	stmt := "UPDATE Posts SET content = $1, title = $2 WHERE id = $3 AND owner_id = $4 AND deletion_date IS null"
	result, err := databases.DbPool.Exec(stmt, content, title, postID, ownerId)
	if err != nil {
		return false, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return false, err
	}
	if rowsAffected != 1 {
		return false, nil
	}

	return true, nil
}

func SoftDeletePost(ownerId UserID, postId PostID, deletionTime time.Time) (bool, error) {
	stmt := "UPDATE Posts SET deletion_date = $3 WHERE id=$1 AND owner_id =$2 AND deletion_date IS null"
	result, err := databases.DbPool.Exec(stmt, postId, ownerId, deletionTime)
	if err != nil {
		return false, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return false, err
	}
	if rowsAffected != 1 {
		return false, nil
	}

	return true, nil

}
