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
/*

Package compliance provides an approved means for reading, consuming, and
analyzing license metadata graphs.

Assuming the license metadata and dependencies are fully and accurately
recorded in the build system, any discrepancy between the official policy for
open source license compliance and this code is a bug in this code.

A few principal types to understand are LicenseGraph, LicenseCondition, and
ResolutionSet.

LicenseGraph
------------

A LicenseGraph is an immutable graph of the targets and dependencies reachable
from a specific set of root targets. In general, the root targets will be the
artifacts in a release or distribution. While conceptually immutable, parts of
the graph may be loaded or evaluated lazily.

LicenseCondition
----------------

A LicenseCondition is an immutable tuple pairing a condition name with an
originating target. e.g. Per current policy, a static library licensed under an
MIT license would pair a "notice" condition with the static library target, and
a dynamic license licensed under GPL would pair a "restricted" condition with
the dynamic library target.

ResolutionSet
-------------

A ResolutionSet is an immutable set of targets and the license conditions which
apply to each of those targets in a given context. Remember: Each license
condition pairs a condition name with an originating target so each resolution
in a ResolutionSet has a target it applies to and a target it originates from,
which may be the same target.

A ResolutionSet may at times require careful consideration. Consider a reciprocal
library. The reciprocal condition originating at the library will apply to every
target that links the library, and the ResolutionSet will attach conditions to
each of those targets. However, no matter how many targets the condition applies
to, only the source-code for the originating library needs to be shared to
satisfy the policy.
*/
package compliance
