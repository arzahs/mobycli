package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/justincampbell/timeago"
	"io"
	"time"
)

// JSONMessage is struct from docker engine pulling image
type JSONMessage struct {
	Status         string `json:"status,omitempty"`
	ProgressDetail struct {
		Current int64 `json:"current"`
		Total   int64 `json:"total"`
	} `json:"progressDetail"`
}

// PortsToStr transforms port list to docker cli represents
func PortsToStr(ports []types.Port) string {
	fmtPorts := ""
	for i, p := range ports {
		if i > 0 {
			fmtPorts += ", "
		}
		if p.PublicPort != 0 {
			fmtPorts += fmt.Sprintf("%s:%d->%d/%s", p.IP, p.PublicPort, p.PrivatePort, p.Type)
			continue
		}
		fmtPorts += fmt.Sprintf("%d/%s", p.PrivatePort, p.Type)
	}
	return fmtPorts
}

// UnixToStr transform unixt imestamp to docker cli represents
func UnixToStr(unixTime int64) string {
	return timeago.FromTime(time.Unix(unixTime, 0))
}

// PrintToWriter reads JSONMessage from reader and writes
func PrintToWriter(reader io.Reader, writer io.Writer) error {
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		var row JSONMessage
		if err := json.Unmarshal(scanner.Bytes(), &row); err != nil {
			return err
		}
		if row.ProgressDetail.Total > 0 {
			fmt.Fprint(writer, fmt.Sprintf("Status: %s\tProgress: (%d/%d)\n", row.Status, row.ProgressDetail.Current, row.ProgressDetail.Total))
			continue
		}
		fmt.Fprint(writer, fmt.Sprintf("Status: %s\n", row.Status))
	}
	return nil
}
