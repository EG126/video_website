package utils

import (
	"time"

	"github.com/bwmarrin/snowflake"
)

var node *snowflake.Node

func InitSnowflake(startTime string, machineID int64) error {
	st, err := time.Parse("2006-01-02", startTime)
	if err != nil {
		return err
	}
	//将解析出的时间转换成毫秒级时间戳，并赋值给snowflake.Epoch(一个库提供的全局变量）
	snowflake.Epoch = st.UnixNano()
	node, err = snowflake.NewNode(machineID)
	return err
}

func GenerateID() string {
	if node == nil {
		panic("snowflake not initialized")
	}
	return node.Generate().String()
}
