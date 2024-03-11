package main

import (
	"log/slog"

	"github.com/labstack/echo/v4"
	"github.com/preguntame/preguntame-backend/controllers"
	"github.com/preguntame/preguntame-backend/databases"
)

func main() {
	if err := databases.InitDatabase(); err != nil {
		slog.Error("Error initializing DataBase connection pool", "error", err)
		return
	}

	e := echo.New()

	e.POST("/users/login", controllers.Login)
	e.POST("/users/register", controllers.Register)

	e.GET("/users/:user_id/questions", controllers.FindQuestionsForUser)
	e.POST("/users/:user_id/questions", controllers.AskQuestionToUser)
	e.PUT("/users/:user_id/questions/:question_id", controllers.ReplyQuestionToUser)
	e.PUT("/users/:user_id/questions/:question_id/fav", controllers.MakeFavourite)
	e.DELETE("/users/:user_id/questions/:question_id", controllers.DeleteQuestion)

	e.Logger.Fatal(e.Start(":8080"))
}
