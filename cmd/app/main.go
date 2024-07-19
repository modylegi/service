package main

import (
	"context"
	"fmt"
	"os"

	_ "github.com/joho/godotenv/autoload"
	_ "github.com/modylegi/service/docs"

	"github.com/modylegi/service/internal/app"
)

//	@title			Content Service API
//	@version		1.0
//	@description	Контент сервис - сервис хранения контента (баннеры, тесты, истории и т.д.) для дальнейшего отображения его на клиенте пользователя.

func main() {
	ctx := context.Background()
	if err := app.Run(ctx, os.Stdout, os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
	}

}
