package main

import (
	"log"
	"net/http"
	"os/exec"
	"strings"

	"github.com/scr34m/go-kvm-ui/domain"
)

// https://github.com/allanrbo/simple-vmcontrol/blob/master/vmcontrol/setautostart.py
func setautostartVm(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/")
	vmname := parts[2]

	domain := domain.Load(vmname)
	if domain == nil {
		http.Redirect(w, r, "/?error=Unknown+virtual+machine+\""+vmname+"\"", 302)
		return
	}

	args := make([]string, 0)
	args = append(args, "autostart")
	args = append(args, vmname)

	if parts[3] == "off" {
		args = append(args, "--disable")
	}

	out, err := exec.Command("/usr/bin/virsh", args...).CombinedOutput()
	if err != nil {
		log.Printf("command: %v, output: %s error: %v", args, out, err)
		http.Redirect(w, r, "/?error=Unable+to+change+autostart+virtual+machine+\""+vmname+"\", error:"+err.Error(), 302)
		return
	}

	http.Redirect(w, r, "/?success=Virtual+machine+\""+vmname+"\"+autostart+status+has+been+changed", 302)
}
