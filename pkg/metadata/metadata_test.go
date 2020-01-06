package metadata

import (
	"testing"
	"time"
)

func testEqStr(t *testing.T, got, expected string) {
	t.Helper()
	if got != expected {
		t.Errorf("Wanted %s, got %s", expected, got)
	}
}

func testEqTime(t *testing.T, got, expected time.Time) {
	t.Helper()
	if got != expected {
		t.Errorf("Wanted %v, got %v", expected, got)
	}
}

func testEqStrs(t *testing.T, got, expected []string) {
	t.Helper()
	if len(got) != len(expected) {
		t.Errorf("Wanted %#v, got %#v", expected, got)
		return
	}

	// Same number of vals in both
	for idx, v := range expected {
		if got[idx] != v {
			t.Errorf("Wanted %#v, got %#v", expected, got)
			return
		}
	}
}

func testEqMap(t *testing.T, got, expected map[string]string) {
	t.Helper()
	if len(got) != len(expected) {
		t.Errorf("Wanted %#v, got %#v", expected, got)
		return
	}

	// Same number of keys in both, so enumerating the keys of either will find any differences
	for k, v := range expected {
		gv, ok := got[k]
		if !ok || (gv != v) {
			t.Errorf("Wanted %#v, got %#v", expected, got)
			return
		}
	}
}

func testEqDsvs(t *testing.T, got, expected map[string]DatasetVersion) {
	t.Helper()
	if len(got) != len(expected) {
		t.Errorf("Wanted %#v, got %#v", expected, got)
		return
	}

	// Same number of keys in both, so enumerating the keys of either will find any differences
	for k, v := range expected {
		gv, ok := got[k]
		if !ok || (gv != v) {
			t.Errorf("Wanted %#v, got %#v", expected, got)
			return
		}
	}
}

func testEqIFs(t *testing.T, got, expected []InputFile) {
	t.Helper()
	if len(got) != len(expected) {
		t.Errorf("Wanted %#v, got %#v", expected, got)
		return
	}

	for idx, v := range expected {
		gv := got[idx]
		if gv != v {
			t.Errorf("Wanted %#v, got %#v", expected, got)
			return
		}
	}
}

func TestParseCommitMetadataThorough(t *testing.T) {
	rm := ParseCommitMetadata(map[string]string{
		"type":                    "dotscience.run.v1",
		"author":                  "452342",
		"date":                    "1538658370073482093",
		"workload.type":           "command",
		"workload.image":          "busybox",
		"workload.image.hash":     "busybox@sha256:2a03a6059f21e150ae84b0973863609494aad70f0a80eaeb64bddd8d92465812",
		"workload.command":        "[\"sh\",\"-c\",\"curl http://localhost/testjob.sh | /bin/sh\"]",
		"workload.environment":    "{\"DEBUG_MODE\": \"YES\"}",
		"runner.version":          "Runner=Dotscience Docker Executor rev. 63db3d0 Agent=Dotscience Agent rev. b1acc85",
		"runner.name":             "bob",
		"runner.platform":         "linux",
		"runner.platform_version": "Linux a1bc10a2fb6e 4.14.60 #1-NixOS SMP Fri Aug 3 05:50:45 UTC 2018 x86_64 GNU/Linux",
		"runner.ram":              "16579702784",
		"runner.cpu":              "[\"Intel(R) Core(TM) i7-7500U CPU @ 2.70GHz\", \"Intel(R) Core(TM) i7-7500U CPU @ 2.70GHz\", \"Intel(R) Core(TM) i7-7500U CPU @ 2.70GHz\", \"Intel(R) Core(TM) i7-7500U CPU @ 2.70GHz\"]",
		"exec.start":              "20181004T130607.101",
		"exec.end":                "20181004T130610.223",
		"exec.logs":               "[\"16204868-ae5a-4574-907b-8d4774aad497/agent-stdout.log\",\"16204868-ae5a-4574-907b-8d4774aad497/pull-workload-stdout.log\",\"16204868-ae5a-4574-907b-8d4774aad497/workload-stdout.log\"]",
		"input-dataset.b":         "<ID of dot B>@<commit ID of dot B before the run>",
		"input-dataset.c":         "<ID of dot C>@<commit ID of dot C before the run>",
		"output-dataset.c":        "<ID of dot C>@<commit ID of dot C created by this run>",
		"output-dataset.d":        "<ID of dot D>@<commit ID of dot D created by this run>",
		"runs":                    "[\"02ecdc67-c49e-4d76-abe8-1ee13f2884b7\", \"cd351be8-3ba9-4c5e-ad26-429d6d6033de\", \"31df506d-c715-4159-99fd-60bb845d4dec\"]",
		"run.02ecdc67-c49e-4d76-abe8-1ee13f2884b7.authority":              "workload",
		"run.02ecdc67-c49e-4d76-abe8-1ee13f2884b7.input-files":            "[\"foo.csv@<some earlier commit ID of workspace dot>\"]",
		"run.02ecdc67-c49e-4d76-abe8-1ee13f2884b7.dataset-input-files.b":  "[\"input.csv@<some earlier commit ID of b>\"]",
		"run.02ecdc67-c49e-4d76-abe8-1ee13f2884b7.dataset-input-files.c":  "[\"cache.sqlite@<some earlier commit ID of c>\"]",
		"run.02ecdc67-c49e-4d76-abe8-1ee13f2884b7.output-files":           "[\"log.txt\"]",
		"run.02ecdc67-c49e-4d76-abe8-1ee13f2884b7.dataset-output-files.c": "[\"cache.sqlite\"]",
		"run.02ecdc67-c49e-4d76-abe8-1ee13f2884b7.dataset-output-files.d": "[\"output.csv\"]",
		"run.02ecdc67-c49e-4d76-abe8-1ee13f2884b7.summary.rms_error":      "0.057",
		"run.02ecdc67-c49e-4d76-abe8-1ee13f2884b7.parameters.smoothing":   "1.0",
		"run.02ecdc67-c49e-4d76-abe8-1ee13f2884b7.start":                  "20181004T130607.225",
		"run.02ecdc67-c49e-4d76-abe8-1ee13f2884b7.end":                    "20181004T130608.225",
		"run.cd351be8-3ba9-4c5e-ad26-429d6d6033de.authority":              "workload",
		"run.cd351be8-3ba9-4c5e-ad26-429d6d6033de.input-files":            "[\"foo.csv@<some earlier commit ID of workspace dot>\"]",
		"run.cd351be8-3ba9-4c5e-ad26-429d6d6033de.dataset-input-files.b":  "[\"input.csv@<some earlier commit ID of b>\"]",
		"run.cd351be8-3ba9-4c5e-ad26-429d6d6033de.dataset-input-files.c":  "[\"cache.sqlite@<some earlier commit ID of c>\"]",
		"run.cd351be8-3ba9-4c5e-ad26-429d6d6033de.output-files":           "[\"log.txt\"]",
		"run.cd351be8-3ba9-4c5e-ad26-429d6d6033de.dataset-output-files.c": "[\"cache.sqlite\"]",
		"run.cd351be8-3ba9-4c5e-ad26-429d6d6033de.dataset-output-files.d": "[\"output.csv\"]",
		"run.cd351be8-3ba9-4c5e-ad26-429d6d6033de.summary.rms_error":      "0.123",
		"run.cd351be8-3ba9-4c5e-ad26-429d6d6033de.parameters.smoothing":   "2",
		"run.cd351be8-3ba9-4c5e-ad26-429d6d6033de.start":                  "20181004T130608.579",
		"run.cd351be8-3ba9-4c5e-ad26-429d6d6033de.end":                    "20181004T130609.579",
		"run.31df506d-c715-4159-99fd-60bb845d4dec.authority":              "correction",
		"run.31df506d-c715-4159-99fd-60bb845d4dec.description":            "File changes were detected that the run metadata did not explain",
		"run.31df506d-c715-4159-99fd-60bb845d4dec.output-files":           "[\"mylibrary.pyc\"]",
	})

	if rm.Success != true {
		t.Errorf("Wanted %t, got %t", true, rm.Success)
	}
	testEqStr(t, rm.SubmitterID, "452342")
	testEqStrs(t, rm.WorkloadCommand, []string{"sh", "-c", "curl http://localhost/testjob.sh | /bin/sh"})
	testEqStr(t, rm.WorkloadType, "command")
	testEqStr(t, rm.WorkloadImage, "busybox")
	testEqStr(t, rm.WorkloadImageHash, "busybox@sha256:2a03a6059f21e150ae84b0973863609494aad70f0a80eaeb64bddd8d92465812")
	testEqMap(t, rm.WorkloadEnvironment, map[string]string{"DEBUG_MODE": "YES"})
	testEqStr(t, rm.RunnerVersion, "Runner=Dotscience Docker Executor rev. 63db3d0 Agent=Dotscience Agent rev. b1acc85")
	testEqStr(t, rm.RunnerName, "bob")
	testEqStr(t, rm.RunnerPlatform, "linux")
	testEqStr(t, rm.RunnerPlatformVersion, "Linux a1bc10a2fb6e 4.14.60 #1-NixOS SMP Fri Aug 3 05:50:45 UTC 2018 x86_64 GNU/Linux")
	if rm.RunnerRAMBytes != 16579702784 {
		t.Errorf("Wanted %d, got %d", 16579702784, rm.RunnerRAMBytes)
	}
	if rm.RunnerRAMECC != MaybeUnknown {
		t.Errorf("Wanted Unknown, got %v", rm.RunnerRAMECC)
	}
	testEqStrs(t, rm.RunnerCPUs, []string{
		"Intel(R) Core(TM) i7-7500U CPU @ 2.70GHz",
		"Intel(R) Core(TM) i7-7500U CPU @ 2.70GHz",
		"Intel(R) Core(TM) i7-7500U CPU @ 2.70GHz",
		"Intel(R) Core(TM) i7-7500U CPU @ 2.70GHz",
	})
	testEqTime(t, rm.ExecStart, time.Date(2018, 10, 4, 13, 6, 7, 101000000, time.UTC))
	testEqTime(t, rm.ExecEnd, time.Date(2018, 10, 4, 13, 6, 10, 223000000, time.UTC))
	testEqStrs(t, rm.ExecLogs, []string{"16204868-ae5a-4574-907b-8d4774aad497/agent-stdout.log", "16204868-ae5a-4574-907b-8d4774aad497/pull-workload-stdout.log", "16204868-ae5a-4574-907b-8d4774aad497/workload-stdout.log"})
	testEqDsvs(t, rm.Inputs, map[string]DatasetVersion{
		"b": DatasetVersion{ID: "<ID of dot B>", Version: "<commit ID of dot B before the run>"},
		"c": DatasetVersion{ID: "<ID of dot C>", Version: "<commit ID of dot C before the run>"},
	})
	testEqDsvs(t, rm.Outputs, map[string]DatasetVersion{
		"c": DatasetVersion{ID: "<ID of dot C>", Version: "<commit ID of dot C created by this run>"},
		"d": DatasetVersion{ID: "<ID of dot D>", Version: "<commit ID of dot D created by this run>"},
	})

	if len(rm.Runs) != 3 {
		t.Errorf("Expected 3 runs, got %d", len(rm.Runs))
	} else {
		testEqStr(t, rm.Runs[0].RunID, "02ecdc67-c49e-4d76-abe8-1ee13f2884b7")
		if rm.Runs[0].Authority != RunAuthority_Workload {
			t.Errorf("Expected authority %d, got %d", RunAuthority_Workload, rm.Runs[0].Authority)
		}
		testEqIFs(t, rm.Runs[0].WorkspaceInputFiles, []InputFile{InputFile{Filename: "foo.csv", Version: "<some earlier commit ID of workspace dot>"}})
		testEqIFs(t, rm.Runs[0].DatasetInputFiles["b"], []InputFile{InputFile{Filename: "input.csv", Version: "<some earlier commit ID of b>"}})
		testEqIFs(t, rm.Runs[0].DatasetInputFiles["c"], []InputFile{InputFile{Filename: "cache.sqlite", Version: "<some earlier commit ID of c>"}})
		testEqStrs(t, rm.Runs[0].WorkspaceOutputFiles, []string{"log.txt"})
		testEqStrs(t, rm.Runs[0].DatasetOutputFiles["c"], []string{"cache.sqlite"})
		testEqStrs(t, rm.Runs[0].DatasetOutputFiles["d"], []string{"output.csv"})
		testEqMap(t, rm.Runs[0].Labels, map[string]string{})
		testEqMap(t, rm.Runs[0].Summary, map[string]string{"rms_error": "0.057"})
		testEqMap(t, rm.Runs[0].Parameters, map[string]string{"smoothing": "1.0"})
		testEqTime(t, rm.Runs[0].ExecStart, time.Date(2018, 10, 4, 13, 6, 7, 225000000, time.UTC))
		testEqTime(t, rm.Runs[0].ExecEnd, time.Date(2018, 10, 4, 13, 6, 8, 225000000, time.UTC))

		testEqStr(t, rm.Runs[1].RunID, "cd351be8-3ba9-4c5e-ad26-429d6d6033de")
		if rm.Runs[1].Authority != RunAuthority_Workload {
			t.Errorf("Expected authority %d, got %d", RunAuthority_Workload, rm.Runs[1].Authority)
		}
		testEqIFs(t, rm.Runs[1].WorkspaceInputFiles, []InputFile{InputFile{Filename: "foo.csv", Version: "<some earlier commit ID of workspace dot>"}})
		testEqIFs(t, rm.Runs[1].DatasetInputFiles["b"], []InputFile{InputFile{Filename: "input.csv", Version: "<some earlier commit ID of b>"}})
		testEqIFs(t, rm.Runs[1].DatasetInputFiles["c"], []InputFile{InputFile{Filename: "cache.sqlite", Version: "<some earlier commit ID of c>"}})
		testEqStrs(t, rm.Runs[1].WorkspaceOutputFiles, []string{"log.txt"})
		testEqStrs(t, rm.Runs[1].DatasetOutputFiles["c"], []string{"cache.sqlite"})
		testEqStrs(t, rm.Runs[1].DatasetOutputFiles["d"], []string{"output.csv"})
		testEqMap(t, rm.Runs[1].Labels, map[string]string{})
		testEqMap(t, rm.Runs[1].Summary, map[string]string{"rms_error": "0.123"})
		testEqMap(t, rm.Runs[1].Parameters, map[string]string{"smoothing": "2"})
		testEqTime(t, rm.Runs[1].ExecStart, time.Date(2018, 10, 4, 13, 6, 8, 579000000, time.UTC))
		testEqTime(t, rm.Runs[1].ExecEnd, time.Date(2018, 10, 4, 13, 6, 9, 579000000, time.UTC))

		testEqStr(t, rm.Runs[2].RunID, "31df506d-c715-4159-99fd-60bb845d4dec")
		if rm.Runs[2].Authority != RunAuthority_Correction {
			t.Errorf("Expected authority %d, got %d", RunAuthority_Workload, rm.Runs[1].Authority)
		}
		testEqStr(t, rm.Runs[2].Description, "File changes were detected that the run metadata did not explain")
		testEqStrs(t, rm.Runs[2].WorkspaceOutputFiles, []string{"mylibrary.pyc"})
	}
}

func TestParseDatasetCommitMetadata(t *testing.T) {
	dcm := ParseDatasetCommitMetadata(map[string]string{
		"type":      "dotscience.run-output.v1",
		"workspace": "ID-of-dot-A",
		"run.02ecdc67-c49e-4d76-abe8-1ee13f2884b7.dataset-output-files": "[\"output.csv\"]",
		"run.cd351be8-3ba9-4c5e-ad26-429d6d6033de.dataset-output-files": "[\"output.csv\"]",
	})

	testEqStr(t, dcm.WorkspaceDotID, "ID-of-dot-A")

	r1, ok := dcm.OutputFiles["02ecdc67-c49e-4d76-abe8-1ee13f2884b7"]
	if ok {
		testEqStrs(t, r1, []string{"output.csv"})
	} else {
		t.Errorf("Did not find run 02ecdc67-c49e-4d76-abe8-1ee13f2884b7 in %#v", dcm.OutputFiles)
	}

	r2, ok := dcm.OutputFiles["cd351be8-3ba9-4c5e-ad26-429d6d6033de"]
	if ok {
		testEqStrs(t, r2, []string{"output.csv"})
	} else {
		t.Errorf("Did not find run cd351be8-3ba9-4c5e-ad26-429d6d6033de in %#v", dcm.OutputFiles)
	}
}
