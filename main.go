package main

import (
	"fmt"
	"github.com/gofiber/fiber"
	"io/ioutil"
	"log"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"
)

type CPUStatsJiffies struct {
	user    int64 // Time spent with normal processing in user mode.
	nice    int64 // Time spent with niced processes in user mode.
	system  int64 // Time spent running in kernel mode.
	idle    int64 // Time spent in idle.
	iowait  int64 // Time spent waiting for I/O to completed. This is considered idle time too.
	irq     int64 // Time spent serving hardware interrupts. See the description of the intr line for more details.
	softirq int64 // Time spent serving software interrupts.
	steal   int64 // Time stolen by other operating systems running in a virtual environment.
	guest   int64 //Time spent for running a virtual CPU or guest OS under the control of the kernel.
}
type CPUStatsHuman struct {
	user    int64 // Time spent with normal processing in user mode.
	nice    int64 // Time spent with niced processes in user mode.
	system  int64 // Time spent running in kernel mode.
	idle    int64 // Time spent in idle.
	iowait  int64 // Time spent waiting for I/O to completed. This is considered idle time too.
	irq     int64 // Time spent serving hardware interrupts. See the description of the intr line for more details.
	softirq int64 // Time spent serving software interrupts.
	steal   int64 // Time stolen by other operating systems running in a virtual environment.
	guest   int64 //Time spent for running a virtual CPU or guest OS under the control of the kernel.
}

func getCPUSample() (idle, total uint64, contents []byte) {
	contents, err := ioutil.ReadFile("/proc/stat")

	//fmt.Println(contents)
	if err != nil {
		return
	}
	lines := strings.Split(string(contents), "\n")
	for _, line := range lines {
		fields := strings.Fields(line)
		if fields[0] == "cpu" {
			numFields := len(fields)
			for i := 1; i < numFields; i++ {
				val, err := strconv.ParseUint(fields[i], 10, 64)
				if err != nil {
					fmt.Println("Error: ", i, fields[i], err)
				}
				total += val // tally up all the numbers to get total ticks
				if i == 4 {  // idle is the 5th field in the cpu line
					idle = val
				}
			}
			return
		}
	}
	return
}

func parseCPUSample(cpuStats string) (idle, total uint64) {
	lines := strings.Split(string(cpuStats), "\n")
	for _, line := range lines {
		fields := strings.Fields(line)
		if fields[0] == "cpu" {
			numFields := len(fields)
			for i := 1; i < numFields; i++ {
				val, err := strconv.ParseUint(fields[i], 10, 64)
				if err != nil {
					fmt.Println("Error: ", i, fields[i], err)
				}
				total += val // tally up all the numbers to get total ticks
				if i == 4 {  // idle is the 5th field in the cpu line
					idle = val
				}
			}
			return
		}
	}
	return
}
func getCPUDelta(idle0, total0, idle1, total1 uint64) (float64, float64) {
	idleTicks := float64(idle1 - idle0)
	totalTicks := float64(total1 - total0)
	cpuUsage := 100 * (totalTicks - idleTicks) / totalTicks
	cpuIdle := 100 * idleTicks / totalTicks
	return cpuUsage, cpuIdle
}

func logCPUUsage() {

	f, err := os.OpenFile("sysstat.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	idle0, total0, contents := getCPUSample()
	if _, err := f.Write(contents); err != nil {
		log.Fatal(err)
	}
	for i := 1; i < 10; i++ {
		time.Sleep(1 * time.Second)
		idle1, total1, contents := getCPUSample()
		if _, err := f.Write(contents); err != nil {
			log.Fatal(err)
		}

		idleTicks := float64(idle1 - idle0)
		totalTicks := float64(total1 - total0)
		cpuUsage := 100 * (totalTicks - idleTicks) / totalTicks

		fmt.Printf("CPU usage is %f%% [busy: %f, total: %f]\n", cpuUsage, totalTicks-idleTicks, totalTicks)
		idle0, total0 = idle1, total1
	}
	if err := f.Close(); err != nil {
		log.Fatal(err)
	}
}

func main() {
	if runtime.GOOS != "linux" {
		log.Fatal("This code must run under linux. Use Docker.")
	}
	// Fiber instance
	app := fiber.New()

	// Routes
	app.Get("/", hello)

	// Start server
	log.Fatal(app.Listen(3000))
}

// Handler
func hello(c *fiber.Ctx) {
	idle0, total0, _ := getCPUSample()
	time.Sleep(1 * time.Second)
	idle1, total1, _ := getCPUSample()
	total_usage, total_idle := getCPUDelta(idle0, total0, idle1, total1)
	c.Send(fmt.Sprintf("Total CPU: %f; Total idle: %f", total_usage, total_idle))
}
