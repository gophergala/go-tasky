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

type ifs struct {
	Interfaces []Ifconfig
}

func (d *Ifconfig) Perform(job []byte, dataCh chan []byte, errCh chan error, quitCh chan bool) {
	done := make(chan bool)

	go func() {
		b := true

		i := ifs{}

		ifaces, err := net.Interfaces()
		if err != nil {
			errCh <- err
			done <- true
			return
		}

		for _, iface := range ifaces {
			addr, err := iface.Addrs()
			if err != nil {
				errCh <- err
				done <- true
				return
			}
			for _, ifaddr := range addr {
				ipnet, ok := ifaddr.(*net.IPNet)
				if !ok {
					continue
				}
				v4 := ipnet.IP.To4()
				if iface.Name[:2] == "lo" || v4 == nil {
					if b == false {
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
			errCh <- err
			done <- true
			return
		}

		dataCh <- jsonStr
		done <- true
	}()

	select {
	case <-done:
		return

	case <-quitCh:
		return
	}
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