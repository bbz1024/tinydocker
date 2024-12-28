package main

import (
	"math/rand"
	"os"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
	"tinydocker/cgroups"
	"tinydocker/cgroups/subsystems"
	"tinydocker/container"
)

// Run 执行具体 command
/*
这里的Start方法是真正开始执行由NewParentProcess构建好的command的调用，它首先会clone出来一个namespace隔离的
进程，然后在子进程中，调用/proc/self/exe,也就是调用自己，发送init参数，调用我们写的init方法，
去初始化容器的一些资源。
*/
func Run(tty bool, comArray []string, volume, containerName string, res *subsystems.ResourceConfig) {
	containerId := container.GenerateContainerID() // 生成 10 位容器 id

	//log.Info("1-", os.Getpid()) // tinydocker的pid

	// -------------------- 创建子进程 --------------------
	parent, writePipe := container.NewParentProcess(tty, volume, containerId,containerName)
	if parent == nil {
		log.Errorf("New parent process error")
		return
	}
	//log.Info("2-", os.Getpid()) // tinydocker的pid

	// 这里会进行clone一个完全隔离的子进程，进程pid==1，start执行的完成
	if err := parent.Start(); err != nil {
		log.Errorf("Run parent.Start err:%v", err)
	}
	//log.Info("3-", os.Getpid())        // tinydocker的pid
	//log.Info("4-", parent.Process.Pid) // sh的pid
	// --------------------  --------------------

	// record container info
	err := container.RecordContainerInfo(parent.Process.Pid, comArray, containerName, containerId)
	if err != nil {
		log.Errorf("Record container info error %v", err)
		return
	}
	// -------------------- cgroup 注册 --------------------
	// 创建cgroup manager, 并通过调用set和apply设置资源限制并使限制在容器上生效
	cgroupManager := cgroups.NewCgroupManager("tinydocker-cgroup")
	defer func() {
		//log.Info("5-", os.Getpid())
		// 当子进程退出时，删除cgroup，执行这个操作
		cgroupManager.Destroy()
		//log.Infof("exit success")
	}()
	_ = cgroupManager.Set(res)
	_ = cgroupManager.Apply(parent.Process.Pid, res)

	// 在子进程创建后才能通过pipe来发送参数
	sendInitCommand(comArray, writePipe)

	// 如果是tty，那么父进程等待，就是前台运行，否则就是跳过，实现后台运行
	if tty {
		_ = parent.Wait()
		//log.Info("6-", os.Getpid())
		container.DeleteWorkSpace("/root/", volume)
		container.DeleteContainerInfo(containerId)

	}
}

// sendInitCommand 通过writePipe将指令发送给子进程
func sendInitCommand(comArray []string, writePipe *os.File) {
	command := strings.Join(comArray, " ")
	log.Infof("command all is %s", command)
	_, _ = writePipe.WriteString(command)
	_ = writePipe.Close()
}
func randStringBytes(n int) string {
	letterBytes := "1234567890"
	rand.Seed(time.Now().UnixNano())
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}
