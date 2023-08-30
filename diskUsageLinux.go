//go:build linux
// +build linux

package main

import (
	"encoding/json"
	"golang.org/x/sys/unix"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

type diskUsageJobStruct struct {
	Folder    string
	DiskThreshold string
}

func checkDiskMinSpace(jobToProcess *job) {
	var diskUsageJob diskUsageJobStruct
	_ = json.Unmarshal([]byte(jobToProcess.JobData), &diskUsageJob)
	var stat unix.Statfs_t
	wd := diskUsageJob.Folder
	unix.Statfs(wd, &stat)
	// Available blocks * size per block = available space in bytes
	freeBytesAvailable := stat.Bavail * uint64(stat.Bsize)
	totalNumberOfBytes := stat.Blocks * uint64(stat.Bsize)

	if strings.HasSuffix(diskUsageJob.DiskThreshold, "%") {
		diskUsageThresholdPercent, _ := strconv.Atoi(strings.Trim(diskUsageJob.DiskThreshold, "%"))
		diskUsageThresholdPercentFloat := float64(diskUsageThresholdPercent)
		usedBytes := int(totalNumberOfBytes) - int(freeBytesAvailable)
		currentShare := float64(usedBytes) / float64(totalNumberOfBytes)
		currentFreeShare := (1 - currentShare) * 100
		if currentFreeShare < diskUsageThresholdPercentFloat && time.Now().After(jobToProcess.lastTriggered.Add(3*time.Hour)) {
			jobToProcess.lastTriggered = time.Now()
			sendData(GetAESEncrypted(hostname + ": Low Free Disk Space: <" + diskUsageJob.DiskThreshold))

		}
	} else {
		var multiplier int
		if strings.HasSuffix(diskUsageJob.DiskThreshold, "B") {
			multiplier = 1
		}
		if strings.HasSuffix(diskUsageJob.DiskThreshold, "K") {
			multiplier = 1024
		}
		if strings.HasSuffix(diskUsageJob.DiskThreshold, "M") {
			multiplier = 1024 * 1024
		}
		if strings.HasSuffix(diskUsageJob.DiskThreshold, "G") {
			multiplier = 1024 * 1024 * 1024
		}
		if strings.HasSuffix(diskUsageJob.DiskThreshold, "T") {
			multiplier = 1024 * 1024 * 1024 * 1024
		}
		if strings.HasSuffix(diskUsageJob.DiskThreshold, "P") {
			multiplier = 1024 * 1024 * 1024 * 1024 * 1024
		}
		diskUsageThresholdBytes := strings.Trim(diskUsageJob.DiskThreshold, "BKMGTP")
		diskUsageThresholdBytesInt, _ := strconv.Atoi(diskUsageThresholdBytes)
		diskUsageThresholdBytesUint64 := uint64(diskUsageThresholdBytesInt) * uint64(multiplier)
		if freeBytesAvailable < diskUsageThresholdBytesUint64 && time.Now().After(jobToProcess.lastTriggered.Add(3*time.Hour)) {
			jobToProcess.lastTriggered = time.Now()
			sendData(GetAESEncrypted(hostname + ": Low Free Disk Space: <" + diskUsageJob.DiskThreshold))
		}
	}
}

func checkRAID(jobToProcess *job) {
	out, _ := exec.Command("cat", "/proc/mdstat").Output()
	outSlice := strings.Split(string(out),"\n")
	hasRAIDError := false
	for _,line := range(outSlice) {
		if strings.Contains(line,"_") {
			hasRAIDError = true
		}
	}
	if hasRAIDError && time.Now().After(jobToProcess.lastTriggered.Add(3*time.Hour)) {
		jobToProcess.lastTriggered = time.Now()
		sendData(GetAESEncrypted(hostname + ": RAID ERROR:" + string(out)))
	}
}
