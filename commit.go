package main

import (
	log "github.com/sirupsen/logrus"
	"os/exec"
)

func commitContainer(imageName string) {
	mntPath := "/root/merged"
	imageTar := "/root/" + imageName + ".tar"
	log.Println("commitContainer imageTar:", imageTar)
	if _, err := exec.Command("tar", "-czf", imageTar, "-C", mntPath, ".").CombinedOutput(); err != nil {
		log.Errorf("tar folder %s error %v", mntPath, err)
	}
}
