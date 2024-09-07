package main

import (
	"fmt"

	"github.com/GoogleCloudPlatform/functions-framework-go/funcframework"
	_ "github.com/ilhamhanif/telegram-msg-tracker-sys/app/cf-tlgrm-act-send-message/function"
)

const PORT = "8080"

func main() {

	/*
		A main function to start function-framework
		in selected port.
	*/

	if err := funcframework.Start(PORT); err != nil {
		fmt.Printf("funcframework.Start: %v\n", err)
	}

}
