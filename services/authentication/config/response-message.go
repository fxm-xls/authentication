package config

import (
	"bigrule/common/global"
	"bigrule/common/logger"
	"encoding/json"
	"github.com/spf13/viper"
	"os"
	"strconv"
)

func InitResponseMessage() {
	f, err := os.OpenFile(ResMsgConfig.path, os.O_RDONLY, 0755)
	defer f.Close()
	if err != nil {
		logger.Error(err)
	}
	responseMessage := make(map[string]string)
	global.ResponseMessage = make(map[int]string)
	decoder := json.NewDecoder(f)
	err = decoder.Decode(&responseMessage)
	if err != nil {
		logger.Error(err)
	}
	for code, message := range responseMessage {
		codeInt, err := strconv.Atoi(code)
		if err != nil {
			logger.Error(err)
		}
		global.ResponseMessage[codeInt] = message
	}
}

type ResMsg struct {
	path string
}

func InitResMsg(cfg *viper.Viper) *ResMsg {
	return &ResMsg{
		path: cfg.GetString("path"),
	}
}

var ResMsgConfig = new(ResMsg)
