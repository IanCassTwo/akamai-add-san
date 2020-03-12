/*
 * Copyright 2018. Akamai Technologies, Inc
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package main

import (
	"os"
	"fmt"
	"bytes"
	"log"
	"net/http"
	"encoding/json"
        "github.com/akamai/AkamaiOPEN-edgegrid-golang/edgegrid"
        "github.com/akamai/AkamaiOPEN-edgegrid-golang/client-v1"
        "github.com/akamai/AkamaiOPEN-edgegrid-golang/cps-v2"

) 

func main() {

	if len(os.Args) != 3 {
		log.Fatal("Usage: ", os.Args[0], " <enrollment>", " <new san>")
	}

        config, err := edgegrid.Init("~/.edgerc", "default")
        if err != nil {
		log.Fatal(err)
        }

	cps.Init(config)

	enrollment, _ := cps.GetEnrollment(fmt.Sprintf("/cps/v2/enrollments/%s", os.Args[1]))

	*enrollment.CertificateSigningRequest.AlternativeNames = append(*enrollment.CertificateSigningRequest.AlternativeNames, os.Args[2])

	s,_ := json.MarshalIndent(enrollment, "", "\t")
	fmt.Println(string(s))

	enrollmentresponse, _ := Update(config, enrollment)

	s,_ = json.MarshalIndent(enrollmentresponse, "", "\t")
	fmt.Println(string(s))

	
}

func Update(config edgegrid.Config, enrollment *cps.Enrollment) (*cps.CreateEnrollmentResponse, error) {

	req, err := newRequest(
		config,
		"PUT",
		*enrollment.Location,
		enrollment,
	)

	if err != nil {
		return nil, err
	}

	res, err := client.Do(config, req)

	if err != nil {
		return nil, err
	}

	if client.IsError(res) {
		return nil, client.NewAPIError(res)
	}

	var response cps.CreateEnrollmentResponse
	if err = client.BodyJSON(res, &response); err != nil {
		return nil, err
	}

	return &response, nil
}

func newRequest(config edgegrid.Config, method, urlStr string, body interface{}) (*http.Request, error) {
	buf := new(bytes.Buffer)
	err := json.NewEncoder(buf).Encode(body)
	if err != nil {
		return nil, err
	}

	log.Printf("[DEBUG] newRequest, buf: %s", string(buf.Bytes()))

	req, err := client.NewRequest(config, method, urlStr, buf)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/vnd.akamai.cps.enrollment.v7+json")
	req.Header.Add("Accept", "application/vnd.akamai.cps.enrollment-status.v1+json")

	return req, nil
}
