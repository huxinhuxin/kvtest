package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"sync"
)

var list1 []context.CancelFunc
var mu sync.Mutex

func startproc(args []string, file string, wg *sync.WaitGroup) {

	wg.Add(1)
	mu.Lock()
	ctx, cancel := context.WithCancel(context.Background())
	list1 = append(list1, cancel)
	//args := []string{"-d","-k", "3"}
	mu.Unlock()
	cmd := exec.CommandContext(ctx, "iostat", args...)
	ss, _ := cmd.CombinedOutput()

	ioutil.WriteFile(file, ss, 0666)
	wg.Done()
}

type disk struct {
	device  string
	tps     float64
	kbread  float64
	kbwrite float64
	kbdscd  float64
}

func parsedisk() map[string]*disk {
	file, err := os.Open("data/disk")
	if err != nil {
		fmt.Println(err.Error())
		return nil
	}
	//ret:=make([]*disk,0)
	ret := make(map[string][]*disk)
	reder := bufio.NewReader(file)
	rows := 0
	for {
		buf, err := reder.ReadString('\n')
		if err != nil {
			if io.EOF == err {
				break
			}
			fmt.Println(err.Error())
			return nil
		}

		str := string(buf)
		//str = strings.ReplaceAll(str, "  ", " ")

		arr := strings.Fields(str)

		if len(arr) == 0 {
			rows = 0
			continue
		}

		if arr[0] == "Device" {
			rows = 1
			continue
		}
		if rows == 0 {
			continue
		}

		dis := &disk{}
		dis.device = arr[0]
		dis.tps, _ = strconv.ParseFloat(arr[1], 64)
		dis.kbread, _ = strconv.ParseFloat(arr[2], 64)
		dis.kbwrite, _ = strconv.ParseFloat(arr[3], 64)
		dis.kbdscd, _ = strconv.ParseFloat(arr[4], 64)

		_, ok := ret[dis.device]
		if ok {
			ret[dis.device] = append(ret[dis.device], dis)
		} else {
			silce := make([]*disk, 0)
			silce = append(silce, dis)
			ret[dis.device] = silce
		}
	}

	//求平均值
	avg := make(map[string]*disk)
	for key, val := range ret {
		val = val[1:]
		avg[key] = &disk{device: key}

		for _, v1 := range val {
			avg[key].kbdscd += v1.kbdscd
			avg[key].tps += v1.tps
			avg[key].kbread += v1.kbread
			avg[key].kbwrite += v1.kbwrite
		}
		size := len(val)
		if size > 0 {
			avg[key].kbdscd /= float64(size)
			avg[key].tps /= float64(size)
			avg[key].kbread /= float64(size)
			avg[key].kbwrite /= float64(size)
		}
	}
	return avg
}

func printdisk(m map[string]*disk) {
	fmt.Printf("\n%-20s%-20s%-20s%-20s%-20s\n", "Device", "tps", "kB_read/s", "kB_wrtn/s", "kB_dscd/s")
	for _, v := range m {
		fmt.Printf("%-20s%-20.2f%-20.2f%-20.2f%-20.2f\n", v.device, v.tps, v.kbread, v.kbwrite, v.kbdscd)
	}
}

type disk1 struct {
	device string
	ws     float64
	wkbs   float64
}

func parsedisk1() map[string]*disk1 {
	file, err := os.Open("data/diskx")
	if err != nil {
		fmt.Println(err.Error())
		return nil
	}
	//ret:=make([]*disk,0)
	ret := make(map[string][]*disk1)
	reder := bufio.NewReader(file)
	rows := 0
	for {
		buf, err := reder.ReadString('\n')
		if err != nil {
			if io.EOF == err {
				break
			}
			fmt.Println(err.Error())
			return nil
		}

		str := string(buf)
		//str = strings.ReplaceAll(str, "  ", " ")

		arr := strings.Fields(str)

		if len(arr) == 0 {
			rows = 0
			continue
		}

		if arr[0] == "Device" {
			rows = 1
			continue
		}
		if rows == 0 {
			continue
		}

		dis := &disk1{}
		dis.device = arr[0]
		dis.ws, _ = strconv.ParseFloat(arr[7], 64)
		dis.wkbs, _ = strconv.ParseFloat(arr[8], 64)

		_, ok := ret[dis.device]
		if ok {
			ret[dis.device] = append(ret[dis.device], dis)
		} else {
			silce := make([]*disk1, 0)
			silce = append(silce, dis)
			ret[dis.device] = silce
		}
	}

	//求平均值
	avg := make(map[string]*disk1)
	for key, val := range ret {
		val = val[1:]
		avg[key] = &disk1{device: key}

		for _, v1 := range val {
			avg[key].ws += v1.ws
			avg[key].wkbs += v1.wkbs

		}
		size := len(val)
		if size > 0 {
			avg[key].ws /= float64(size)
			avg[key].wkbs /= float64(size)
		}
	}
	return avg
}

func printdisk1(m map[string]*disk1) {
	fmt.Printf("\n%-20s%-20s%-20s\n", "Device", "w/s", "wk/s")
	for _, v := range m {
		fmt.Printf("%-20s%-20.2f%-20.2f\n", v.device, v.ws, v.wkbs)
	}
}

type cpu struct {
	user   float64
	nice   float64
	system float64
	iowait float64
	steal  float64
	idle   float64
}

func parsecpu() *cpu {
	file, err := os.Open("data/cpu")
	if err != nil {
		fmt.Println(err.Error())
		return nil
	}
	//ret:=make([]*disk,0)
	ret := make([]*cpu, 0)
	reder := bufio.NewReader(file)
	rows := 0
	for {
		buf, err := reder.ReadString('\n')
		if err != nil {
			if io.EOF == err {
				break
			}
			fmt.Println(err.Error())
			return nil
		}

		str := string(buf)
		//str = strings.ReplaceAll(str, "  ", " ")

		arr := strings.Fields(str)

		if len(arr) == 0 {
			rows = 0
			continue
		}

		if arr[0] == "avg-cpu:" {
			rows = 1
			continue
		}
		if rows == 0 {
			continue
		}

		cpu := &cpu{}
		cpu.idle, _ = strconv.ParseFloat(arr[5], 64)
		cpu.iowait, _ = strconv.ParseFloat(arr[3], 64)
		cpu.nice, _ = strconv.ParseFloat(arr[1], 64)
		cpu.steal, _ = strconv.ParseFloat(arr[4], 64)
		cpu.system, _ = strconv.ParseFloat(arr[2], 64)
		cpu.user, _ = strconv.ParseFloat(arr[0], 64)

		ret = append(ret, cpu)
	}

	//求平均值

	avg := &cpu{}

	for _, v1 := range ret {
		avg.idle += v1.idle
		avg.iowait += v1.iowait
		avg.nice += v1.nice
		avg.steal += v1.steal
		avg.system += v1.system
		avg.user += v1.user
	}
	size := len(ret)
	if size > 0 {
		avg.idle /= float64(size)
		avg.iowait /= float64(size)
		avg.nice /= float64(size)
		avg.steal /= float64(size)
		avg.system /= float64(size)
		avg.user /= float64(size)
	}

	return avg
}

func printcpu(m *cpu) {
	fmt.Printf("\n%-20s%-20s%-20s%-20s%-20s%-20s%-20s\n", "avg-cpu:", "%user", "%nice", "%system", "%iowait", "%steal", "%idle")

	fmt.Printf("%-20s%-20.2f%-20.2f%-20.2f%-20.2f%-20.2f%-20.2f\n", "", m.user, m.nice, m.system, m.iowait, m.steal, m.idle)

}
