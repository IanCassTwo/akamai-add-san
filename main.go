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
	"time"
	"log"
	"encoding/json"
        "github.com/akamai/AkamaiOPEN-edgegrid-golang/edgegrid"
        "github.com/akamai/AkamaiOPEN-edgegrid-golang/cps-v2"
//        "github.com/akamai/AkamaiOPEN-edgegrid-golang/client-v1"
//        "github.com/akamai/AkamaiOPEN-edgegrid-golang/configdns-v2"

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

	//
	// STEP 1
	// Get Cert Enrollment
	//

	var enrollment cps.Enrollment
	enrollment.Location = cps.GetLocation(os.Args[1])
	err = enrollment.GetEnrollment()
	if err != nil {
		log.Fatal(err)
	}

	// Step 2
	// Update enrollment with new SAN
	//

	*enrollment.CertificateSigningRequest.AlternativeNames = append(*enrollment.CertificateSigningRequest.AlternativeNames, os.Args[2])

	s,_ := json.MarshalIndent(enrollment, "", "\t")
	fmt.Println(string(s))

	updateresponse, err := enrollment.Update()
	if err != nil {
		log.Fatal(err)
	}

	s,_ = json.MarshalIndent(updateresponse, "", "\t")
	fmt.Println(string(s))


	//
	// Step 3
	// Wait for validation to happen & for DV challenges to become available
	//

        currentstatus, err := enrollment.GetStatus(*updateresponse)
        if err != nil {
		log.Fatal(err)
        }

        s,_ = json.MarshalIndent(currentstatus, "", "\t")
        fmt.Println(string(s))

        for currentstatus.StatusInfo.Status != "coodinate-domain-validation" {
                time.Sleep(10 * time.Second)

                var err error

                currentstatus, err = enrollment.GetStatus(*updateresponse)
                if err != nil {
			log.Fatal(err)
                }
                s,_ := json.MarshalIndent(currentstatus, "", "\t")
                fmt.Println(string(s))
        }

	// Note, next 2 steps are commented out because we're assuming the site is currently live (on http) and
	// redirecting path /.well-known/acme-challenge/* to dcv.akamai.com to complete the HTTP based validation

	//
	// Step 4
	// Retrieve challenges
	//

	/*
	domainvalidations, err := enrollment.GetDVChallenges(*currentstatus)
        if err != nil {
		log.Fatal(err)
        }
	*/

	//
	// Step 5
	// Update DNS with each challenge
	//

	/*
	for _, element := range domainvalidations.Dv {
		if (element.ValidationStatus != "VALIDATED") {
			for _, challenge := range element.Challenges {
				if (challenge.Type == "dns-01") {

					// Update DNS
					var recordset dnsv2.Recordset
					recordset.Name = challenge.FullPath
					recordset.Type = "CNAME"
					recordset.TTL = 30
					recordset.Rdata = append(recordset.Rdata, challenge.Token)


					req, _ := client.NewRequest(
						config, 
						"POST", 
						fmt.Sprintf("/config-dns/v2/zones/{zone}/names/{name}/types/{type}", zone, element.domain, "CNAME"
						// Incomplete
					
				}
			}
		}
	}
	*/

	//
	// STEP 6
	// Tell CPS to validate
	// Note, if we don't do this, validation will happen anyway. This just makes it a little quicker
	//
	err = enrollment.AcknowledgeDVChallenges(*currentstatus)
        if err != nil {
		log.Fatal(err)
        }

	fmt.Println("Finished")
	
}

