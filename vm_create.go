package main

import (
	"fmt"
	"log"
	"net/http"
	"os/exec"
	"strconv"

	"github.com/scr34m/go-kvm-ui/domain"
)

func createVm(w http.ResponseWriter, r *http.Request) {

	vmname := r.FormValue("vmname")

	domain := domain.Load(vmname)
	if domain != nil {
		http.Redirect(w, r, "/?error=Virtual+machine+\""+vmname+"\"+already+exists", http.StatusTemporaryRedirect)
		return
	}

	if len(vmname) < 1 || !IsLetter(vmname) {
		http.Redirect(w, r, "/?error=Wrong characters in name \""+vmname+"\"", http.StatusTemporaryRedirect)
		return
	}

	cores, _ := strconv.ParseInt(r.FormValue("cores")[0:], 10, 64)
	if cores < 1 || cores > MAX_CORES {
		http.Redirect(w, r, "/?error=Cores count not in 1 and "+fmt.Sprintf("%d", MAX_CORES), http.StatusTemporaryRedirect)
		return
	}

	memory, _ := strconv.ParseInt(r.FormValue("memory")[0:], 10, 64)
	if memory < 1 || memory > MAX_MEMORY {
		http.Redirect(w, r, "/?error=Memory size not in 1 and "+fmt.Sprintf("%d", MAX_MEMORY), http.StatusTemporaryRedirect)
		return
	}

	osdisksize, _ := strconv.ParseInt(r.FormValue("osdisksize")[0:], 10, 64)
	if osdisksize < 1 || osdisksize > MAX_DISK {
		http.Redirect(w, r, "/?error=Disk size not in 1 and "+fmt.Sprintf("%d", MAX_DISK), http.StatusTemporaryRedirect)
		return
	}

	installiso := r.FormValue("installiso")

	// Create the OS disk image
	// TODO check already exsist
	args := make([]string, 0)
	args = append(args, "create")
	args = append(args, []string{"-f", "qcow2"}...)
	args = append(args, DIR_VM+vmname+".os.img")
	args = append(args, strconv.FormatInt(osdisksize, 10)+"G")

	out, err := exec.Command("/usr/bin/qemu-img", args...).CombinedOutput()
	if err != nil {
		log.Printf("command: %v, output: %s error: %v", args, out, err)
		http.Redirect(w, r, "/?error=Unable+to+create+virtual+machine+\""+vmname+"\"+OS+disk, error: "+err.Error(), http.StatusTemporaryRedirect)
		return
	}

	// Create the VM
	args = make([]string, 0)
	args = append(args, "--connect=qemu:///system")
	args = append(args, []string{"--name", vmname}...)
	args = append(args, []string{"--ram", strconv.FormatInt(memory, 10)}...)
	args = append(args, []string{"--vcpus", strconv.FormatInt(cores, 10)}...)
	args = append(args, []string{"--disk", "path=" + DIR_VM + vmname + ".os.img,format=qcow2,bus=virtio,cache=writeback"}...)
	args = append(args, []string{"-c", DIR_ISO + installiso}...)
	args = append(args, "--network="+NETWORK)
	args = append(args, []string{"--graphics", "vnc,passwd=" + VNC_PASSWORD + ",listen=" + VNC_ADDRESS + ",keymap=" + KEYMAP}...)
	args = append(args, "--noautoconsole")
	args = append(args, "--accelerate")
	args = append(args, "--boot=cdrom,hd")
	args = append(args, "--noapic")
	args = append(args, "--hvm")

	out, err = exec.Command("/usr/bin/virt-install", args...).CombinedOutput()
	if err != nil {
		log.Printf("command: %v, output: %s error: %v", args, out, err)
		http.Redirect(w, r, "/?error=Unable+to+install+virtual+machine+\""+vmname+"\", error: "+err.Error(), http.StatusTemporaryRedirect)
		return
	}

	http.Redirect(w, r, "/?success=Virtual+machine+\""+vmname+"\"+has+been+created", http.StatusTemporaryRedirect)
}
