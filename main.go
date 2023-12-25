package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"

	"github.com/narasimha-1511/go-microservices/appilication"
)

// to start the redis serve us e the following
// command :: docker run --name=redis-devel --publish=6379:6379 --hostname=redis --restart=on-failure --detach redis:latest

func main(){	
	// TODO: Implement
	app:= appilication.New(appilication.LoadConfig());

	ctx ,cancel := signal.NotifyContext(context.Background(), os.Interrupt);

	defer cancel()

	error:= app.Start(ctx);

	if error!=nil{
		fmt.Println("Error starting the application",error);
	}

}	
