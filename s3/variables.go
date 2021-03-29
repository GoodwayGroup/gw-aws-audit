package s3

import (
	"github.com/GoodwayGroup/gw-aws-audit/lib"
	"github.com/clok/kemba"
)

var (
	k        = kemba.New("gw-aws-audit:s3")
	kcbo     = k.Extend("ClearBucketObjects")
	khr      = kcbo.Extend("handleResponse")
	kact     = k.Extend("AddCostTag")
	ktag     = k.Extend("checkCostTag")
	kup      = k.Extend("updateTags")
	khtr     = k.Extend("handleGetTagsResponse")
	kproc    = kact.Extend("processBucket")
	kmetrics = k.Extend("GetBucketMetrics")
	kpbm     = kmetrics.Extend("processBucketMetrics")
	session  = lib.Session{}
)
