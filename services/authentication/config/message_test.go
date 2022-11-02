package config

import (
	"bigrule/common/global"
	"bigrule/common/logger"
	"encoding/json"
	"fmt"
	"os"
	"testing"
)

func TestInitResponseMessage(t *testing.T) {
	f, err := os.OpenFile("../../../config/response-message.json", os.O_RDONLY, 0755)
	defer f.Close()
	if err != nil {
		logger.Error(err)
	}
	//global.ResponseMessage = make(map[string]string)
	decoder := json.NewDecoder(f)
	err = decoder.Decode(&global.ResponseMessage)
	if err != nil {
		fmt.Printf("json decode has error:%v\n", err)
	}
	fmt.Print(global.ResponseMessage)
}
