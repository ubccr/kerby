// Copyright 2015 Andrew E. Bruno
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

package kerby

import (
	"os"
	"testing"
)

func TestGssClientAuth(t *testing.T) {
	kc := new(KerbClient)

	service := os.Getenv("KERBY_TEST_SERVICE")
	princ := os.Getenv("KERBY_TEST_PRINC")

	if len(service) == 0 {
		t.Error("Missing KERBY_TEST_SERVICE env variable. Please set before running test")
	}
	if len(princ) == 0 {
		t.Error("Missing KERBY_TEST_PRINC env variable. Please set before running test")
	}

	err := kc.Init(service, princ)
	if err != nil {
		t.Error(err)
	}

	err = kc.Step("")
	if err != nil {
		t.Error(err)
	}

	if len(kc.Response()) == 0 {
		t.Error("Invalid kerberos ticket. Empty response")
	}
}
