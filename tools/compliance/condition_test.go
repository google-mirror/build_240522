// Copyright 2021 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package compliance

import (
	"testing"
)

func TestConditionNames(t *testing.T) {
	impliesShare := ConditionNames([]string{"restricted", "reciprocal"})

	if impliesShare.Contains("notice") {
		t.Errorf("impliesShare.Contains(\"notice\") got true, want false")
	}

	if !impliesShare.Contains("restricted") {
		t.Errorf("impliesShare.Contains(\"restricted\") got false, want true")
	}

	if !impliesShare.Contains("reciprocal") {
		t.Errorf("impliesShare.Contains(\"reciprocal\") got false, want true")
	}

	if impliesShare.Contains("") {
		t.Errorf("impliesShare.Contains(\"\") got true, want false")
	}
}
