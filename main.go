package main

import (
	"flag"
	"log"
	"net/http"
)

const MAX_CORES = 4
const MAX_MEMORY = 8 * 1204
const MAX_DISK = 400 * 1204

const DIR_ISO = "/media/disk/kvm/iso/"
const DIR_VM = "/media/disk/kvm/vm/"

var listen = flag.String("listen", ":2015", "Location to listen for connections")

func main() {
	flag.Parse()

	http.HandleFunc("/", listVm)
	http.HandleFunc("/create", createVm)
	http.HandleFunc("/stop/", stopVm)
	http.HandleFunc("/start/", startVm)
	http.HandleFunc("/setautostart/", setautostartVm)
	http.HandleFunc("/delete/", deleteVm)
	http.HandleFunc("/createdisk/", createDisk)
	http.HandleFunc("/deletedisk/", deleteDisk)
	http.HandleFunc("/mountiso/", mountDisk)

	http.HandleFunc("/websockify/", wsh)

	fs := http.FileServer(http.Dir("assets"))
	http.Handle("/assets/", http.StripPrefix("/assets/", fs))

	log.Fatal(http.ListenAndServe(*listen, nil))
}
