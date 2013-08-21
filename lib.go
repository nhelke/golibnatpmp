// This package is not thread safe. All calls to it should happen from a single
// thread or goroutine or be synchronised by some other mechanism
package golibnatpmp

// #cgo CFLAGS: -DENABLE_STRNATPMPERR
// #include "natpmp.h"
// #include "waitAndReadNATPMPResp.h"
import "C"

import (
	"errors"
)

// Values to pass to sendnewportmappingrequest()
const (
	UDP protocol = 1
	TCP protocol = 2
)

type protocol int

type NATPMP struct {
	ctype *C.natpmp_t
}

// Initialize a natpmp_t object.
//
// The gateway is not detected automaticaly and the passed gateway address is
// used unless forcedgw == 0
func NewNATPMPWithGW(forcedgw uint32) (natpmp NATPMP, err error) {
	var force C.int = 1
	if forcedgw == 0 {
		force = 0
	}
	natpmp.ctype = new(C.natpmp_t)
	ret, err := C.initnatpmp(natpmp.ctype, force, C.in_addr_t(forcedgw))
	if err == nil && ret != 0 {
		err = strNATPMPerr(ret)
	}
	return
}

// Initialize a natpmp_t object.
//
// The gateway is detected automaticaly. Equivalent to `NewNATPMPWithGW(0)`
func NewNATPMP() (NATPMP, error) {
	return NewNATPMPWithGW(0)
}

// Close resources associated with a natpmp_t object
func (self NATPMP) Close() (err error) {
	ret, err := C.closenatpmp(self.ctype)
	if err == nil && ret != 0 {
		err = strNATPMPerr(ret)
	}
	return
}

// IP is only valid if Err is nil
type PublicAddressResponse struct {
	IP  uint32
	Err error
}

// Send a public address NAT-PMP request to the network gateway
func (self NATPMP) SendPublicAddressRequest() <-chan PublicAddressResponse {
	c := make(chan PublicAddressResponse, 1)
	ret, err := C.sendpublicaddressrequest(self.ctype)
	if err == nil && ret != 2 {
		err = strNATPMPerr(ret)
	}
	if err != nil {
		c <- PublicAddressResponse{
			Err: err,
		}
	} else {
		go func() {
			resp, err := self.waitForAndReadNATPMPResponse()
			c <- PublicAddressResponse{
				IP:  uint32(resp.addr),
				Err: err,
			}
		}()
	}
	return c
}

// This struct's fields may only be considered valid if Err is nil
type NewPortMappingResponse struct {
	PrivatePort       uint16
	MappedPublicPort  uint16
	LifetimeInSeconds uint32
	Err               error
}

// Send a new port mapping NAT-PMP request to the network gateway
//
// Arguments :
// protocol is either golibnatpmp.TCP or golibnatpmp.UDP,
// lifetime is in seconds.
//
// To remove a port mapping, set lifetime to zero.
//
// To remove all port mappings to the host, set lifetime and both ports
// to zero.
func (self NATPMP) SendNewPortMappingRequest(protocol protocol,
	privateport, publicport uint16,
	lifetimeInSeconds uint32) <-chan NewPortMappingResponse {
	c := make(chan NewPortMappingResponse, 1)
	ret, err := C.sendnewportmappingrequest(self.ctype, C.int(protocol),
		C.uint16_t(privateport), C.uint16_t(publicport),
		C.uint32_t(lifetimeInSeconds))
	if err == nil && ret != 12 {
		err = strNATPMPerr(ret)
	}
	if err != nil {
		c <- NewPortMappingResponse{
			Err: err,
		}
	} else {
		go func() {
			resp, err := self.waitForAndReadNATPMPResponse()
			c <- NewPortMappingResponse{
				PrivatePort:       uint16(resp.privateport),
				MappedPublicPort:  uint16(resp.mappedpublicport),
				LifetimeInSeconds: uint32(resp.lifetime),
				Err:               err,
			}
		}()
	}
	return c
}

// resp is guarenteed not to be nil (except if memory runs out) but its fields
// may only be considered valid if err is nil
func (self NATPMP) waitForAndReadNATPMPResponse() (resp *C.resp_s, err error) {
	resp = new(C.resp_s)
	ret, err := C.waitAndReadNATPMPResp(self.ctype, resp)
	if err == nil && ret != 0 {
		err = strNATPMPerr(ret)
	}
	return
}

func strNATPMPerr(ret C.int) error {
	return errors.New(C.GoString(C.strnatpmperr(ret)))
}
