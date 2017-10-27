package domain

type OS struct {
	Type OSType `xml:"type"`
	Boot []Boot `xml:"boot"`
}

type OSType struct {
	Type        string `xml:",chardata"`
	TypeArch    string `xml:"arch,attr"`
	TypeMachine string `xml:"machine,attr"`
}

type Boot struct {
	Dev string `xml:"dev,attr"`
}

type Target struct {
	Dev string `xml:"dev,attr"`
	Bus string `xml:"bus,attr"`
}

type Disk struct {
	Target Target `xml:"target"`
}

type CPU struct {
	Mode  string   `xml:"mode,attr"`
	Match string   `xml:"match,attr"`
	Model CPUModel `xml:"model"`
}

type CPUModel struct {
	Fallback string `xml:"fallback,attr"`
	Value    string `xml:",chardata"`
}

type Clock struct {
	Offset string       `xml:"offset,attr"`
	Timer  []ClockTimer `xml:"timer"`
}

type ClockTimer struct {
	Name       string `xml:"name,attr"`
	Tickpolicy string `xml:"tickpolicy,attr"`
	Present    string `xml:"present,attr"`
}

type Device struct {
	Emulator string           `xml:"emulator"`
	Disks    []DeviceDisk     `xml:"disk"`
	Graphics []DeviceGraphics `xml:"graphics"`
}

type DeviceDisk struct {
	Type   string       `xml:"type,attr"`
	Device string       `xml:"device,attr"`
	Driver DeviceDriver `xml:"driver"`
	Source DeviceSource `xml:"source"`
	Target DeviceTarget `xml:"target"`
	// <readonly/>
	Address DeviceAddress `xml:"address"`
	Size    int64
	MaxSize int64
}

type DeviceDriver struct {
	Name  string `xml:"name,attr"`
	Type  string `xml:"type,attr"`
	Cache string `xml:"cache,attr"`
}

type DeviceSource struct {
	File string `xml:"file,attr"`
}

type DeviceTarget struct {
	Dev string `xml:"dev,attr"`
	Bus string `xml:"bus,attr"`
}

type DeviceAddress struct {
	Type     string `xml:"type,attr"`
	Domain   string `xml:"domain,attr"`
	Bus      string `xml:"bus,attr"`
	Slot     string `xml:"slot,attr"`
	Function string `xml:"function,attr"`
}

type DeviceGraphics struct {
	Type     string `xml:"type,attr"`
	Port     string `xml:"port,attr"`
	AutoPort string `xml:"autoport,attr"`
	Listen   string `xml:"listen,attr"`
	Keymap   string `xml:"keymap,attr"`
	// +child <listen type='address' address='0.0.0.0'/>
}

type Domain struct {
	DomainType    string `xml:"type,attr"`
	Name          string `xml:"name"`
	UUID          string `xml:"uuid"`
	Memory        int    `xml:"memory"`        // attr unit=KiB
	CurrentMemory int    `xml:"currentMemory"` // attr unit=KiB
	Cores         int    `xml:"vcpu"`          // attr placement=static
	OS            OS     `xml:"os"`
	/*
	  <features>
	    <acpi/>
	  </features>
	*/
	CPU        CPU    `xml:"cpu"`
	Clock      Clock  `xml:"clock"`
	OnPowerOff string `xml:"on_poweroff"`
	OnReboot   string `xml:"on_reboot"`
	OnCrash    string `xml:"on_crash"`
	/*
	  <pm>
	    <suspend-to-mem enabled='no'/>
	    <suspend-to-disk enabled='no'/>
	  </pm>
	*/
	Devices Device `xml:"devices"`
	/*
	   <controller type='usb' index='0' model='ich9-ehci1'>
	     <address type='pci' domain='0x0000' bus='0x00' slot='0x04' function='0x7'/>
	   </controller>
	   <controller type='usb' index='0' model='ich9-uhci1'>
	     <master startport='0'/>
	     <address type='pci' domain='0x0000' bus='0x00' slot='0x04' function='0x0' multifunction='on'/>
	   </controller>
	   <controller type='usb' index='0' model='ich9-uhci2'>
	     <master startport='2'/>
	     <address type='pci' domain='0x0000' bus='0x00' slot='0x04' function='0x1'/>
	   </controller>
	   <controller type='usb' index='0' model='ich9-uhci3'>
	     <master startport='4'/>
	     <address type='pci' domain='0x0000' bus='0x00' slot='0x04' function='0x2'/>
	   </controller>
	   <controller type='pci' index='0' model='pci-root'/>
	   <controller type='ide' index='0'>
	     <address type='pci' domain='0x0000' bus='0x00' slot='0x01' function='0x1'/>
	   </controller>
	   <interface type='bridge'>
	     <mac address='52:54:00:e1:ec:16'/>
	     <source bridge='br0'/>
	     <model type='virtio'/>
	     <address type='pci' domain='0x0000' bus='0x00' slot='0x03' function='0x0'/>
	   </interface>
	   <serial type='pty'>
	     <target port='0'/>
	   </serial>
	   <console type='pty'>
	     <target type='serial' port='0'/>
	   </console>
	   <input type='tablet' bus='usb'>
	     <address type='usb' bus='0' port='1'/>
	   </input>
	   <input type='mouse' bus='ps2'/>
	   <input type='keyboard' bus='ps2'/>
	   <video>
	     <model type='cirrus' vram='16384' heads='1' primary='yes'/>
	     <address type='pci' domain='0x0000' bus='0x00' slot='0x02' function='0x0'/>
	   </video>
	   <memballoon model='virtio'>
	     <address type='pci' domain='0x0000' bus='0x00' slot='0x06' function='0x0'/>
	   </memballoon>
	*/
	State     string
	Autostart string
}
