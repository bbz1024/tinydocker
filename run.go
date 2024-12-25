package main

import (
	"os"
	"strings"

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
func Run(tty bool, comArray []string, res *subsystems.ResourceConfig) {
	// -------------------- 创建子进程 --------------------
	parent, writePipe := container.NewParentProcess(tty)
	if parent == nil {
		log.Errorf("New parent process error")
		return
	}
	if err := parent.Start(); err != nil {
		log.Errorf("Run parent.Start err:%v", err)
	}
	//fmt.Println("4", os.Getpid())        // 这里是 tinydocker的程序id
	//fmt.Println("5", parent.Process.Pid) // run -it /bin/sh 这里的id指的是 /bin/sh 的id
	// --------------------  --------------------

	// -------------------- cgroup 注册 --------------------
	// 创建cgroup manager, 并通过调用set和apply设置资源限制并使限制在容器上生效
	cgroupManager := cgroups.NewCgroupManager("tinydocker-cgroup")
	defer func() {
		// 当子进程退出时，删除cgroup，执行这个操作
		cgroupManager.Destroy()
		log.Infof("exit success")
	}()
	_ = cgroupManager.Set(res)
	_ = cgroupManager.Apply(parent.Process.Pid, res)

	// 在子进程创建后才能通过pipe来发送参数
	sendInitCommand(comArray, writePipe)
	_ = parent.Wait()
}

// sendInitCommand 通过writePipe将指令发送给子进程
func sendInitCommand(comArray []string, writePipe *os.File) {
	command := strings.Join(comArray, " ")
	log.Infof("command all is %s", command)
	_, _ = writePipe.WriteString(command)
	_ = writePipe.Close()
}
