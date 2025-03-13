package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"

	"github.com/scr34m/go-kvm-ui/domain"
)

func createDisk(w http.ResponseWriter, r *http.Request) {

	parts := strings.Split(r.URL.Path, "/")

	var file string
	vmname := parts[2]
	size := parts[3]

	domain := domain.Load(vmname)
	if domain == nil {
		http.Redirect(w, r, "/?error=Unknown+virtual+machine+\""+vmname+"\"", http.StatusTemporaryRedirect)
		return
	}

	i := 1
	for {
		file = fmt.Sprintf(DIR_VM+vmname+".data%d.img", i)
		if _, err := os.Stat(file); os.IsNotExist(err) {
			break
		}
		i += 1
	}

	// Find unique device name
	var names = []string{"vda", "vdb", "vdc", "vdd", "vde", "vdf", "vdg", "vdh", "vdi", "vdj", "vdk", "vdl", "vdm", "vdn", "vdo", "vdp", "vdq", "vdr", "vds", "vdt", "vdu", "vdv", "vdw", "vdx", "vdy", "vdz"}
	var usednames []string

	for _, disk := range domain.Devices.Disks {
		usednames = append(usednames, disk.Target.Dev)
	}

	avail := difference(names, usednames)

	if len(avail) > 1 {
		// Create image file
		args := make([]string, 0)
		args = append(args, "create")
		args = append(args, []string{"-f", "qcow2"}...)
		args = append(args, file)
		args = append(args, size+"G")

		out, err := exec.Command("/usr/bin/qemu-img", args...).CombinedOutput()
		if err != nil {
			log.Printf("command: %v, output: %s error: %v", args, out, err)
			http.Redirect(w, r, "/?error=Unable+to+create+virtual+machine+\""+vmname+"\"+disk, error: "+err.Error(), http.StatusTemporaryRedirect)
			return
		}

		// Attach image file
		args = make([]string, 0)
		args = append(args, "attach-disk")
		args = append(args, vmname)
		args = append(args, file)
		args = append(args, avail[0])
		args = append(args, []string{"--driver", "qemu"}...)
		args = append(args, []string{"--subdriver", "qcow2"}...)
		args = append(args, []string{"--cache", "writeback"}...)
		args = append(args, "--persistent")

		out, err = exec.Command("/usr/bin/virsh", args...).CombinedOutput()
		if err != nil {
			log.Printf("command: %v, output: %s error: %v", args, out, err)
			http.Redirect(w, r, "/?error=Unable+to+attach+virtual+machine+\""+vmname+"\"+disk, error: "+err.Error(), http.StatusTemporaryRedirect)
			return
		}
	}

	http.Redirect(w, r, "/?success=Virtual+machine+\""+vmname+"\"+disk+has+been+created", http.StatusTemporaryRedirect)
}
