package aws

import (
	"fmt"
	"testing"

	"github.com/awslabs/aws-sdk-go/aws"
	"github.com/awslabs/aws-sdk-go/service/ec2"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccAWSEBSVolume(t *testing.T) {
	var v ec2.Volume
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccAwsEbsVolumeConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVolumeExists("aws_ebs_volume.test", &v),
				),
			},
		},
	})
}

func testAccCheckVolumeExists(n string, v *ec2.Volume) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		conn := testAccProvider.Meta().(*AWSClient).ec2conn

		request := &ec2.DescribeVolumesInput{
			VolumeIDs: []*string{aws.String(rs.Primary.ID)},
		}

		response, err := conn.DescribeVolumes(request)
		if err == nil {
			if response.Volumes != nil && len(response.Volumes) > 0 {
				*v = *response.Volumes[0]
				return nil
			}
		}
		return fmt.Errorf("Error finding EC2 volume %s", rs.Primary.ID)
	}
}

const testAccAwsEbsVolumeConfig = `
resource "aws_ebs_volume" "test" {
	availability_zone = "us-west-2a"
	size = 1
}
`
