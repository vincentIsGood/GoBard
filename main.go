package main

import (
	"bardtest/bardapp"
	"flag"
	"fmt"
	"os"
)

func main(){
    infoCommand := flag.NewFlagSet("info", flag.ContinueOnError)
    sessionForInfo := infoCommand.String("session", "", "'__Secure-1PSID' cookie from google")
    chatCommand := flag.NewFlagSet("chat", flag.ContinueOnError)
    sessionForChat := chatCommand.String("session", "", "'__Secure-1PSID' cookie from google")
    messageForChat := chatCommand.String("text", "こんにちは", "your message to bard")

    if len(os.Args) == 1{
        showHelpAndExit(infoCommand, chatCommand)
    }

    switch os.Args[1]{
    case "info":
        infoCommand.Parse(os.Args[2:])
        bard := bardapp.New(*sessionForInfo)
        bard.FetchInfo()
    case "chat":
        chatCommand.Parse(os.Args[2:])
        bard := bardapp.New(*sessionForChat)
        bard.FetchInfo()
        bard.SendMessage(*messageForChat)
        bard.PrintResponses()
    default:
        showHelpAndExit(infoCommand, chatCommand)
    }

    // Then, search for "WIZ_global_data"
}

func showHelpAndExit(flags... *flag.FlagSet){
    fmt.Printf("[!] Usage: %s <info / chat> [options]\n", os.Args[0])
    for _, flag := range flags{
        flag.Usage()
    }
    os.Exit(1)
}

/// Note on Go arrays
// Java Syntax: Cookie[] cookies = new Cookie[]{new Cookie(...)}
// cookies := []*http.Cookie{{
//     Name: "__Secure-1PSID",
//     Value: sessionId,
// }}