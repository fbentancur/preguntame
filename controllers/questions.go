package controllers

import (
	"database/sql"
	"log/slog"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/preguntame/preguntame-backend/auth"
	"github.com/preguntame/preguntame-backend/models"
)

type findQuestionsDTO struct {
	UserId string `param:"user_id"`
}

type askQuestionDTO struct {
	UserId    string `param:"user_id"`
	Message   string `json:"message"`
	Signature string `json:"signature"`
}

type replyQuestionDTO struct {
	UserId     string `param:"user_id"`
	QuestionId string `param:"question_id"`
	Message    string `json:"message"`
}

type questionDTO struct {
	Id      string  `json:"id"`
	Message string  `json:"message"`
	Reply   *string `json:"reply"`
}

type favouriteDTO struct {
	UserId     string `param:"user_id"`
	QuestionId string `param:"question_id"`
}

type deleteDTO struct {
	QuestionId string `param:"question_id"`
	UserId     string `param:"user_id"`
}

func FindQuestionsForUser(e echo.Context) error {
	params := findQuestionsDTO{}

	if err := e.Bind(&params); err != nil {
		slog.WarnContext(e.Request().Context(), "Error binding to request", "error", err)
		return err
	}

	questions, err := models.FindQuestionsByUserId(params.UserId)
	if err != nil {
		slog.Error("Error getting questions from db", "error", err)
		return err
	}

	response := make([]questionDTO, len(questions))

	for i, question := range questions {
		response[i] = questionToDto(question)
	}

	return e.JSON(http.StatusOK, response)
}

func AskQuestionToUser(e echo.Context) error {
	params := askQuestionDTO{}

	if err := e.Bind(&params); err != nil {
		slog.WarnContext(e.Request().Context(), "Error binding to request", "error", err)
		return err
	}

	uuid, err := uuid.NewUUID()
	if err != nil {
		slog.Error("Error generating uuid", "error", err)
		return err
	}

	if len(params.Message) < 10 {
		slog.Info("Length of questions must be greater or equals than 10", "question", params.Message)
		return e.String(http.StatusBadRequest, "Length of questions must be greater or equals than 10")
	}

	if len(params.Message) > 1000 {
		slog.Info("The length of questions must be less or equal than 1000", "question", params.Message)
		return e.String(http.StatusBadRequest, "The length of questions must be less or equal than 1000")
	}

	signature := sql.NullString{Valid: false}

	if params.Signature != "" {
		signature.Valid = true
		signature.String = params.Signature
	}

	question := models.Question{
		Id:        uuid.String(),
		UserId:    params.UserId,
		Message:   params.Message,
		Reply:     sql.NullString{Valid: false},
		Favourite: false,
		Signature: signature,
	}

	err = models.InsertQuestion(question)
	if err != nil {
		slog.Error("Error inserting question into the database", "error", err)
		return err
	}

	return e.String(http.StatusOK, "Question asked successfuly")
}

func ReplyQuestionToUser(e echo.Context) error {
	params := replyQuestionDTO{}

	if err := e.Bind(&params); err != nil {
		slog.WarnContext(e.Request().Context(), "Error binding to request", "error", err)
		return err
	}

	loggedUser, err := auth.DecodeUserToken(e)
	if err != nil {
		slog.Warn("Invalid/Missing jwt", "error", err)
		return e.String(http.StatusUnauthorized, "Invalid/Missing jwt")
	}

	if len(params.Message) < 10 {
		slog.Info("Length of reply must be greater or equals than 10", "question", params.Message)
		return e.String(http.StatusBadRequest, "Length of reply must be greater or equals than 10")
	}

	if len(params.Message) > 1000 {
		slog.Info("The length of reply must be less or equal than 1000", "question", params.Message)
		return e.String(http.StatusBadRequest, "The length of reply must be less or equal than 1000")
	}

	if loggedUser.Id != params.UserId {
		slog.Warn(
			"Can't reply another's question",
			"logged_user", loggedUser,
			"user_id", params.UserId,
			"question_id", params.QuestionId,
		)
		return e.String(http.StatusForbidden, "Can't reply another's question")
	}

	updated, err := models.UpdateQuestionReply(params.UserId, params.QuestionId, params.Message)
	if err != nil {
		slog.Error("Error updating question in database")
		return err
	}

	if !updated {
		slog.Warn("Tried to answer non existing question", "user_id", params.UserId, "question_id", params.QuestionId)
		return e.String(http.StatusBadRequest, "Question doesn't exists")
	}

	return e.String(http.StatusOK, "Question updated successfuly")
}

func questionToDto(question models.Question) questionDTO {
	var reply *string = nil
	if question.Reply.Valid {
		reply = &question.Reply.String
	}

	return questionDTO{
		Id:      question.Id,
		Message: question.Message,
		Reply:   reply,
	}
}

func MakeFavourite(e echo.Context) error {
	params := favouriteDTO{}

	if err := e.Bind(&params); err != nil {
		slog.WarnContext(e.Request().Context(), "Error binding to request", "error", err)
		return err
	}

	loggedUser, err := auth.DecodeUserToken(e)
	if err != nil {
		slog.Warn("Invalid/Missing jwt", "error", err)
		return e.String(http.StatusUnauthorized, "Invalid/Missing jwt")
	}

	if loggedUser.Id != params.UserId {
		slog.Warn(
			"Can't set to favourite another's question",
			"logged_user", loggedUser,
			"user_id", params.UserId,
			"question_id", params.QuestionId,
		)

		return e.String(http.StatusForbidden, "Can't set to favourite another's question")
	}

	updated, err := models.UpdateQuestionFavourite(params.UserId, params.QuestionId, true)
	if err != nil {
		return e.String(http.StatusInternalServerError, "Error updating question favourite in database")
	}
	if !updated {
		return e.String(http.StatusBadRequest, "Invalid question or user")
	}

	return e.String(http.StatusOK, "Question updated successfuly")
}

func DeleteQuestion(e echo.Context) error {
	params := deleteDTO{}

	if err := e.Bind(&params); err != nil {
		slog.WarnContext(e.Request().Context(), "Error binding to request", "error", err)
		return err
	}

	loggedUser, err := auth.DecodeUserToken(e)
	if err != nil {
		slog.Warn("Invalid/Missing jwt", "error", err)
		return e.String(http.StatusUnauthorized, "Invalid/Missing jwt")
	}

	if loggedUser.Id != params.UserId {
		slog.Warn(
			"Can't delete another's question",
			"logged_user", loggedUser,
			"user_id", params.UserId,
			"question_id", params.QuestionId,
		)

		return e.String(http.StatusForbidden, "Can't delete another's question")
	}

	updated, err := models.DeleteQuestion(params.UserId, params.QuestionId)
	if err != nil {
		return e.String(http.StatusInternalServerError, "Error deleting question in database")
	}

	if !updated {
		return e.String(http.StatusBadRequest, "Invalid question or user")
	}

	return e.String(http.StatusOK, "Question deleted successfuly")
}
