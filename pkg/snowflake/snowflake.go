package snowflake

import (
	"time"

	sf "github.com/bwmarrin/snowflake"
)

// 在一台机器上启动一个node节点
var node *sf.Node

func Init(startTime string, machineID int64) (err error) {
	var st time.Time
	st, err = time.Parse("2006-01-02", startTime) //开始时间
	if err != nil {
		return
	}
	sf.Epoch = st.UnixNano() / 1000000
	node, err = sf.NewNode(machineID) //通过机器id生成node节点
	return
}
func GenID() int64 {
	return node.Generate().Int64()
}
