package main

import (
	"fmt"
	"log"
	"net/http"
	"os/exec"
	"strconv"

	"github.com/scr34m/go-kvm-ui/domain"
)

// https://github.com/allanrbo/simple-vmcontrol/blob/master/vmcontrol/createvm.py
func createVm(w http.ResponseWriter, r *http.Request) {

	vmname := r.FormValue("vmname")

	domain := domain.Load(vmname)
	if domain != nil {
		http.Redirect(w, r, "/?error=Virtual+machine+\""+vmname+"\"+already+exists", 302)
		return
	}

	if len(vmname) < 1 || !IsLetter(vmname) {
		http.Redirect(w, r, "/?error=Wrong characters in name \""+vmname+"\"", 302)
		return
	}

	cores, _ := strconv.ParseInt(r.FormValue("cores")[0:], 10, 64)
	if cores < 1 || cores > MAX_CORES {
		http.Redirect(w, r, "/?error=Cores count not in 1 and "+fmt.Sprintf("%d", MAX_CORES), 302)
		return
	}

	memory, _ := strconv.ParseInt(r.FormValue("memory")[0:], 10, 64)
	if memory < 1 || memory > MAX_MEMORY {
		http.Redirect(w, r, "/?error=Memory size not in 1 and "+fmt.Sprintf("%d", MAX_MEMORY), 302)
		return
	}

	osdisksize, _ := strconv.ParseInt(r.FormValue("osdisksize")[0:], 10, 64)
	if osdisksize < 1 || osdisksize > MAX_DISK {
		http.Redirect(w, r, "/?error=Disk size not in 1 and "+fmt.Sprintf("%d", MAX_DISK), 302)
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
		http.Redirect(w, r, "/?error=Unable+to+create+virtual+machine+\""+vmname+"\"+OS+disk, error: "+err.Error(), 302)
		return
	}

	// Create the VM
	args = make([]string, 0)
	args = append(args, "--connect=qemu:///system")
	args = append(args, []string{"--name", vmname}...)
	args = append(args, []string{"--ram", strconv.FormatInt(memory, 10)}...)
	args = append(args, []string{"--vcpus", strconv.FormatInt(cores, 10)}...)
	// --disk path=/media/kvm/squeeze1.img,device=disk,bus=virtio,size=4 \
	args = append(args, []string{"--disk", "path=" + DIR_VM + vmname + ".os.img,format=qcow2,bus=virtio,cache=writeback"}...)
	// --disk path=/media/kvm/debian-6.0.10-amd64-netinst.iso,device=cdrom \
	// --cdrom /media/kvm/debian-6.0.10-amd64-netinst.iso \
	args = append(args, []string{"-c", DIR_ISO + installiso}...)
	args = append(args, "--network=bridge:br0")
	args = append(args, []string{"--graphics", "vnc,listen=0.0.0.0,keymap=hu-hu"}...)
	// --os-type=linux
	args = append(args, "--noautoconsole")
	args = append(args, "--accelerate")
	args = append(args, "--boot=cdrom,hd")
	args = append(args, "--noapic")
	args = append(args, "--hvm")

	out, err = exec.Command("/usr/bin/virt-install", args...).CombinedOutput()
	if err != nil {
		log.Printf("command: %v, output: %s error: %v", args, out, err)
		http.Redirect(w, r, "/?error=Unable+to+install+virtual+machine+\""+vmname+"\", error: "+err.Error(), 302)
		return
	}

	http.Redirect(w, r, "/?success=Virtual+machine+\""+vmname+"\"+has+been+created", 302)
}
