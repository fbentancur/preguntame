package controllers

import (
	"database/sql"
	"log/slog"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/preguntame/preguntame-backend/auth"
	"github.com/preguntame/preguntame-backend/models"
)

type createPostDTO struct {
	OwnerId string `param:"user_id"`
	Title   string `json:"title"`
	Content string `json:"content"`
}

type modifyPostDTO struct {
	OwnerId string `param:"user_id"`
	Title   string `json:"title"`
	Content string `json:"content"`
	PostId  string `param:"post_id"`
}

type deletePostDTO struct {
	PostId string `param:"post_id"`
	UserId string `param:"user_id"`
}

func CreatePost(e echo.Context) error {
	params := createPostDTO{}

	if err := e.Bind(&params); err != nil {
		slog.WarnContext(e.Request().Context(), "Error binding to request", "error", err)
		return err
	}

	uuid, err := uuid.NewUUID()
	if err != nil {
		slog.Error("Error generating uuid", "error", err)
		return err
	}

	loggedUser, err := auth.DecodeUserToken(e)
	if err != nil {
		slog.Warn("Invalid/Missing jwt", "error", err)
		return e.String(http.StatusUnauthorized, "Invalid/Missing jwt")
	}

	if loggedUser.Id != params.OwnerId {
		slog.Warn(
			"Can't post in another's feed",
			"logged_user", loggedUser,
			"user_id", params.OwnerId,
		)
		return e.String(http.StatusForbidden, "Can't post in another's feed")
	}

	post := models.Post{
		Id:           uuid.String(),
		OwnerId:      params.OwnerId,
		Content:      params.Content,
		Title:        params.Title,
		CreationDate: time.Now(),
		DeletionDate: sql.NullTime{Valid: false},
	}

	err = models.InsertPost(post)
	if err != nil {
		slog.Error("Error inserting post into the database", "error", err)
		return err
	}

	return e.String(http.StatusOK, "Post added successfuly")
}

func ModifyPosts(e echo.Context) error {
	params := modifyPostDTO{}

	if err := e.Bind(&params); err != nil {
		slog.WarnContext(e.Request().Context(), "Error binding to request", "error", err)
		return err
	}
	loggedUser, err := auth.DecodeUserToken(e)
	if err != nil {
		slog.Warn("Invalid/Missing jwt", "error", err)
		return e.String(http.StatusUnauthorized, "Invalid/Missing jwt")
	}

	if loggedUser.Id != params.OwnerId {
		slog.Warn(
			"Can't modify another's post",
			"logged_user", loggedUser,
			"user_id", params.OwnerId,
		)
		return e.String(http.StatusForbidden, "Can't modify another's post")
	}

	updated, err := models.UpdatePost(params.OwnerId, params.PostId, params.Content, params.Title)
	if err != nil {
		slog.Error("Error updating post in database", "error", err)
		return err
	}

	if !updated {
		slog.Warn("Tried to modify non existing post", "user_id", params.OwnerId, "post_id", params.PostId)
		return e.String(http.StatusBadRequest, "Post doesn't exists")
	}

	return e.String(http.StatusOK, "Post updated successfuly")

}

func DeletePosts(e echo.Context) error {
	params := deletePostDTO{}

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
			"Can't delete another's post",
			"logged_user", loggedUser,
			"user_id", params.UserId,
			"post_id", params.PostId,
		)

		return e.String(http.StatusForbidden, "Can't delete another's post")
	}

	updated, err := models.SoftDeletePost(params.UserId, params.PostId, time.Now())
	if err != nil {
		slog.Error("Error deleting post in database", "error", err)
		return e.String(http.StatusInternalServerError, "Error deleting post in database")
	}

	if !updated {
		return e.String(http.StatusBadRequest, "Invalid post or user")
	}

	return e.String(http.StatusOK, "Post deleted successfuly")
}
