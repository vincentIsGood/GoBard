package utils

import (
	"encoding/json"
	"fmt"
)

func PrintObj(obj any){
    bytes, _ := json.Marshal(obj)
    fmt.Println(string(bytes))
}