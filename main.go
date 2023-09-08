package main

import (
	"bardtest/bardapp"
	"flag"
	"fmt"
	"os"
)

func main(){
    infoCommand := flag.NewFlagSet("info", flag.ContinueOnError)
    session1ForInfo := infoCommand.String("session1", "", "'__Secure-1PSID' cookie from google")
    session2ForInfo := infoCommand.String("session2", "", "'__Secure-1PSIDTS' cookie from google")

    chatCommand := flag.NewFlagSet("chat", flag.ContinueOnError)
    session1ForChat := chatCommand.String("session1", "", "'__Secure-1PSID' cookie from google")
    session2ForChat := chatCommand.String("session2", "", "'__Secure-1PSIDTS' cookie from google")
    messageForChat := chatCommand.String("text", "こんにちは", "your message to bard")

    if len(os.Args) == 1{
        showHelpAndExit(infoCommand, chatCommand)
    }

    switch os.Args[1]{
    case "info":
        infoCommand.Parse(os.Args[2:])
        bard := bardapp.New(*session1ForInfo, *session2ForInfo)
        bard.FetchInfo()
    case "chat":
        chatCommand.Parse(os.Args[2:])
        bard := bardapp.New(*session1ForChat, *session2ForChat)
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