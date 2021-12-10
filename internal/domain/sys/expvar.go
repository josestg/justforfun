package sys

import "expvar"

// BuildXXX are exported variables.
// These variables will be sets at build time.
var (
	BuildRef  = expvar.NewString("build_ref")
	BuildDate = expvar.NewString("build_date")
	BuildName = expvar.NewString("build_name")
)

// Info represents the system info.
type Info struct {
	InstanceID   string `json:"instance_id"`
	InstanceName string `json:"instance_name"`
	BuildRef     string `json:"build_ref"`
	BuildDate    string `json:"build_date"`
	Environment  string `json:"environment"`
}

// NewInfo creates a new system info with build info.
func NewInfo(id string, env string) Info {
	return Info{
		InstanceID:   id,
		Environment:  env,
		InstanceName: BuildName.Value(),
		BuildRef:     BuildRef.Value(),
		BuildDate:    BuildDate.Value(),
	}
}

// Support represents the support system info.
type Support struct {
	Database string `json:"database"`
	Storage  string `json:"storage"`
	Service  string `json:"service"`
}
