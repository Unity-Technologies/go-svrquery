package titanfall

import (
	"github.com/multiplay/go-svrquery/lib/svrquery/common"
	"github.com/netdata/go-orchestrator/module"
)

var (
	charts = module.Charts{
		{
			ID:    "%d_frame_time",
			Title: "Frame Time",
			Ctx:   "clanforge.frame_time",
			Units: "milliseconds",
			Fam:   "serverid %d",
			Type:  module.Area,
			Dims: module.Dims{
				{ID: "%d_avg_frame_time", Name: "average", Div: common.Dim3DP},
				{ID: "%d_max_frame_time", Name: "max", Div: common.Dim3DP},
			},
		},
		{
			ID:    "%d_user_cmd_time",
			Title: "User Command Time",
			Ctx:   "clanforge.user_cmd_time",
			Units: "milliseconds",
			Fam:   "serverid %d",
			Type:  module.Area,
			Dims: module.Dims{
				{ID: "%d_avg_user_cmd_time", Name: "average", Div: common.Dim3DP},
				{ID: "%d_max_user_cmd_time", Name: "max", Div: common.Dim3DP},
			},
		},
		{
			ID:    "%d_phase",
			Title: "Phase",
			Ctx:   "clanforge.phase",
			Units: "Phase",
			Fam:   "serverid %d",
			Type:  module.Line,
			Dims: module.Dims{
				{ID: "%d_phase", Name: "phase"},
			},
		},
	}
)
