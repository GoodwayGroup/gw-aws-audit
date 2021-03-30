package iam

import (
	"github.com/GoodwayGroup/gw-aws-audit/lib"
	"github.com/clok/kemba"
)

var (
	kiam     = kemba.New("gw-aws-audit:iam")
	kbud     = kiam.Extend("buildUserData")
	kbuduser = kbud.Extend("fullUser")
	session  = lib.Session{}
)
