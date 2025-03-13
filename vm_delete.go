package main

import (
	"log"
	"net/http"
	"os"
	"os/exec"

	"github.com/scr34m/go-kvm-ui/domain"
)

func deleteVm(w http.ResponseWriter, r *http.Request) {

	vmname := r.URL.Path[len("/delete/"):]

	domain := domain.Load(vmname)
	if domain == nil {
		http.Redirect(w, r, "/?error=Unknown+virtual+machine+\""+vmname+"\"", http.StatusTemporaryRedirect)
		return
	}

	// Stop the VM
	args := make([]string, 0)
	args = append(args, "destroy")
	args = append(args, vmname)

	out, err := exec.Command("/usr/bin/virsh", args...).CombinedOutput()
	if err != nil {
		log.Printf("command: %v, output: %s error: %v", args, out, err)
		http.Redirect(w, r, "/?error=Unable+to+stop+virtual+machine+\""+vmname+"\", error:"+err.Error(), http.StatusTemporaryRedirect)
		return
	}

	// Delete the vm
	args = make([]string, 0)
	args = append(args, "undefine")
	args = append(args, vmname)

	out, err = exec.Command("/usr/bin/virsh", args...).CombinedOutput()
	if err != nil {
		log.Printf("command: %v, output: %s error: %v", args, out, err)
		http.Redirect(w, r, "/?error=Unable+to+delete+virtual+machine+\""+vmname+"\", error:"+err.Error(), http.StatusTemporaryRedirect)
		return
	}

	for _, disk := range domain.Devices.Disks {
		if disk.Device == "disk" {
			err = os.Remove(disk.Source.File)
			if err != nil {
				log.Printf("command: %v, output: %s error: %v", args, out, err)
				http.Redirect(w, r, "/?error=Unable+to+delete+virtual+machine+\""+vmname+"\" disk, error:"+err.Error(), http.StatusTemporaryRedirect)
				return
			}
		}
	}

	http.Redirect(w, r, "/?success=Virtual+machine+\""+vmname+"\"+has+been+deleted", http.StatusTemporaryRedirect)
}
