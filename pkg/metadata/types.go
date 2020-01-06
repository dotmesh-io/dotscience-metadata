package metadata

import (
	"time"
)

// type DotMode represents an attachment mode when using a Dot as a
// Dataset. It can be DotMode_Input for an input dataset,
// DotMode_Output for an output dataset, or DotMode_ReadWrite for a
// dataset that is used for both input and output.
type DotMode int

const (
	DotMode_Input     DotMode = 1
	DotMode_Output    DotMode = 2
	DotMode_ReadWrite DotMode = 3
)

// type DotID represents the internal ID of a Dot.
type DotID string

// type DatasetVersion is a record of a specific version of a specific
// dot being used as a dataset. This type is used when recording the
// history of a run; the similar type Dataset is used to specify the
// desire to use a Dot as a Dataset for a run.
type DatasetVersion struct {
	ID      DotID  `json:"id"`
	Version string `json:"version"`
}

// type MaybeBool represents a trinary logic value - True, False, or Unknown.
type MaybeBool int

const (
	MaybeUnknown MaybeBool = iota
	MaybeTrue
	MaybeFalse
)

// type DatasetCommitMetadata records the results of a dataset commit.
type DatasetCommitMetadata struct {
	// The ID of the workspace dot, in which the run IDs in OutputFiles may be found.
	WorkspaceDotID string

	// OutputFiles is a map from run IDs to arrays of filenames modified by that run in this dataset.
	OutputFiles map[string][]string
}

// type CommitMetadata records the results of a commit, containing one or more runs.
type CommitMetadata struct {
	SubmitterID string `json:"submitter_id"`

	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`

	WorkloadType        string                    `json:"workload_type,omitempty"`
	WorkloadImage       string                    `json:"workload_image,omitempty"`
	WorkloadImageHash   string                    `json:"workload_image_hash,omitempty"`
	WorkloadCommand     []string                  `json:"workload_command,omitempty"`
	WorkloadEnvironment map[string]string         `json:"workload_environment,omitempty"`
	Inputs              map[string]DatasetVersion `json:"inputs,omitempty"`
	Outputs             map[string]DatasetVersion `json:"outputs,omitempty"`

	ExecLogs           []string  `json:"exec_logs"`
	ExecStart          time.Time `json:"exec_start"`
	ExecEnd            time.Time `json:"exec_end"`
	ExecCPUSecondsUsed float64   `json:"exec_cpu_seconds_used,omitempty"`
	ExecPeakRAMBytes   int64     `json:"exec_peak_ram_bytes,omitempty"`

	RunnerName            string    `json:"runner_name,omitempty"`
	RunnerVersion         string    `json:"runner_version,omitempty"`
	RunnerPlatform        string    `json:"runner_platform,omitempty"`
	RunnerPlatformVersion string    `json:"runner_platform_version,omitempty"`
	RunnerCPUs            []string  `json:"runner_cpus,omitempty"`
	RunnerGPUs            []string  `json:"runner_gpus,omitempty"`
	RunnerRAMBytes        int64     `json:"runner_ram_bytes,omitempty"`
	RunnerRAMECC          MaybeBool `json:"runner_ram_ecc,omitempty"`

	Runs []RunMetadata `json:"runs"`
}

type RunAuthority int

const (
	RunAuthority_Workload RunAuthority = iota
	RunAuthority_Derived
	RunAuthority_Correction
)

// type InputFile records the version of a file used as input
type InputFile struct {
	Filename string `json:"filename"`
	Version  string `json:"version"`
}

// type RunMetadata records the final result of a run.
type RunMetadata struct {
	RunID     string       `json:"run_id"`
	CommitID  string       `json:"commit_id"`
	Authority RunAuthority `json:"authority"`

	Description  string `json:"description,omitempty"`
	WorkloadFile string `json:"workload_file,omitempty"`

	Success      bool    `json:"success"`
	ErrorMessage *string `json:"error_message,omitempty"`

	WorkspaceInputFiles  []InputFile `json:"workspace_input_files,omitempty"`
	WorkspaceOutputFiles []string    `json:"workspace_output_files,omitempty"`

	DatasetInputFiles  map[string][]InputFile `json:"dataset_input_files,omitempty"`
	DatasetOutputFiles map[string][]string    `json:"dataset_output_files,omitempty"`

	Labels     map[string]string `json:"labels,omitempty"`
	Summary    map[string]string `json:"summary,omitempty"`
	Parameters map[string]string `json:"parameters,omitempty"`

	ExecStart time.Time `json:"exec_start,omitempty"`
	ExecEnd   time.Time `json:"exec_end,omitempty"`
}
