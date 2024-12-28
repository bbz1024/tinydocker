package main

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"os"
	"tinydocker/container"
)

func removeContainer(containerId string, force bool) {
	containerInfo, err := container.GetInfoByContainerId(containerId)
	if err != nil {
		log.Errorf("Get container %s info error %v", containerId, err)
		return
	}

	switch containerInfo.Status {
	case container.STOP: // STOP 状态容器直接删除即可
		// 获取容器元数据
		dirPath := fmt.Sprintf(container.InfoLocFormat, containerId)
		// 删除挂载
		container.DeleteWorkSpace(dirPath, containerInfo.Volume)

		if err = os.RemoveAll(dirPath); err != nil {
			log.Errorf("Remove file %s error %v", dirPath, err)
			return
		}

	case container.RUNNING: // RUNNING 状态容器如果指定了 force 则先 stop 然后再删除
		if !force {
			log.Errorf("Couldn't remove running container [%s], Stop the container before attempting removal or"+
				" force remove", containerId)
			return
		}
		stopContainer(containerId)
		removeContainer(containerId, force)
	default:
		log.Errorf("Couldn't remove container,invalid status %s", containerInfo.Status)
		return
	}
}
