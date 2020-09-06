//    Copyright 2020 SHIFT72
//
//    Licensed under the Apache License, Version 2.0 (the "License");
//    you may not use this file except in compliance with the License.
//    You may obtain a copy of the License at
//
//        http://www.apache.org/licenses/LICENSE-2.0
//
//    Unless required by applicable law or agreed to in writing, software
//    distributed under the License is distributed on an "AS IS" BASIS,
//    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//    See the License for the specific language governing permissions and
//    limitations under the License.

package models

import "kibble/utils"

// CrewMember -
type CrewMember struct {
	Name string
	Job  string
}

// CrewMembers is an array of CrewMember
type CrewMembers []CrewMember

// GetJobNames - get a list of all the unique job titles in the list of crew members
func (crew CrewMembers) GetJobNames() StringCollection {
	var result StringCollection
	for _, c := range crew {
		result = utils.AppendUnique(c.Job, result)
	}
	return result
}

// GetMembers - get a list of all the unique people who have a particular job
func (crew CrewMembers) GetMembers(job string) StringCollection {
	var result StringCollection
	for _, c := range crew {
		if c.Job == job {
			result = utils.AppendUnique(c.Name, result)
		}

	}
	return result
}
