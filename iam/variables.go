package iam

import (
	"github.com/GoodwayGroup/gw-aws-audit/lib"
	"github.com/clok/kemba"
)

var (
	kiam    = kemba.New("gw-aws-audit:iam")
	session = lib.Session{}
)
