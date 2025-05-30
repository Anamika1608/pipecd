// Copyright 2024 The PipeCD Authors.
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

package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMakeCommentBody(t *testing.T) {
	testcases := []struct {
		name      string
		proposals []ReleaseProposal
		exists    []ReleaseProposal
		expected  string
	}{
		{
			name:     "no release",
			expected: "testdata/no-release-comment.txt",
		},
		{
			name: "one release",
			proposals: []ReleaseProposal{
				{
					ReleaseNote: "Release note for tag 1",
				},
			},
			expected: "testdata/one-release-comment.txt",
		},
		{
			name: "multiple releases",
			proposals: []ReleaseProposal{
				{
					ReleaseNote: "Release note for tag 1",
				},
				{
					ReleaseNote: "Release note for tag 2",
				},
			},
			expected: "testdata/multi-release-comment.txt",
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			got := makeCommentBody(tc.proposals, tc.exists)
			expected, err := testdata.ReadFile(tc.expected)
			require.NoError(t, err)

			assert.Equal(t, string(expected), got)
		})
	}
}
