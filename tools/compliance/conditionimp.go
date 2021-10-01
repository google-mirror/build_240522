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

// licenseConditionImp implements LicenseCondition.
type licenseConditionImp struct {
	name string
	origin targetNodeImp
}

// Name returns the name of the condition. e.g. "restricted" or "notice"
func (c licenseConditionImp) Name() string {
	return c.name
}

// Origin identifies the TargetNode where the condition originates.
func (c licenseConditionImp) Origin() TargetNode {
	return c.origin
}
