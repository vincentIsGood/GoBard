package tests_test

import (
	"bardtest/utils"
	"fmt"
	"regexp"
	"testing"
)

func TestParameterizedRegex(t *testing.T){
    sampleTestData := `,"SNlM0e":"TestingMyCodeToMakeSureIamNotWastingMyOwnTime",`

    fmt.Print(utils.GetParams(regexp.MustCompile(`"SNlM0e":"(?P<SNlM0e>[^\"]*)",`), sampleTestData))
}