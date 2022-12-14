package probers

import (
	pb "250-monitor/pinger"
	"fmt"
	"net"
)

type ICMPProbe struct {
	target  string
	replyCh chan *pb.PingReply
}

func NewICMPProbe(target string, replyCh chan *pb.PingReply) *ICMPProbe {
	icmpP := new(ICMPProbe)
	icmpP.target = target
	icmpP.replyCh = replyCh
	return icmpP
}

func (icmpP *ICMPProbe) Start() {
	if net.ParseIP(icmpP.target) == nil {
		fmt.Printf("An invalid IP (%s) was provided to ICMP Probe, aborting\n", icmpP.target)
		return
	}
	go icmpP.ping()
}

func (icmpP *ICMPProbe) ping() {
	ping(icmpP.target, icmpP.replyCh)
}
