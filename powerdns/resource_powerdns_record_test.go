package powerdns

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

//
// Tests
//

// TODO: Move to helper block
func testPDNSRecordCommonTestCore(t *testing.T, recordConfig &PowerDNSRecordResource) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPDNSRecordDestroy,
		Steps: []resource.TestStep{
			{
				Config: recordConfig.ResourceDelcaration(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPDNSRecordExists(recordConfig.ResourceDelcaration),
				),
			},
			{
				ResourceName:      recordConfig.ResourceName,
				ImportStateId:     recordConfig.ResourceID,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// Tests
func TestAccPDNSRecord_Empty(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      testPDNSRecordConfigRecordEmpty,
				ExpectError: regexp.MustCompile("'records' must not be empty"),
			},
		},
	})
}

func TestAccPDNSRecord_A(t *testing.T) {
	testPDNSRecordCommonTestCore(t, testPDNSRecordConfigA())
}

func TestAccPDNSRecord_WithPtr(t *testing.T) {
	recordConfig = testPDNSRecordConfigAWithPtr()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPDNSRecordDestroy,
		Steps: []resource.TestStep{
			{
				Config: recordConfig.ResourceDeclaration(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPDNSRecordExists(recordConfig.ResourceName),
				),
			},
			{
				ResourceName:            recordConfig.ResourceName,
				ImportStateId:           recordConfig.ResourceID,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"set_ptr"}, // Variance from common function
			},
		},
	})
}

// TODO: Resource ID variance
/*
func TestAccPDNSRecord_WithCount(t *testing.T) {
	resourceID0 := `{"zone":"sysa.xyz.","id":"testpdnsrecordconfighyphenedwithcount-0.sysa.xyz.:::A"}`
	resourceID1 := `{"zone":"sysa.xyz.","id":"testpdnsrecordconfighyphenedwithcount-1.sysa.xyz.:::A"}`

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPDNSRecordDestroy,
		Steps: []resource.TestStep{
			{
				Config: testPDNSRecordConfigHyphenedWithCount,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPDNSRecordExists("powerdns_record.test-counted.0"),
					testAccCheckPDNSRecordExists("powerdns_record.test-counted.1"),
				),
			},
			{
				ResourceName:      "powerdns_record.test-counted[0]",
				ImportStateId:     resourceID0,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				ResourceName:      "powerdns_record.test-counted[1]",
				ImportStateId:     resourceID1,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
*/

func TestAccPDNSRecord_AAAA(t *testing.T) {
	testPDNSRecordCommonTestCore(t, testPDNSRecordConfigAAAA())
}

func TestAccPDNSRecord_CNAME(t *testing.T) {
	testPDNSRecordCommonTestCore(t, testPDNSRecordConfigCNAME())
}

func TestAccPDNSRecord_HINFO(t *testing.T) {
	testPDNSRecordCommonTestCore(t, testPDNSRecordConfigHINFO())
}

func TestAccPDNSRecord_LOC(t *testing.T) {
	testPDNSRecordCommonTestCore(t, testPDNSRecordConfigLOC())
}

func TestAccPDNSRecord_MX(t *testing.T) {
	testPDNSRecordCommonTestCore(t, testPDNSRecordConfigMX())
}

func TestAccPDNSRecord_MX(t *testing.T) {
	testPDNSRecordCommonTestCore(t, testPDNSRecordConfigMXMulti())
}

func TestAccPDNSRecord_NAPTR(t *testing.T) {
	testPDNSRecordCommonTestCore(t, testPDNSRecordConfigNAPTR())
}

func TestAccPDNSRecord_NS(t *testing.T) {
	testPDNSRecordCommonTestCore(t, testPDNSRecordConfigNS())
}

func TestAccPDNSRecord_SPF(t *testing.T) {
	testPDNSRecordCommonTestCore(t, testPDNSRecordConfigSPF())
}

func TestAccPDNSRecord_SSHFP(t *testing.T) {
	testPDNSRecordCommonTestCore(t, testPDNSRecordConfigSSHFP())
}

func TestAccPDNSRecord_SRV(t *testing.T) {
	testPDNSRecordCommonTestCore(t, testPDNSRecordConfigSRV())
}

func TestAccPDNSRecord_TXT(t *testing.T) {
	testPDNSRecordCommonTestCore(t, testPDNSRecordConfigTXT())
}

func TestAccPDNSRecord_ALIAS(t *testing.T) {
	testPDNSRecordCommonTestCore(t, testPDNSRecordConfigALIAS())
}

func TestAccPDNSRecord_SOA(t *testing.T) {
	testPDNSRecordCommonTestCore(t, testPDNSRecordConfigSOA())
}

//
// Test resource declaration functions
//
// Pattern: testPDNSRecordConfigXXX() returns a PowerDNSRecordResource struct
// The PowerDNSRecordResource struct can be used to query test config, update attributes for update tests,
// and can have ResourceDelcaration() called against it to generate the Terraform DSL resource block string.
//
type PowerDNSRecordResourceArguments struct {
	// TODO: Not following project UpCamelCase convention.
	// https://github.com/pan-net/terraform-provider-powerdns/blob/b26c0eb85f0cbe49869de72339b08ce701136a26/powerdns/client.go#L170
	count       int
	zone        string
	name        string
	record_type string
	ttl         int      // type is a Go reserved word
	records     []string
	set_ptr     bool
}

type PowerDNSRecordResource struct {
	name      string
	arguments *PowerDNSRecordResourceArguments
}

func (resourceConfig *PowerDNSRecordResource) ResourceDelcaration() (string) {
	var quotedRecords []string
	for _, record := range resourceConfig.arguments.records {
		quotedRecords = append(quotedRecords, ('"' + record + '"'))
	}

	// TODO: Can use backtics instead of double quotes to not hve to escape
	resourceDeclaration := "resource \"powerdns_record\" \"" + resourceConfig.name +"\" + {\n"
	if resourceConfig.arguments.count > 0 {
		resourceDeclaration += "  count = " + strconv.Itoa(resourceConfig.arguments.count) + "\n"
	}

	// zone, name, type, ttl, and records are mandatory
	resourceDeclaration += "  zone = \""    + resourceConfig.arguments.zone              + "\"\n"
	resourceDeclaration += "  name = \""    + resourceConfig.arguments.name              + "\"\n"
	resourceDeclaration += "  type = \""    + resourceConfig.arguments.record_type       + "\"\n"
	resourceDeclaration += "  ttl = "       + strconv.Itoa(resourceConfig.arguments.ttl) + "\n"
	resourceDeclaration += "  records = [ " + strings.Join(quotedRecords,", ")           + " ]\n"
	
	if resourceConfig.arguments.set_ptr {
		resourceDeclaration += "  set_ptr = true\n"
	}

	resourceDeclaration += "}"

	return resourceDeclaration
}

func (resourceConfig *PowerDNSRecordResource) ResourceName() (string) {
	return "powerdns_record." + resourceConfig.name
}

func (resourceConfig *PowerDNSRecordResource) ResourceID() (string) {
	return `{"zone":"` + resourceConfig.arguments.zone + `","id":"` + resourceConfig.arguments.name + ":::" + resourceConfig.arguments.record_type + `"}`
}

// Test Configs
func NewPowerDNSRecordResource() (*PowerDNSRecordResource) {
	record := &PowerDNSRecordResource{}
	record.arguments = &PowerDNSRecordResourceArguments{}

	// The zone argument is common across all the tests
	record.arguments.zone = "sysa.xyz."
	// TTL is set to 60 in the majority of the tests, default to 60 do deduplicate code.
	record.arguments.ttl = 60

	return record
}

func testPDNSRecordConfigRecordEmpty() (*PowerDNSRecordResource) {
	record := NewPowerDNSRecordResource()
	record.name                  = "test-a"
	record.arguments.name        = "testpdnsrecordconfigrecordempty.sysa.xyz."
	record.arguments.record_type = "A"
	record.arguments.records     = []
	return record
}

func testPDNSRecordConfigA() (*PowerDNSResourceArguments) {
	record := NewPowerDNSRecordResource()
	record.name                  = "test-a"
	record.arguments.name        = "testpdnsrecordconfigrecorda.sysa.xyz."
	record.arguments.record_type = "A"
	record.arguments.records = [ "1.1.1.1", "2.2.2.2" ]
	return record
}

func testPDNSRecordConfigAWithPtr() (*PowerDNSResourceArguments) {
	record := NewPowerDNSRecordResource()
	record.name                  = "test-a"
	record.arguments.name        = "testpdnsrecordconfigrecordawithptr.sysa.xyz."
	record.arguments.record_type = "A"
	record.arguments.records     = [ "1.1.1.1" ]
	record.arguments.set_ptr     = true
	return record
}

func testPDNSRecordConfigHyphenedWithCount() (*PowerDNSResourceArguments) {
	record := NewPowerDNSRecordResource()
	record.name                  = "test-counted"
	record.arguments.count       = 2
	record.arguments.name        = "testpdnsrecordconfighyphenedwithcount-${count.index}.sysa.xyz."
	record.arguments.record_type = "A"
	record.arguments.records     = [ "1.1.1.${count.index}" ]
	return record
}

func testPDNSRecordConfigAAAA() (*PowerDNSResourceArguments) {
	record := NewPowerDNSRecordResource()
	record.name                  = "test-aaaa"
	record.arguments.name        = "testpdnsrecordconfigaaaa.sysa.xyz."
	record.arguments.record_type = "AAAA"
	record.arguments.records     = [ "2001:db8:2000:bf0::1", "2001:db8:2000:bf1::1" ]
	return record
}

func testPDNSRecordConfigCNAME() (*PowerDNSResourceArguments) {
	record := NewPowerDNSRecordResource()
	record.name                  = "test-cname"
	record.arguments.name = "testpdnsrecordconfigcname.sysa.xyz."
	record.arguments.record_type = "CNAME"
	record.arguments.records = [ "redis.example.com." ]
	return record
}

func testPDNSRecordConfigCNAME() (*PowerDNSResourceArguments) {
	record := NewPowerDNSRecordResource()
	record.name                  = "test-hinfo"
	record.arguments.name = "testpdnsrecordconfighinfo.sysa.xyz."
	record.arguments.record_type = "HINFO"
	record.arguments.records = [ "\"PC-Intel-2.4ghz\" \"Linux\"" ]
	return record
}

func testPDNSRecordConfigLOC() (*PowerDNSResourceArguments) {
	record := NewPowerDNSRecordResource()
	record.name                  = "test-loc"
	record.arguments.name = "testpdnsrecordconfigloc.sysa.xyz."
	record.arguments.record_type = "LOC"
	record.arguments.records = [ "51 56 0.123 N 5 54 0.000 E 4.00m 1.00m 10000.00m 10.00m" ]
	return record
}

func testPDNSRecordConfigMX() (*PowerDNSResourceArguments) {
	record := NewPowerDNSRecordResource()
	record.name                  = "test-mx"
	record.arguments.name = "sysa.xyz."
	record.arguments.record_type = "MX"
	record.arguments.records = [ "10 mail.example.com." ]
	return record
}

func testPDNSRecordConfigMXMulti() (*PowerDNSResourceArguments) {
	record := NewPowerDNSRecordResource()
	record.name                  = "test-mx-multi"
	record.arguments.name = "multi.sysa.xyz."
	record.arguments.record_ttype = "MX"
	record.arguments.records = [ "10 mail1.example.com.", "20 mail2.example.com." ]
	return record
}

func testPDNSRecordConfigNAPTR() (*PowerDNSResourceArguments) {
	record := NewPowerDNSRecordResource()
	record.name                  = "test-naptr"
	record.arguments.name = "sysa.xyz."
	record.arguments.record_type = "NAPTR"
	record.arguments.records = [ "100 50 \"s\" \"z3950+I2L+I2C\" \"\" _z3950._tcp.gatech.edu'." ]
	return record
}

func testPDNSRecordConfigNS() (*PowerDNSResourceArguments) {
	record := NewPowerDNSRecordResource()
	record.name                  = "test-ns"
	record.arguments.name = "lab.sysa.xyz."
	record.arguments.record_type = "NS"
	record.arguments.records = [ "ns1.sysa.xyz.", "ns2.sysa.xyz." ]
	return record
}

func testPDNSRecordConfigSPF() (*PowerDNSResourceArguments) {
	record := NewPowerDNSRecordResource()
	record.name                  = "test-spf"
	record.arguments.name = "sysa.xyz."
	record.arguments.record_type = "SPF"
	record.arguments.records = [ "\"v=spf1 +all\"" ]
	return record
}

func testPDNSRecordConfigSSHFP() (*PowerDNSResourceArguments) {
	record := NewPowerDNSRecordResource()
	record.name                  = "test-sshfp"
	record.arguments.name = "ssh.sysa.xyz."
	record.arguments.record_type = "SSHFP"
	record.arguments.records = [ "1 1 123456789abcdef67890123456789abcdef67890" ]
	return record
}

func testPDNSRecordConfigSRV() (*PowerDNSResourceArguments) {
	record := NewPowerDNSRecordResource()
	record.name                  = "test-srv"
	record.arguments.name = "_redis._tcp.sysa.xyz."
	record.arguments.record_type = "SRV"
	record.arguments.records = [ "0 10 6379 redis1.sysa.xyz.", "0 10 6379 redis2.sysa.xyz.", "10 10 6379 redis-replica.sysa.xyz." ]
	return record
}

func testPDNSRecordConfigTXT() (*PowerDNSResourceArguments) {
	record := NewPowerDNSRecordResource()
	record.name                  = "test-txt"
	record.arguments.name = "text.sysa.xyz."
	record.arguments.record_type = "TXT"
	record.arguments.records = [ "\"text record payload\"" ]
	return record
}

func testPDNSRecordConfigALIAS() (*PowerDNSResourceArguments) {
	record := NewPowerDNSRecordResource()
	record.name                  = "test-alias"
	record.arguments.name = "alias.sysa.xyz."
	record.arguments.record_type = "ALIAS"
	record.arguments.ttl = 3600
	record.arguments.records = [ "www.some-alias.com." ]
	return record
}

func testPDNSRecordConfigSOAS() (*PowerDNSResourceArguments) {
	record := NewPowerDNSRecordResource()
	record.name                  = "test-soa"
	record.arguments.zone = "test-soa-sysa.xyz."
	record.arguments.name = "test-soa-sysa.xyz."
	record.arguments.record_type = "SOA"
	record.arguments.ttl = 3600
	record.arguments.records = [ "something.something. hostmaster.sysa.xyz. 2019090301 10800 3600 604800 3600" ]
	return record
}

//
// Test Helper Functions
//
func testAccCheckPDNSRecordDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "powerdns_record" {
			continue
		}

		client := testAccProvider.Meta().(*Client)
		exists, err := client.RecordExistsByID(rs.Primary.Attributes["zone"], rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("Error checking if record still exists: %#v", rs.Primary.ID)
		}
		if exists {
			return fmt.Errorf("Record still exists: %#v", rs.Primary.ID)
		}

	}
	return nil
}

func testAccCheckPDNSRecordContents(recordConfig &PowerDNSRecordResource) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[recordConfig.ResourceName()]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Record ID is set")
		}

		client := testAccProvider.Meta().(*Client)
		foundRecords, err := client.ListRecordsByID(rs.Primary.Attributes["zone"], rs.Primary.ID)
		if err != nil {
			return err
		}
		if len(foundRecords) == 0 {
			return fmt.Errorf("Record does not exist")
		}
		for _, rec := range foundRecords {
			if rec.ID() == rs.Primary.ID {
				if rec.Name != recordConfig.name {
					return fmt.Errorf("Record name field does not match: %#v : %#v", rec.Name, recordConfig.name)
				}

				if rec.Type != recordConfig.record_type {
					return fmt.Errorf("Record type field does not match: %#v : %#v", rec.Type, recordConfig.record_type)
				}

				if rec.Content != recordConfig.records {
					return fmt.Errorf("Record content field does not match: %#v : %#v", rec.Content, recordConfig.records)
				}

				if rec.TTL != recordConfig.ttl {
					return fmt.Errorf("Record TTL field does not match: %#v : %#v", rec.TTL, recordConfig.ttl)
				}

				if rec.SetPTR != recordConfig.set_ptr {
					return fmt.Errorf("Record set ptr field does not match: %#v : %#v", rec.SetPtr, recordConfig.set_ptr)
				}

				return nil
			}
		}
		return fmt.Errorf("Record does not exist: %#v", rs.Primary.ID)
	}
}
