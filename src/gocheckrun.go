package gocheck

import (
    "testing"
    "flag"
    "fmt"
)

// -----------------------------------------------------------------------
// Test suite registry.

var allSuites []interface{}

func Suite(suite interface{}) interface{} {
    return Suites(suite)[0]
}

func Suites(suites ...interface{}) []interface{} {
    lenAllSuites := len(allSuites)
    lenSuites := len(suites)
    if lenAllSuites + lenSuites > cap(allSuites) {
        newAllSuites := make([]interface{}, (lenAllSuites+lenSuites)*2)
        copy(newAllSuites, allSuites)
        allSuites = newAllSuites
    }
    allSuites = allSuites[0:lenAllSuites+lenSuites]
    for i, suite := range suites {
        allSuites[lenAllSuites+i] = suite
    }
    return suites
}


// -----------------------------------------------------------------------
// Public running interface.

var filterFlag = flag.String("f", "",
                             "Regular expression to select " +
                             "what to run (gocheck)")

func TestingT(testingT *testing.T) {
    result := RunAll(&RunConf{Filter: *filterFlag})
    println(result.String())
    if !result.Passed() {
        testingT.Fail()
    }
}

func RunAll(runConf *RunConf) *Result {
    result := Result{}
    for _, suite := range allSuites {
        result.Add(Run(suite, runConf))
    }
    return &result
}

func Run(suite interface{}, runConf *RunConf) *Result {
    runner := newSuiteRunner(suite, runConf)
    return runner.run()
}


// -----------------------------------------------------------------------
// Result methods.

func (r *Result) Add(other *Result) {
    r.Succeeded += other.Succeeded
    r.Skipped += other.Skipped
    r.Failed += other.Failed
    r.Panicked += other.Panicked
    r.FixturePanicked += other.FixturePanicked
    r.Missed += other.Missed
}

func (r *Result) Passed() bool {
    return (r.Failed == 0 && r.Panicked == 0 &&
            r.FixturePanicked == 0 && r.Missed == 0 &&
            r.RunError == nil)
}

func (r *Result) String() string {
    if r.RunError != nil {
        return "ERROR: " + r.RunError.String()
    }

    var value string
    if r.Failed == 0 && r.Panicked == 0 && r.FixturePanicked == 0 &&
       r.Missed == 0 {
        value = "OK: "
    } else {
        value = "OOPS: "
    }
    value += fmt.Sprintf("%d passed", r.Succeeded)
    if r.Skipped != 0 {
        value += fmt.Sprintf(", %d skipped", r.Skipped)
    }
    if r.Failed != 0 {
        value += fmt.Sprintf(", %d FAILED", r.Failed)
    }
    if r.Panicked != 0 {
        value += fmt.Sprintf(", %d PANICKED", r.Panicked)
    }
    if r.FixturePanicked != 0 {
        value += fmt.Sprintf(", %d FIXTURE PANICKED", r.FixturePanicked)
    }
    if r.Missed != 0 {
        value += fmt.Sprintf(", %d MISSED", r.Missed)
    }
    return value
}
