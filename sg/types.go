package sg

import (
	"fmt"
	"github.com/aws/aws-sdk-go/service/ec2"
	"net"
	"strings"
)

type portToIP struct {
	port string
	ip   string
}

type securityGroup struct {
	id       string
	name     string
	attached map[string]int
	rules    map[string][]*ec2.IpRange
}

func (s *securityGroup) getAttachmentsAsString() string {
	var attachments []string
	for t, cnt := range s.attached {
		attachments = append(attachments, fmt.Sprintf("%s: %d", t, cnt))
	}
	return strings.Join(attachments, " ")
}

type groupedSecurityGroups struct {
	alert    map[*securityGroup][]*portToIP
	warning  map[*securityGroup][]*portToIP
	unknown  map[*securityGroup][]*portToIP
	wideOpen map[*securityGroup][]*portToIP
}

func newGroupedSecurityGroups() *groupedSecurityGroups {
	var gsg groupedSecurityGroups
	gsg.alert = map[*securityGroup][]*portToIP{}
	gsg.warning = map[*securityGroup][]*portToIP{}
	gsg.unknown = map[*securityGroup][]*portToIP{}
	gsg.wideOpen = map[*securityGroup][]*portToIP{}
	return &gsg
}

func (csg *groupedSecurityGroups) addToAlert(sec *securityGroup, value *portToIP) {
	if _, ok := csg.alert[sec]; !ok {
		csg.alert[sec] = []*portToIP{value}
	} else {
		csg.alert[sec] = append(csg.alert[sec], value)
	}
}

func (csg *groupedSecurityGroups) addToWarning(sec *securityGroup, value *portToIP) {
	if _, ok := csg.warning[sec]; !ok {
		csg.warning[sec] = []*portToIP{value}
	} else {
		csg.warning[sec] = append(csg.warning[sec], value)
	}
}

func (csg *groupedSecurityGroups) addToUnknown(sec *securityGroup, value *portToIP) {
	if _, ok := csg.unknown[sec]; !ok {
		csg.unknown[sec] = []*portToIP{value}
	} else {
		csg.unknown[sec] = append(csg.unknown[sec], value)
	}
}

func (csg *groupedSecurityGroups) addToWideOpen(sec *securityGroup, value *portToIP) {
	if _, ok := csg.wideOpen[sec]; !ok {
		csg.wideOpen[sec] = []*portToIP{value}
	} else {
		csg.wideOpen[sec] = append(csg.wideOpen[sec], value)
	}
}

type groupedIPBlockRules struct {
	approved []*net.IPNet
	warning  []*net.IPNet
	alert    []*net.IPNet
}

func newGroupedIPBlockRules() *groupedIPBlockRules {
	var g groupedIPBlockRules
	g.approved = []*net.IPNet{}
	g.warning = []*net.IPNet{}
	g.alert = []*net.IPNet{}
	return &g
}

func (g *groupedIPBlockRules) addToApproved(cidr *net.IPNet) {
	g.approved = append(g.approved, cidr)
}

func (g *groupedIPBlockRules) addToWarning(cidr *net.IPNet) {
	g.warning = append(g.warning, cidr)
}

func (g *groupedIPBlockRules) addToAlert(cidr *net.IPNet) {
	g.alert = append(g.alert, cidr)
}
