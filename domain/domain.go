package domain

import (
	"encoding/xml"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

func (domain *Domain) GetVNC() int64 {
	var port int64
	for _, g := range domain.Devices.Graphics {
		if g.Type == "vnc" {
			port, _ = strconv.ParseInt(g.Port, 10, 64)
		}
	}
	return port
}

func (disk *DeviceDisk) IsImageData() bool {
	rx := regexp.MustCompile(".data([0-9]+).img$")
	return rx.MatchString(disk.Source.File)
}

func getXml(name string) []byte {
	out, err := exec.Command("/usr/bin/virsh", "dumpxml", name).CombinedOutput()
	if err != nil {
		log.Fatalf("output: %s error: %v", out, err)
	}
	return out
}

func getDiskSize(file string) (int64, int64) {
	var size int64
	var maxsize int64

	if file == "" || file == "-" {
		return 0, 0
	}

	if strings.HasSuffix(file, ".iso") {
		if fi, err := os.Stat(file); !os.IsNotExist(err) {
			size = fi.Size()
		}
	}

	if _, err := os.Stat(file); !os.IsNotExist(err) {
		out, err := exec.Command("/usr/bin/qemu-img", "info", file).CombinedOutput()
		if err != nil {
			log.Fatalf("output: %s error: %v", out, err)
		}

		rx := regexp.MustCompile(`virtual size: [^\(]+\((\d+)`)
		x := rx.FindStringSubmatch(string(out))
		if len(x) == 2 {
			maxsize, _ = strconv.ParseInt(x[1], 10, 64)
		} else {
			maxsize = size
		}
	} else {
		maxsize = size
	}

	return size, maxsize
}

func parse(name string, state string) Domain {
	var d Domain

	xml.Unmarshal(getXml(name), &d)

	d.State = state

	d.Autostart = "no"
	if _, err := os.Stat("/etc/libvirt/qemu/autostart/" + d.Name + ".xml"); !os.IsNotExist(err) {
		d.Autostart = "yes"
	}

	for idx, disk := range d.Devices.Disks {
		s, m := getDiskSize(disk.Source.File)
		d.Devices.Disks[idx].Size = s
		d.Devices.Disks[idx].MaxSize = m
	}
	return d
}

func Load(vmname string) *Domain {
	var domain Domain

	out, err := exec.Command("/usr/bin/virsh", "list", "--all").CombinedOutput()
	if err != nil {
		log.Fatalf("output: %s error: %v", out, err)
	}

	rx := regexp.MustCompile(`([0-9\-]+)\s+([a-z0-9]+)\s+(running|blocked|paused|shutdown|shut off|crashed|inactive)`)
	lines := strings.Split(string(out), string('\n'))
	for _, l := range lines[2:] {
		x := rx.FindStringSubmatch(l)
		if len(x) != 4 {
			continue
		}

		if x[2] == vmname {
			domain = parse(x[2], x[3])
			return &domain
		}
	}

	return nil
}

func LoadAll() []Domain {
	var domains []Domain

	out, err := exec.Command("/usr/bin/virsh", "list", "--all").CombinedOutput()
	if err != nil {
		log.Fatalf("output: %s error: %v", out, err)
	}

	rx := regexp.MustCompile(`([0-9\-]+)\s+([a-z0-9]+)\s+(running|blocked|paused|shutdown|shut off|crashed|inactive)`)
	lines := strings.Split(string(out), string('\n'))
	for _, l := range lines[2:] {
		x := rx.FindStringSubmatch(l)
		if len(x) != 4 {
			continue
		}

		d := parse(x[2], x[3])
		domains = append(domains, d)
	}

	return domains
}
