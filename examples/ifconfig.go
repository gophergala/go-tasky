package main

import (
	"encoding/json"
	"net"

	"github.com/gophergala/go-tasky/tasky"
)

type Ifconfig struct {
	Name         string
	HardwareAddr string
	Ipnet        string
}

func (d *Ifconfig) Info() []byte {
	s := `{
	}`

	return []byte(s)
}

func (d *Ifconfig) Services() []byte {
	return nil
}

type ifErr struct {
	Error string
}

type ifs struct {
	Interfaces []Ifconfig
}

func (d *Ifconfig) Perform(job []byte) ([]byte, bool) {
	b := true
	aflag := &b

	i := ifs{}

	ifaces, err := net.Interfaces()
	if err != nil {
		e := ifErr{err.Error()}
		estr, _ := json.Marshal(e)
		return estr, false
	}

	for _, iface := range ifaces {
		addr, err := iface.Addrs()
		if err != nil {
			e := ifErr{err.Error()}
			estr, _ := json.Marshal(e)
			return estr, false
		}
		for _, ifaddr := range addr {
			ipnet, ok := ifaddr.(*net.IPNet)
			if !ok {
				continue
			}
			v4 := ipnet.IP.To4()
			if iface.Name[:2] == "lo" || v4 == nil {
				if *aflag == false {
					continue
				}
			}

			ii := Ifconfig{}
			ii.Name = iface.Name
			ii.HardwareAddr = iface.HardwareAddr.String()
			ii.Ipnet = ipnet.String()

			if len(i.Interfaces) <= 0 {
				i.Interfaces = make([]Ifconfig, 1)
				i.Interfaces[0] = ii
			} else {
				i.Interfaces = append(i.Interfaces, ii)
			}
		}
	}

	jsonStr, err := json.Marshal(&i)
	if err != nil {
		e := ifErr{err.Error()}
		estr, _ := json.Marshal(e)
		return estr, false
	}

	return jsonStr, true
}

func (d *Ifconfig) Status() []byte {
	return nil
}

func (d *Ifconfig) Signal(act tasky.Action) bool {
	return true
}

func (d *Ifconfig) Statistics() []byte {
	return nil
}
