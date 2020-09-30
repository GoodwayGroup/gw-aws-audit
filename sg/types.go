package sg

import (
	"fmt"
	"github.com/aws/aws-sdk-go/service/ec2"
	"net"
	"strings"
)

type portToIP struct {
	port   string
	ip     string
	proto  string
	prefix *Prefix
}

// SecurityGroup defines the struct for common SG properties used by this tool.
type SecurityGroup struct {
	id       string
	name     string
	attached map[string]int
	rules    map[string][]*ec2.IpRange
}

// GetAttachmentsAsString will return a formatted list of AWS attachments
func (s *SecurityGroup) GetAttachmentsAsString() string {
	var attachments []string
	for t, cnt := range s.Attachments() {
		attachments = append(attachments, fmt.Sprintf("%s: %d", t, cnt))
	}
	return strings.Join(attachments, " ")
}

// Attachments will return the map of Attachments
func (s *SecurityGroup) Attachments() map[string]int {
	return s.attached
}

// ID will return the SecurityGroup ID
func (s *SecurityGroup) ID() string {
	return s.id
}

// Name will return the SecurityGroup Name
func (s *SecurityGroup) Name() string {
	return s.name
}

// Rules will return the SecurityGroup Rules map
func (s *SecurityGroup) Rules() map[string][]*ec2.IpRange {
	return s.rules
}

// ParseRuleToken break the Rules token key from the Rules map and return
// the component parts of [port, protocol, security group IDs]
func (s SecurityGroup) ParseRuleToken(token string) (port string, protocol string, sgIDs string) {
	parts := strings.Split(token, "::")
	return parts[0], parts[1], parts[2]
}

type groupedSecurityGroups struct {
	alert    map[*SecurityGroup][]*portToIP
	warning  map[*SecurityGroup][]*portToIP
	unknown  map[*SecurityGroup][]*portToIP
	wideOpen map[*SecurityGroup][]*portToIP
	amazon   map[*SecurityGroup][]*portToIP
	sg       map[*SecurityGroup][]*portToIP
}

func newGroupedSecurityGroups() *groupedSecurityGroups {
	var gsg groupedSecurityGroups
	gsg.alert = map[*SecurityGroup][]*portToIP{}
	gsg.warning = map[*SecurityGroup][]*portToIP{}
	gsg.unknown = map[*SecurityGroup][]*portToIP{}
	gsg.wideOpen = map[*SecurityGroup][]*portToIP{}
	gsg.amazon = map[*SecurityGroup][]*portToIP{}
	gsg.sg = map[*SecurityGroup][]*portToIP{}
	return &gsg
}

func (csg *groupedSecurityGroups) addToAlert(sec *SecurityGroup, value *portToIP) {
	if _, ok := csg.alert[sec]; !ok {
		csg.alert[sec] = []*portToIP{value}
	} else {
		csg.alert[sec] = append(csg.alert[sec], value)
	}
}

func (csg *groupedSecurityGroups) addToWarning(sec *SecurityGroup, value *portToIP) {
	if _, ok := csg.warning[sec]; !ok {
		csg.warning[sec] = []*portToIP{value}
	} else {
		csg.warning[sec] = append(csg.warning[sec], value)
	}
}

func (csg *groupedSecurityGroups) addToUnknown(sec *SecurityGroup, value *portToIP) {
	if _, ok := csg.unknown[sec]; !ok {
		csg.unknown[sec] = []*portToIP{value}
	} else {
		csg.unknown[sec] = append(csg.unknown[sec], value)
	}
}

func (csg *groupedSecurityGroups) addToWideOpen(sec *SecurityGroup, value *portToIP) {
	if _, ok := csg.wideOpen[sec]; !ok {
		csg.wideOpen[sec] = []*portToIP{value}
	} else {
		csg.wideOpen[sec] = append(csg.wideOpen[sec], value)
	}
}

func (csg *groupedSecurityGroups) addToAmazon(sec *SecurityGroup, value *portToIP) {
	if _, ok := csg.amazon[sec]; !ok {
		csg.amazon[sec] = []*portToIP{value}
	} else {
		csg.amazon[sec] = append(csg.amazon[sec], value)
	}
}

func (csg *groupedSecurityGroups) addToSG(sec *SecurityGroup, value *portToIP) {
	if _, ok := csg.sg[sec]; !ok {
		csg.sg[sec] = []*portToIP{value}
	} else {
		csg.sg[sec] = append(csg.sg[sec], value)
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

// AWSIPRanges is the JSON struct used to parse the AWS IP Range file.
type AWSIPRanges struct {
	SyncToken    string        `json:"syncToken"`
	CreateDate   string        `json:"createDate"`
	Prefixes     []*Prefix     `json:"prefixes"`
	IPv6Prefixes []*IPv6Prefix `json:"ipv6_prefixes"`
}

// Prefix is used with AWSIPRanges.
type Prefix struct {
	IPPrefix           string `json:"ip_prefix"`
	Region             string `json:"region"`
	NetworkBorderGroup string `json:"network_border_group"`
	Service            string `json:"service"`
}

// GetCIDR will extract the IPPrefix as a CIDR definition.
func (p *Prefix) GetCIDR() (*net.IPNet, error) {
	_, ipv4Net, err := net.ParseCIDR(p.IPPrefix)
	if err != nil {
		return nil, err
	}
	return ipv4Net, err
}

// GetService will extract the AWS service name that the IP is associated with.
func (p *Prefix) GetService() string {
	switch p.Service {
	case "AMAZON":
		return "EC2"
	default:
		return p.Service
	}
}

// IPv6Prefix is used with AWSIPRanges.
type IPv6Prefix struct {
	IPv6Prefix         string `json:"ipv6_prefix"`
	Region             string `json:"region"`
	NetworkBorderGroup string `json:"network_border_group"`
	Service            string `json:"service"`
}

// GetCIDR will extract the IPPrefix as a CIDR definition.
func (p *IPv6Prefix) GetCIDR() (*net.IPNet, error) {
	_, ipv6Net, err := net.ParseCIDR(p.IPv6Prefix)
	if err != nil {
		return nil, err
	}
	return ipv6Net, err
}

// AWSIPs is a
type AWSIPs struct {
	list  []*net.IPNet
	table map[*net.IPNet]*Prefix
}
