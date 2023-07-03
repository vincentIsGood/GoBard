package utils

import (
	"math/rand"
	"regexp"
	"strings"
)

func PanicOnError(err error){
    if err != nil{
        panic(err)
    }
}

func RandomChoices(choices []string, outputLen int) string{
    choicesLen := len(choices)
    result := make([]string, 0)
    for i := 0; i < outputLen; i++{
        result = append(result, choices[rand.Intn(choicesLen)])
    }
    return strings.Join(result, "")
}

func GetParams(regex *regexp.Regexp, matchAgainst string) map[string]string{
    matches := regex.FindStringSubmatch(matchAgainst)
    resultParams := make(map[string]string)
    for i, name := range regex.SubexpNames(){
        if i > 0 && i < len(matches){
            resultParams[name] = matches[i]
        }
    }
    return resultParams
}