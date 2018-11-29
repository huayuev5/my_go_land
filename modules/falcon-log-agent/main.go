package main

import (
	"github.com/huayuev5/my_go_land/modules/falcon-log-agent/http"

	"github.com/huayuev5/my_go_land/modules/falcon-log-agent/common/proc/metric"
	"github.com/huayuev5/my_go_land/modules/falcon-log-agent/common/utils"

	"github.com/huayuev5/my_go_land/modules/falcon-log-agent/common/dlog"
	"github.com/huayuev5/my_go_land/modules/falcon-log-agent/common/g"
	"github.com/huayuev5/my_go_land/modules/falcon-log-agent/worker"

	"runtime"
)

func main() {
	g.InitAll()
	defer g.CloseLog()

	maxCoreNum := utils.GetCPULimitNum(g.Conf().MaxCPURate)
	dlog.Infof("bind [%d] cpu core", maxCoreNum)
	runtime.GOMAXPROCS(maxCoreNum)

	go metric.MetricLoop(60)
	go worker.UpdateConfigsLoop()
	//go patrol.PatrolLoop()
	go worker.PusherStart()

	http.Start()
}
