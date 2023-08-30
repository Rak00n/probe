//go:build windows
// +build windows

package main

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"golang.org/x/sys/windows"
)

type diskUsageJobStruct struct {
	Folder    string
	DiskThreshold string
}

func checkDiskMinSpace(jobToProcess *job) {
	var freeBytesAvailable uint64
	var totalNumberOfBytes uint64
	var totalNumberOfFreeBytes uint64

	var diskUsageJob diskUsageJobStruct
	_ = json.Unmarshal([]byte(jobToProcess.JobData), &diskUsageJob)

	_ = windows.GetDiskFreeSpaceEx(windows.StringToUTF16Ptr(diskUsageJob.Folder), &freeBytesAvailable, &totalNumberOfBytes, &totalNumberOfFreeBytes)
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
	fmt.Println("Function is not implemented. It won't be unless anyone request it :)")
}
