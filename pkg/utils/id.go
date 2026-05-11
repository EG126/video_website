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
	snowflake.Epoch = st.UnixMilli()
	node, err = snowflake.NewNode(machineID)
	return err
}

func GenerateID() string {
	if node == nil {
		panic("snowflake not initialized")
	}
	return node.Generate().String()
}
