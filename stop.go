package main

import (
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"os"
	"path"
	"strconv"
	"syscall"
	"tinydocker/constant"
	"tinydocker/container"
)

func stopContainer(containerId string) {
	// 1. 根据容器Id查询容器信息
	containerInfo, err := container.GetInfoByContainerId(containerId)
	if err != nil {
		log.Errorf("Get container %s info error %v", containerId, err)
		return
	}
	pidInt, err := strconv.Atoi(containerInfo.Pid)
	if err != nil {
		log.Errorf("Conver pid from string to int error %v", err)
		return
	}
	// 2.发送SIGTERM信号
	if err = syscall.Kill(pidInt, syscall.SIGTERM); err != nil {
		log.Errorf("Stop container %s error %v", containerId, err)
		return
	}
	// 3.修改容器信息，将容器置为STOP状态，并清空PID
	containerInfo.Status = container.STOP
	containerInfo.Pid = " "
	newContentBytes, err := json.Marshal(containerInfo)
	if err != nil {
		log.Errorf("Json marshal %s error %v", containerId, err)
		return
	}
	// 4.重新写回存储容器信息的文件
	dirPath := fmt.Sprintf(container.InfoLocFormat, containerId)
	configFilePath := path.Join(dirPath, container.ConfigName)
	if err := os.WriteFile(configFilePath, newContentBytes, constant.Perm0622); err != nil {
		log.Errorf("Write file %s error:%v", configFilePath, err)
	}
}

