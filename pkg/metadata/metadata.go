package metadata

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"
)

// ParseDatasetCommitMetadata converts a string->string map, in the
// Dotscience Run Dataset Commit Metadata format, into a
// DatasetCommitMetadata struct. Any unrecognised keys in the map are
// ignored.
func ParseDatasetCommitMetadata(input map[string]string) DatasetCommitMetadata {
	var r DatasetCommitMetadata

	commitType, ok := input["type"]
	if ok {
		switch commitType {
		case "dotscience.run-output.v1":
			r.WorkspaceDotID, ok = input["workspace"]
			r.OutputFiles = map[string][]string{}
			for k, v := range input {
				if strings.HasPrefix(k, "run.") {
					parts := strings.Split(k, ".")
					if len(parts) == 3 && parts[2] == "dataset-output-files" {
						dotId := parts[1]
						filenames := []string{}
						err := json.Unmarshal([]byte(v), &filenames)
						if err == nil {
							r.OutputFiles[dotId] = filenames
						} // else, ignore errors
					}
				}
			}
		default:
			// Leave r as the empty value
		}
	}
	return r
}

// FIXME: most of this function should be handled by json parsing.

// ParseCommitMetadata converts a string->string map, in the Dotscience
// Run Commit Metadata format, into a CommitMetadata struct. Any unrecognised
// keys in the map are ignored.
func ParseCommitMetadata(input map[string]string) CommitMetadata {
	// Get a simple string value, or def if missing
	get := func(key string, def string) string {
		val, ok := input[key]
		if !ok {
			return def
		} else {
			return val
		}
	}

	// Get a simple string value as a pointer to a string, or nil if missing
	getP := func(key string) *string {
		val, ok := input[key]
		if !ok {
			return nil
		} else {
			return &val
		}
	}

	getBool := func(key string, def bool) bool {
		v := get(key, "")
		if v == "" {
			return def
		} else {
			return v == "true"
		}
	}

	getMaybeBool := func(key string) MaybeBool {
		v := get(key, "")
		if v == "true" {
			return MaybeTrue
		} else if v == "false" {
			return MaybeFalse
		} else {
			return MaybeUnknown
		}
	}

	getRunAuthority := func(key string) RunAuthority {
		v := get(key, "")
		if v == "workload" {
			return RunAuthority_Workload
		} else if v == "derived" {
			return RunAuthority_Derived
		} else if v == "correction" {
			return RunAuthority_Correction
		} else {
			// Unknown, so call it a correction, it'll attract attention.
			return RunAuthority_Correction
		}
	}

	getInt64 := func(key string, def int64) int64 {
		v := get(key, "")
		val, err := strconv.ParseInt(v, 10, 64)
		if err == nil {
			return val
		} else {
			return def
		}
	}

	getFloat64 := func(key string, def float64) float64 {
		v := get(key, "")
		val, err := strconv.ParseFloat(v, 64)
		if err == nil {
			return val
		} else {
			return def
		}
	}

	getTime := func(key string) time.Time {
		v := get(key, "")
		if v != "" {
			t, err := time.Parse("20060102T150405.999999999", v)
			if err != nil {
				return time.Time{}
			} else {
				return t
			}
		} else {
			return time.Time{}
		}
	}

	// key="[...json list of strings...]"
	getStringSlice := func(key string) []string {
		v := get(key, "[]")
		result := []string{}
		err := json.Unmarshal([]byte(v), &result)
		if err == nil {
			return result
		} else {
			return []string{}
		}
	}

	// key="{...json map from string to string...}"
	getDirectStringMap := func(key string) map[string]string {
		v := get(key, "{}")
		result := map[string]string{}
		err := json.Unmarshal([]byte(v), &result)
		if err == nil {
			return result
		} else {
			return map[string]string{}
		}
	}

	// wantedKey<FOO>=<BAR>   -> map with entry of "FOO": "BAR"
	getStringMap := func(wantedKey string) map[string]string {
		parsedResult := map[string]string{}
		for key, val := range input {
			// key strings are {WANTEDKEY}{NAME}
			if strings.HasPrefix(key, wantedKey) {
				name := key[len(wantedKey):]

				parsedResult[name] = val
			}
		}
		return parsedResult
	}

	// wantedKey<FOO>="[<BAR>, ...]"   -> map with entry of "FOO": ["BAR", ...]
	getStringSliceMap := func(wantedKey string) map[string][]string {
		parsedResult := map[string][]string{}
		for key, val := range input {
			// key strings are {WANTEDKEY}{NAME}
			if strings.HasPrefix(key, wantedKey) {
				name := key[len(wantedKey):]

				result := []string{}
				err := json.Unmarshal([]byte(val), &result)
				if err == nil {
					parsedResult[name] = result
				}
			}
		}
		return parsedResult
	}

	// wantedKey<FOO>=<DOT>@<VERSION> -> map with entry of "FOO": DatasetVersion{DOT, VERSION}
	getDSVMap := func(wantedKey string) map[string]DatasetVersion {
		parsedResult := map[string]DatasetVersion{}
		for key, dsv := range input {
			// key strings are {WANTEDKEY}{NAME}
			if strings.HasPrefix(key, wantedKey) {
				name := key[len(wantedKey):]

				// dsv strings are DS@VERSION
				parts := strings.Split(dsv, "@")
				if len(parts) == 2 {
					parsedResult[name] = DatasetVersion{ID: DotID(parts[0]), Version: parts[1]}
				}
			}
		}
		return parsedResult
	}

	// wantedKey=JSON LIST OF DOT@VERSION STRINGS
	// -> []InputFile{DOT, VERSION}
	getIFs := func(wantedKey string) []InputFile {
		infs := get(wantedKey, "[]")
		result := []string{}
		err := json.Unmarshal([]byte(infs), &result)

		parsedInfs := make([]InputFile, len(result))

		if err == nil {
			// inf strings are FILE@VERSION
			for idx, inf := range result {
				parts := strings.Split(inf, "@")
				if len(parts) == 2 {
					parsedInfs[idx] = InputFile{Filename: parts[0], Version: parts[1]}
				}
			}
		}

		return parsedInfs
	}

	// wantedKey<REF> = JSON LIST OF DOT@VERSION STRINGS
	// -> map with entry of "REF": []InputFile{DOT, VERSION}
	getIFMap := func(wantedKey string) map[string][]InputFile {
		parsedResult := map[string][]InputFile{}
		for key, _ := range input {
			// key strings are {WANTEDKEY}{NAME}
			if strings.HasPrefix(key, wantedKey) {
				name := key[len(wantedKey):]
				parsedResult[name] = getIFs(key)
			}
		}
		return parsedResult
	}

	getRuns := func(runIds []string) []RunMetadata {
		runs := make([]RunMetadata, len(runIds))
		for idx, runId := range runIds {
			prefix := fmt.Sprintf("run.%s.", runId)
			runs[idx] = RunMetadata{
				RunID:                runId,
				Authority:            getRunAuthority(prefix + "authority"),
				Description:          get(prefix+"description", ""),
				WorkloadFile:         get(prefix+"workload-file", ""),
				Success:              getP(prefix+"error") == nil,
				ErrorMessage:         getP(prefix + "error"),
				Labels:               getStringMap(prefix + "label."),
				Summary:              getStringMap(prefix + "summary."),
				Parameters:           getStringMap(prefix + "parameters."),
				ExecStart:            getTime(prefix + "start"),
				ExecEnd:              getTime(prefix + "end"),
				WorkspaceInputFiles:  getIFs(prefix + "input-files"),
				WorkspaceOutputFiles: getStringSlice(prefix + "output-files"),
				DatasetInputFiles:    getIFMap(prefix + "dataset-input-files."),
				DatasetOutputFiles:   getStringSliceMap(prefix + "dataset-output-files."),
			}
		}
		return runs
	}

	var r CommitMetadata

	_, hasRuns := input["runs"]

	if !hasRuns {
		r = CommitMetadata{
			Success: false,
			Message: "No run metadata was returned",
		}
	} else {
		r = CommitMetadata{
			SubmitterID: get("author", ""),

			Inputs:  getDSVMap("input-dataset."),
			Outputs: getDSVMap("output-dataset."),

			WorkloadType:        get("workload.type", ""),
			WorkloadImage:       get("workload.image", ""),
			WorkloadImageHash:   get("workload.image.hash", ""),
			WorkloadCommand:     getStringSlice("workload.command"),
			WorkloadEnvironment: getDirectStringMap("workload.environment"),

			ExecStart:          getTime("exec.start"),
			ExecEnd:            getTime("exec.end"),
			ExecLogs:           getStringSlice("exec.logs"),
			ExecCPUSecondsUsed: getFloat64("exec.cpu-seconds", -1.0),
			ExecPeakRAMBytes:   getInt64("exec.ram", -1),

			RunnerName:            get("runner.name", ""),
			RunnerVersion:         get("runner.version", ""),
			RunnerPlatform:        get("runner.platform", ""),
			RunnerPlatformVersion: get("runner.platform_version", ""),
			RunnerCPUs:            getStringSlice("runner.cpu"),
			RunnerGPUs:            getStringSlice("runner.gpu"),
			RunnerRAMBytes:        getInt64("runner.ram", -1),
			RunnerRAMECC:          getMaybeBool("runner.ram.ecc"),

			Runs: getRuns(getStringSlice("runs")),

			Success: getBool("success", true),
			Message: get("message", ""),
		}
	}

	return r
}
