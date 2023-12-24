package main

import (
	"context"
	"fmt"

	"github.com/narasimha-1511/go-microservices/appilication"
)

func main(){	
	// TODO: Implement
	app:= appilication.New();
	error:= app.Start(context.TODO());

	if error!=nil{
		fmt.Println("Error starting the application",error);
	}

}	
