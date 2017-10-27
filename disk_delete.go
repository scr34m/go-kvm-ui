package main

import (
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"

	"github.com/scr34m/go-kvm-ui/domain"
)

// https://github.com/allanrbo/simple-vmcontrol/blob/master/vmcontrol/deletedatadisk.py
func deleteDisk(w http.ResponseWriter, r *http.Request) {

	parts := strings.Split(r.URL.Path, "/")

	vmname := parts[2]
	dev := parts[3]

	domain := domain.Load(vmname)
	if domain == nil {
		http.Redirect(w, r, "/?error=Unknown+virtual+machine+\""+vmname+"\"", 302)
		return
	}

	var file string
	for _, disk := range domain.Devices.Disks {
		if disk.Target.Dev == dev && disk.IsImageData() {
			file = disk.Source.File
		}
	}

	if file != "" {
		// Detach the disk image
		args := make([]string, 0)
		args = append(args, "detach-disk")
		args = append(args, vmname)
		args = append(args, dev)
		args = append(args, "--config")

		out, err := exec.Command("/usr/bin/virsh", args...).CombinedOutput()
		if err != nil {
			log.Printf("command: %v, output: %s error: %v", args, out, err)
			http.Redirect(w, r, "/?error=Unable+to+detach+virtual+machine+\""+vmname+"\"+disk, error: "+err.Error(), 302)
			return
		}

		err = os.Remove(file)
		if err != nil {
			log.Printf("command: %v, output: %s error: %v", args, out, err)
			http.Redirect(w, r, "/?error=Unable+to+delete+virtual+machine+\""+vmname+"\"+disk, error: "+err.Error(), 302)
			return
		}
	}

	http.Redirect(w, r, "/?success=Virtual+machine+\""+vmname+"\"+disk+has+been+created", 302)
}
