package program

import (
	"myetcd-manage/program/config"
	"net/http"
	"time"
	"fmt"
	"os/exec"
	"runtime"
	"github.com/opentracing/opentracing-go/log"
	"myetcd-manage/program/logger"
)

// Program 主程序
type Program struct {
	cfg *config.Config
	s   *http.Server
}

// New 创建主程序
func New() (*Program, error) {
	// 配置文件
	cfg, err := config.LoadConfig("")
	if err != nil {
		return nil, err
	}

	// 日志对象
	_, err = logger.InitLogger(cfg.LogPath, cfg.Debug)
	if err != nil {
		return nil, err
	}

	return &Program{
		cfg: cfg,
	}, nil
}


// 打开url
func openURL(urlAddr string) {
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", " /c start "+urlAddr)
	} else if runtime.GOOS == "darwin" {
		cmd = exec.Command("open", urlAddr)
	} else {
		return
	}
	err := cmd.Start()
	if err != nil {
		//logger.Log.Errorw("打开浏览器错误", "err", err)
		log.Error(err)
	}
}


// Run 启动程序
func (p *Program) Run() error {
	// 启动http服务
	go p.startAPI()

	// 打开浏览器
	go func() {
		time.Sleep(100 * time.Millisecond)
		openURL(fmt.Sprintf("http://127.0.0.1:%d/ui/", p.cfg.HTTP.Port))
	}()

	return nil
}



// Stop 停止服务
func (p *Program) Stop() {
	if p.s != nil {
		p.s.Close()
	}
}