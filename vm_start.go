package main

import (
	"log"
	"net/http"
	"os/exec"

	"github.com/scr34m/go-kvm-ui/domain"
)

func startVm(w http.ResponseWriter, r *http.Request) {
	vmname := r.URL.Path[len("/start/"):]

	domain := domain.Load(vmname)
	if domain == nil {
		http.Redirect(w, r, "/?error=Unknown+virtual+machine+\""+vmname+"\"", http.StatusTemporaryRedirect)
		return
	}

	if domain.State != "running" {
		args := make([]string, 0)
		args = append(args, "start")
		args = append(args, vmname)

		out, err := exec.Command("/usr/bin/virsh", args...).CombinedOutput()
		if err != nil {
			log.Printf("command: %v, output: %s error: %v", args, out, err)
			http.Redirect(w, r, "/?error=Unable+to+start+virtual+machine+\""+vmname+"\", error:"+err.Error(),
				http.StatusTemporaryRedirect)
			return
		}
	}

	http.Redirect(w, r, "/?success=Virtual+machine+\""+vmname+"\"+has+been+started", http.StatusTemporaryRedirect)
}
