package main

import (
	"log"
	"net/http"
	"os/exec"
	"strings"

	"github.com/scr34m/go-kvm-ui/domain"
)

// https://github.com/allanrbo/simple-vmcontrol/blob/master/vmcontrol/createdatadisk.py
func mountDisk(w http.ResponseWriter, r *http.Request) {

	parts := strings.Split(r.URL.Path, "/")

	vmname := parts[2]
	isoname := parts[3]

	// TODO check already exsist
	domain := domain.Load(vmname)

	var media = ""
	var cdromdev = ""
	for _, disk := range domain.Devices.Disks {
		if disk.Device == "cdrom" {
			cdromdev = disk.Target.Dev
			media = disk.Source.File
		}
	}

	if cdromdev != "" {
		// Always eject existing ISO if there is any
		if media != "" {
			args := make([]string, 0)
			args = append(args, "change-media")
			args = append(args, vmname)
			args = append(args, cdromdev)
			args = append(args, "--eject")
			args = append(args, "--config")

			out, err := exec.Command("/usr/bin/virsh", args...).CombinedOutput()
			if err != nil {
				log.Fatalf("command: %v, output: %s error: %v", args, out, err)
			}
		}

		if isoname != "" {
			if domain.State == "running" {
				args := make([]string, 0)
				args = append(args, "attach-disk")
				args = append(args, vmname)
				args = append(args, DIR_ISO+isoname)
				args = append(args, cdromdev)
				args = append(args, []string{"--type", "cdrom"}...)

				out, err := exec.Command("/usr/bin/virsh", args...).CombinedOutput()
				if err != nil {
					log.Fatalf("command: %v, output: %s error: %v", args, out, err)
				}
			} else {
				args := make([]string, 0)
				args = append(args, "detach-disk")
				args = append(args, vmname)
				args = append(args, cdromdev)
				args = append(args, "--config")

				out, err := exec.Command("/usr/bin/virsh", args...).CombinedOutput()
				if err != nil {
					log.Fatalf("command: %v, output: %s error: %v", args, out, err)
				}

				args = make([]string, 0)
				args = append(args, "attach-disk")
				args = append(args, vmname)
				args = append(args, DIR_ISO+isoname)
				args = append(args, cdromdev)
				args = append(args, []string{"--type", "cdrom"}...)
				args = append(args, "--config")

				out, err = exec.Command("/usr/bin/virsh", args...).CombinedOutput()
				if err != nil {
					log.Fatalf("command: %v, output: %s error: %v", args, out, err)
				}
			}
		}
	}

	http.Redirect(w, r, "/?error=false", 302)
}
