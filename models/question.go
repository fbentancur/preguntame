package models

import (
	"database/sql"

	"github.com/preguntame/preguntame-backend/databases"
)

type QuestionID = string

type Question struct {
	Id        QuestionID
	UserId    UserID
	Message   string
	Reply     sql.NullString
	Favourite bool
	Signature sql.NullString
}

func FindQuestionsByUserId(userId UserID) ([]Question, error) {
	questions := make([]Question, 0, 16)

	query := "SELECT id, target_id, message, reply, favourite FROM Questions WHERE target_id = $1"
	cursor, err := databases.DbPool.Query(query, userId)
	if err != nil {
		return questions, err
	}
	defer cursor.Close()

	for cursor.Next() {
		question := Question{}

		err = cursor.Scan(&question.Id, &question.UserId, &question.Message, &question.Reply, &question.Favourite)
		if err != nil {
			return questions, err
		}

		questions = append(questions, question)
	}

	err = cursor.Err()
	if err != nil {
		return questions, err
	}

	return questions, nil
}

func InsertQuestion(question Question) error {
	stmt := "INSERT INTO Questions(id, target_id, message, reply, favourite, signature) VALUES ($1, $2, $3, $4, $5, $6)"
	_, err := databases.DbPool.Exec(stmt, question.Id, question.UserId, question.Message, question.Reply, question.Favourite, question.Signature)
	return err
}

func UpdateQuestionReply(userId UserID, questionId QuestionID, reply string) (bool, error) {
	stmt := "UPDATE Questions SET reply = $1 WHERE id = $2 AND target_id = $3 AND reply IS NULL"
	result, err := databases.DbPool.Exec(stmt, reply, questionId, userId)
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
func UpdateQuestionFavourite(userId UserID, questionId QuestionID, favourite bool) (bool, error) {
	stmt := "UPDATE Questions SET favourite = $1 WHERE id = $2 AND target_id = $3"
	result, err := databases.DbPool.Exec(stmt, favourite, questionId, userId)
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

func DeleteQuestion(userId UserID, questionId QuestionID) (bool, error) {
	stmt := "DELETE from Questions where id = $1 and target_id =$2"
	result, err := databases.DbPool.Exec(stmt, questionId, userId)
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
