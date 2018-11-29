package worker

import (
	"sync"

	"github.com/huayuev5/my_go_land/modules/falcon-log-agent/common/dlog"
	"github.com/huayuev5/my_go_land/modules/falcon-log-agent/common/g"
	"github.com/huayuev5/my_go_land/modules/falcon-log-agent/common/scheme"

	"time"

	"github.com/huayuev5/my_go_land/modules/falcon-log-agent/reader"
	"github.com/huayuev5/my_go_land/modules/falcon-log-agent/strategy"
	"encoding/json"
)

// ConfigInfo to control config
type ConfigInfo struct {
	ID       int64
	FilePath string
}

// Job to control job
type Job struct {
	r *reader.Reader
	w *WorkerGroup
}

// ManagerJob to manage jobs
var ManagerJob map[string]*Job //管理job,文件路径为key
// ManagerJobLock is a global lock
var ManagerJobLock *sync.RWMutex

// ManagerConfig to manage configs
var ManagerConfig map[int64]*ConfigInfo

func init() {
	ManagerJob = make(map[string]*Job)
	ManagerJobLock = new(sync.RWMutex)
	ManagerConfig = make(map[int64]*ConfigInfo)
}

// UpdateConfigsLoop to update strategys
func UpdateConfigsLoop() {
	for {
		strategy.Update()
		strategyMap := strategy.GetAll() //最新策略
		ManagerJobLock.Lock()
		strStrategyMap, _ := json.Marshal(strategyMap)
		dlog.Infof("Get new strategyMap %s", string(strStrategyMap))

		for id, st := range strategyMap {
			config := &ConfigInfo{
				ID:       id,
				FilePath: st.FilePath,
			}
			cache := make(chan string, g.Conf().Worker.QueueSize)
			if err := createJob(config, cache, st); err != nil {
				dlog.Errorf("create job fail [id:%d][filePath:%s][err:%v]", config.ID, config.FilePath, err)
			}
		}

		for id := range ManagerConfig {
			if _, ok := strategyMap[id]; !ok { //如果策略中不存在，说明用户已删除
				config := &ConfigInfo{
					ID:       id,
					FilePath: ManagerConfig[id].FilePath,
				}
				deleteJob(config)
			}
		}
		ManagerJobLock.Unlock()

		//更新counter
		GlobalCount.UpdateByStrategy(strategyMap)
		time.Sleep(time.Second * time.Duration(g.Conf().Strategy.UpdateDuration))
	}
}

// GetOldestTms to analysis the oldes tms
func GetOldestTms(filepath string) (int64, bool) {
	ManagerJobLock.RLock()
	defer ManagerJobLock.RUnlock()
	job, ok := ManagerJob[filepath]
	if !ok {
		return 0, false
	}

	tms, allFree := job.w.GetOldestTms()
	nowTms := time.Now().Unix()
	//如果worker全都空闲，且当前时间戳已经领先1min
	//则将标记的时间戳设置为当前时间戳-1min
	if allFree && nowTms-tms > 60 {
		tms = nowTms - 60
	}
	return tms, true
}

//添加任务到管理map( managerjob managerconfig) 启动reader和worker
func createJob(config *ConfigInfo, cache chan string, st *scheme.Strategy) error {
	if _, ok := ManagerJob[config.FilePath]; ok {
		if _, ok := ManagerConfig[config.ID]; !ok {
			ManagerConfig[config.ID] = config
		}
		return nil
	}

	ManagerConfig[config.ID] = config
	//启动reader
	r, err := reader.NewReader(config.FilePath, cache)
	if err != nil {
		return err
	}
	dlog.Infof("Add Reader : [%s]", config.FilePath)
	//启动worker
	w := NewWorkerGroup(config.FilePath, cache, st)
	ManagerJob[config.FilePath] = &Job{
		r: r,
		w: w,
	}
	w.Start()
	//启动reader
	go r.Start()

	dlog.Infof("Create job success [filePath:%s][sid:%d]", config.FilePath, config.ID)
	return nil
}

//先stop worker reader再从管理map中删除
func deleteJob(config *ConfigInfo) {
	//删除jobs
	tag := 0
	for _, cg := range ManagerConfig {
		if config.FilePath == cg.FilePath {
			tag++
		}
	}
	if tag <= 1 {
		dlog.Infof("Del Reader : [%s]", config.FilePath)
		if job, ok := ManagerJob[config.FilePath]; ok {
			job.w.Stop() //先stop worker
			job.r.Stop()
			delete(ManagerJob, config.FilePath)
		}
	}
	dlog.Infof("Stop reader & worker success [filePath:%s][sid:%d]", config.FilePath, config.ID)

	//删除config
	if _, ok := ManagerConfig[config.ID]; ok {
		delete(ManagerConfig, config.ID)
	}
}
