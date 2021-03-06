/*
Copyright (c) 2019 the Octant contributors. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package component

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"testing"

	"github.com/vmware-tanzu/octant/internal/util/json"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_PodStatus_Marshal(t *testing.T) {
	cases := []struct {
		name         string
		input        Component
		expectedFile string
		isErr        bool
	}{
		{
			name: "in general",
			input: &PodStatus{
				Config: PodStatusConfig{
					Pods: map[string]PodSummary{},
				},
			},
			expectedFile: "pod-status-marshal-general.json",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			expected, err := ioutil.ReadFile(filepath.Join("testdata", tc.expectedFile))
			require.NoError(t, err)

			got, err := json.Marshal(tc.input)
			if tc.isErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)

			assert.JSONEq(t, string(expected), string(got))
		})
	}
}

func Test_PodStatus_AddSummary(t *testing.T) {
	podStatus := NewPodStatus()

	podStatus.AddSummary("name", []Component{NewText("details")}, nil, NodeStatusOK)

	expected := NewPodStatus()
	expected.Config.Pods["name"] = PodSummary{
		Details: []Component{NewText("details")},
		Status:  NodeStatusOK,
	}
}

func Test_PodStatus_Status(t *testing.T) {
	cases := []struct {
		name     string
		statuses []NodeStatus
		expected NodeStatus
	}{
		{
			name:     "ok",
			statuses: []NodeStatus{NodeStatusOK},
			expected: NodeStatusOK,
		},
		{
			name:     "warning",
			statuses: []NodeStatus{NodeStatusOK, NodeStatusWarning},
			expected: NodeStatusWarning,
		},
		{
			name:     "error",
			statuses: []NodeStatus{NodeStatusOK, NodeStatusError},
			expected: NodeStatusError,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			podStatus := NewPodStatus()

			for i, status := range tc.statuses {
				podStatus.AddSummary(fmt.Sprintf("%d", i), []Component{NewText("details")}, nil, status)
			}

			assert.Equal(t, tc.expected, podStatus.Status())
		})
	}

}
