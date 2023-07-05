package bardapp

import (
	"bardtest/httpclient"
	"bardtest/utils"
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

const BARD_URL string = "https://bard.google.com"
const BARD_CHAT_URL string = "https://bard.google.com/_/BardChatUi/data/assistant.lamda.BardFrontendService/StreamGenerate"

type BardResponse struct{
    choiceId string
    text string
}

func (resp *BardResponse) String() string{
    return fmt.Sprintf("BardResponse{choiceId: %s, text: '%s'}", resp.choiceId, resp.text)
}

type BardClient struct{
    client *httpclient.SimpleHttpClient
    SNlM0e string
    model string
    reqId int

    conversationId string
    responseId string
    selectedChoiceId string

    lastMessageSent string
    interpretedMessage string
    lastResponses map[string](*BardResponse)
}

func New(sessionId string) *BardClient{
    fmt.Printf("[*] Using session id '%s' for Bard\n", sessionId)
    client := httpclient.New()
    client.AddCookie(BARD_URL, "__Secure-1PSID", sessionId)

    reqId, _ := strconv.Atoi(utils.RandomChoices(strings.Split("0123456789", ""), 6))

    bard := &BardClient{
        client: client,
        model: "boq_assistant-bard-web-server_20230627.10_p1",
        reqId: reqId,
        conversationId: "",
        responseId: "",
        selectedChoiceId: "",
        
        lastMessageSent: "",
        interpretedMessage: "",
        lastResponses: make(map[string]*BardResponse),
    }
    return bard
}

func (bard *BardClient) FetchInfo(){
    // res, err := http.Get(BARD_URL)
    fmt.Printf("[+] Sending request to '%s'\n", BARD_URL)
    res, err := bard.client.SendRequest("GET", BARD_URL)
    utils.PanicOnError(err)
    defer res.Body.Close()

    respBodyBytes, err := io.ReadAll(res.Body)
    utils.PanicOnError(err)

    bard.SNlM0e = utils.GetParams(regexp.MustCompile(`"SNlM0e":"(?P<SNlM0e>[^\"]*)",`), string(respBodyBytes))["SNlM0e"]
    fmt.Printf("[+] Found 'SNlM0e' value: '%s'\n", bard.SNlM0e)

    if bard.SNlM0e == "" || len(bard.SNlM0e) > 200{
        panic("[-] 'SNlM0e' value should not be empty or too large (Is your session token correct?)")
    }
}

func (bard *BardClient) SendMessage(userMessage string){
    userMessage = strings.TrimSpace(userMessage)
    if len(userMessage) < 2{
        panic("[-] Cannot send a message with length < 2 to Bard")
    }

    bardChatUrlParams := fmt.Sprintf("?bl=%s&_reqid=%d&rt=c", bard.model, bard.reqId)
    bardChatUrl := BARD_CHAT_URL + bardChatUrlParams

    fmt.Printf("[+] Sending another request to '%s'\n", bardChatUrl)

    messageArray := []interface{}{
        []interface{}{userMessage, 0, nil, []interface{}{}},
        []interface{}{"jp"},
        []interface{}{bard.conversationId, bard.responseId, bard.selectedChoiceId},
        nil, nil, nil, 
        []interface{}{1}, 0, []interface{}{},
    }
    messageArrayBytes, _ := json.Marshal(messageArray)
    postBody, _ := json.Marshal([]interface{}{nil, string(messageArrayBytes)})
    
    res, _ := bard.client.SendRequestWithHeadersAndBody("POST", bardChatUrl, map[string]string{
        "Content-Type": "application/x-www-form-urlencoded;charset=UTF-8",
    }, httpclient.EncodePostForm(map[string]string{
        "f.req": string(postBody),
        "at": bard.SNlM0e,
    }))

    bard.lastMessageSent = userMessage
    bard.reqId += 10000
    bard.parsePostResponse(res)
}

func (bard *BardClient) PrintResponses(){
    fmt.Println("======== Printing Response ========")
    fmt.Printf("[+] Current conversation id: '%s'\n", bard.conversationId)
    fmt.Printf("[+] Current response id: '%s'\n", bard.responseId)
    fmt.Printf("[+] Current choice id: '%s'\n", bard.selectedChoiceId)
    fmt.Println("[+] Message Sent:")
    fmt.Println(bard.lastMessageSent)
    fmt.Println("====================================")
    // counter := 0
    // for choiceId, response := range bard.lastResponses{
    //     fmt.Printf("[+] Response %d: id('%s')\n", counter+1, choiceId)
    //     fmt.Println(response.text)
    //     counter++
    // }
    fmt.Println("[+] Chosen Response:")
    fmt.Println(bard.lastResponses[bard.selectedChoiceId].text)
}

func (bard *BardClient) parsePostResponse(res *http.Response){
    defer res.Body.Close()

    bodyScanner := bufio.NewScanner(res.Body)
    bodyScanner.Split(bufio.ScanLines)
    
    var responseJson string
    for bodyScanner.Scan(){
        line := bodyScanner.Text()
        if strings.HasPrefix(line, "["){
            // take the 1st occurance ONLY
            responseJson = line
            break
        }
    }

    // [["wrb.fr", null, "[...]"]]
    respJsonObj := make([]interface{}, 0)
    utils.PanicOnError(json.Unmarshal([]byte(responseJson), &respJsonObj))

    msgObj := respJsonObj[0].([]interface{})

    if msgObj[0] != "wrb.fr"{
        panic("[-] Received an invalid json response")
    }

    // [null, ["c", "r"], [["userinput"], ...], null, [["choice1", ["resp"]], ...], ...]
    payloadObj := make([]interface{}, 0)
    utils.PanicOnError(json.Unmarshal([]byte(msgObj[2].(string)), &payloadObj))

    basicIds := payloadObj[1].([]interface{})
    bard.conversationId = basicIds[0].(string)
    bard.responseId = basicIds[1].(string)

    if payloadObj[2] == nil{
        panic(fmt.Sprintf("[-] Bard do not support the language for '%s'", bard.lastMessageSent))
    }

    userInputs := payloadObj[2].([]interface{})
    bard.interpretedMessage = userInputs[0].([]interface{})[0].(string)

    bardResponse := payloadObj[4].([]interface{})

    for i := 0; i < len(bardResponse); i++{
        responseArray := bardResponse[i].([]interface{})
        choiceId := responseArray[0].(string)
        responseText := responseArray[1].([]interface{})[0].(string)
        bard.lastResponses[choiceId] = &BardResponse{
            choiceId: choiceId,
            text: responseText,
        }
        
        if i == 0{
            bard.selectedChoiceId = choiceId
        }
    }
}