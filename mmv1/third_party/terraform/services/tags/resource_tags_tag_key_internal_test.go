// Copyright 2026 Google Inc.
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package tags

import (
	"context"
	"io"
	"net/http"
	"strings"
	"testing"

	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

// This file lives in package tags (not tags_test) so it can call the generated, unexported
// flattenTagsTagKeyPurposeData.

type roundTripperFunc func(*http.Request) (*http.Response, error)

func (f roundTripperFunc) RoundTrip(r *http.Request) (*http.Response, error) {
	return f(r)
}

func okNetworkClient(selfLink string) *http.Client {
	return &http.Client{
		Transport: roundTripperFunc(func(r *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: http.StatusOK,
				Header:     http.Header{"Content-Type": []string{"application/json"}},
				Body:       io.NopCloser(strings.NewReader(`{"selfLink": "` + selfLink + `"}`)),
				Request:    r,
			}, nil
		}),
	}
}

// TestTagsTagKeyUpgradeAndFlattenAgree asserts the property that makes this migration safe: a state
// the upgrader produces must be identical to what a read writes for the same network. If they
// diverge, an upgraded state would still diff and force a replacement on the next plan.
func TestTagsTagKeyUpgradeAndFlattenAgree(t *testing.T) {
	const selfLink = "https://www.googleapis.com/compute/v1/projects/my-project/global/networks/123456789"

	rawState := map[string]interface{}{
		"purpose_data": map[string]interface{}{"network": "my-project/vpc-us-west1"},
	}
	upgraded, err := ResourceTagsTagKeyUpgradeV0(
		context.Background(),
		rawState,
		&transport_tpg.Config{Client: okNetworkClient(selfLink)},
	)
	if err != nil {
		t.Fatalf("error migrating state: %s", err)
	}
	upgradedNetwork := upgraded["purpose_data"].(map[string]interface{})["network"]

	// What the API returns for that same network, as the flattener receives it on read.
	flattened := flattenTagsTagKeyPurposeData(
		map[string]interface{}{"network": selfLink},
		nil,
		nil,
	)
	flattenedNetwork := flattened.(map[string]interface{})["network"]

	if upgradedNetwork != flattenedNetwork {
		t.Fatalf("upgraded state would still diff against a read:\n\nupgrader:  %#v\nflattener: %#v\n",
			upgradedNetwork, flattenedNetwork)
	}
}
