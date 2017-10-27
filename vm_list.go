package main

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/scr34m/go-kvm-ui/domain"
)

func getISO() []string {
	var isoimages []string
	files, err := ioutil.ReadDir(DIR_ISO)
	if err != nil {
		log.Fatal(err)
	}
	isoimages = isoimages[:0]
	for _, f := range files {
		isoimages = append(isoimages, f.Name())
	}
	return isoimages
}

type DiskWeb struct {
	Dev         string
	File        string
	Current     string
	Max         string
	IsImageData bool
}

type DomainWeb struct {
	Name           string
	Memory         string
	Cores          string
	State          string
	Autostart      string
	Vnc            string
	IsRunning      bool
	IsNotRunning   bool
	IsAutostart    bool
	IsNotAutostart bool
	Disks          []DiskWeb
}

func listVm(w http.ResponseWriter, r *http.Request) {
	isoimages := getISO()

	var domainsweb []DomainWeb
	domains := domain.LoadAll()
	for _, d := range domains {
		var diskweb []DiskWeb
		for _, disk := range d.Devices.Disks {
			diskweb = append(diskweb, DiskWeb{
				Dev:         disk.Target.Dev,
				File:        disk.Source.File,
				Current:     fmt.Sprintf("%d", disk.Size/1024),
				Max:         fmt.Sprintf("%d", disk.MaxSize/1024),
				IsImageData: disk.IsImageData(),
			})
		}
		domainsweb = append(domainsweb, DomainWeb{
			Name:           d.Name,
			Memory:         fmt.Sprintf("%d", d.Memory/1024),
			Cores:          fmt.Sprintf("%d", d.Cores),
			State:          d.State,
			Autostart:      d.Autostart,
			Vnc:            fmt.Sprintf("%d", d.GetVNC()),
			IsRunning:      d.State == "running",
			IsNotRunning:   d.State != "running",
			IsAutostart:    d.Autostart == "yes",
			IsNotAutostart: d.Autostart != "yes",
			Disks:          diskweb,
		})
	}

	data := struct {
		IsoImages      []string
		Domains        []DomainWeb
		MessageSuccess string
		MessageError   string
	}{
		IsoImages:      isoimages,
		Domains:        domainsweb,
		MessageSuccess: r.FormValue("success"),
		MessageError:   r.FormValue("error"),
	}
	tmpl, err := template.ParseFiles("tpl/index.html")
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	tmpl.Execute(w, data)
}
