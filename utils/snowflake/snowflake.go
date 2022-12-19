package snowflake

import (
	sf "github.com/bwmarrin/snowflake"
	"time"
)

var node *sf.Node

func Init(startTime string, machineID int64) (err error) {
	t, err := time.Parse("2006-01-02", startTime)
	if err != nil {
		return
	}
	// 修改开始时间
	sf.Epoch = t.UnixMilli()
	node, err = sf.NewNode(machineID)
	if err != nil {
		return
	}
	return nil
}

func GenID() int64 {
	return node.Generate().Int64()
}
