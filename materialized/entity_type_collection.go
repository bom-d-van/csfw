// Copyright 2015 CoreStore Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package materialized

import "github.com/corestoreio/csfw/eav"

// csEntityTypeCollection contains all entity types mapped to their Go types/interfaces
var csEntityTypeCollection eav.CSEntityTypeSlice

// GetEntityTypeCollection to avoid leaking global variable
func GetEntityTypeCollection() eav.CSEntityTypeSlice {
	return csEntityTypeCollection
}