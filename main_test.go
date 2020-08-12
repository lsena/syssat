package main

import (
	"bufio"
	"bytes"
	"log"
	"math"
	"os"
	"testing"
)

func getMockProcStat(scanner *bufio.Scanner) string {
	var buffer bytes.Buffer
	for i := 0; i < 9; i++ {
		scanner.Scan()
		buffer.WriteString(scanner.Text())
		buffer.WriteString("\n")
	}
	return buffer.String()
}

func TestCPUUsage(t *testing.T) {
	f, err := os.Open("sysstat_static.log")
	expectedIdleOutputs := [10]float64{19.76, 26.80, 22.33, 15.46, 16.49}
	expectedCpuOutputs := [10]float64{80.23, 73.19, 77.66, 84.53, 83.50}
	if err != nil {
		t.Errorf(err.Error())
		log.Fatal(err)
	}

	scanner := bufio.NewScanner(f)

	idle0, total0 := parseCPUSample(getMockProcStat(scanner))
	for i := 0; i < 5; i++ {
		idle1, total1 := parseCPUSample(getMockProcStat(scanner))
		if idle, total := getCPUDelta(idle0, total0, idle1, total1);
			math.Floor(idle*100)/100 != expectedIdleOutputs[i] || math.Floor(total*100)/100 != expectedCpuOutputs[i] {
			t.Errorf("parseCPUSample = %f %f; want %f %f", idle, total, expectedIdleOutputs[i], expectedCpuOutputs[i])
		}
		if err := scanner.Err(); err != nil {
			log.Fatal(err)
		}
		idle0, total0 = idle1, total1
	}
	if err := f.Close(); err != nil {
		log.Fatal(err)
	}
}
