// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package docdb_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/YakDriver/regexache"
	"github.com/aws/aws-sdk-go/service/docdb"
	sdkacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-provider-aws/internal/acctest"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	tfdocdb "github.com/hashicorp/terraform-provider-aws/internal/service/docdb"
	"github.com/hashicorp/terraform-provider-aws/internal/tfresource"
)

func TestAccDocDBClusterInstance_basic(t *testing.T) {
	ctx := acctest.Context(t)
	var v docdb.DBInstance
	resourceName := "aws_docdb_cluster_instance.test"
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, docdb.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckClusterInstanceDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccClusterInstanceConfig_basic(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckClusterInstanceExists(ctx, resourceName, &v),
					resource.TestCheckNoResourceAttr(resourceName, "apply_immediately"),
					acctest.MatchResourceAttrRegionalARN(resourceName, "arn", "rds", regexache.MustCompile(fmt.Sprintf("db:%s", rName))),
					resource.TestCheckResourceAttr(resourceName, "auto_minor_version_upgrade", "true"),
					resource.TestCheckResourceAttrSet(resourceName, "availability_zone"),
					resource.TestCheckResourceAttrSet(resourceName, "ca_cert_identifier"),
					resource.TestCheckResourceAttrSet(resourceName, "cluster_identifier"),
					resource.TestCheckResourceAttr(resourceName, "copy_tags_to_snapshot", "false"),
					resource.TestCheckResourceAttrSet(resourceName, "db_subnet_group_name"),
					resource.TestCheckResourceAttrSet(resourceName, "dbi_resource_id"),
					resource.TestCheckNoResourceAttr(resourceName, "enable_performance_insights"),
					resource.TestCheckResourceAttrSet(resourceName, "endpoint"),
					resource.TestCheckResourceAttr(resourceName, "engine", "docdb"),
					resource.TestCheckResourceAttrSet(resourceName, "engine_version"),
					resource.TestCheckResourceAttr(resourceName, "identifier", rName),
					resource.TestCheckNoResourceAttr(resourceName, "identifier_prefix"),
					resource.TestCheckResourceAttrSet(resourceName, "instance_class"),
					resource.TestCheckResourceAttr(resourceName, "kms_key_id", ""),
					resource.TestCheckNoResourceAttr(resourceName, "performance_insights_kms_key_id"),
					resource.TestCheckResourceAttrSet(resourceName, "port"),
					resource.TestCheckResourceAttrSet(resourceName, "preferred_backup_window"),
					resource.TestCheckResourceAttrSet(resourceName, "preferred_maintenance_window"),
					resource.TestCheckResourceAttr(resourceName, "promotion_tier", "0"),
					resource.TestCheckResourceAttrSet(resourceName, "publicly_accessible"),
					resource.TestCheckResourceAttrSet(resourceName, "storage_encrypted"),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "0"),
					resource.TestCheckResourceAttr(resourceName, "writer", "true"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"apply_immediately",
					"identifier_prefix",
				},
			},
			{
				Config: testAccClusterInstanceConfig_modified(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckClusterInstanceExists(ctx, resourceName, &v),
					resource.TestCheckResourceAttr(resourceName, "auto_minor_version_upgrade", "false"),
					resource.TestCheckResourceAttr(resourceName, "promotion_tier", "3"),
				),
			},
		},
	})
}

func TestAccDocDBClusterInstance_performanceInsights(t *testing.T) {
	ctx := acctest.Context(t)
	var v docdb.DBInstance
	resourceName := "aws_docdb_cluster_instance.test"
	rNamePrefix := acctest.ResourcePrefix
	rName := sdkacctest.RandomWithPrefix(rNamePrefix)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, docdb.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckClusterInstanceDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccClusterInstanceConfig_performanceInsights(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckClusterInstanceExists(ctx, resourceName, &v),
					resource.TestCheckResourceAttrSet(resourceName, "enable_performance_insights"),
					resource.TestCheckResourceAttrSet(resourceName, "performance_insights_kms_key_id"),
				),
			},

			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"apply_immediately",
					"identifier_prefix",
					"enable_performance_insights",
					"performance_insights_kms_key_id",
				},
			},
		},
	})
}

func TestAccDocDBClusterInstance_az(t *testing.T) {
	ctx := acctest.Context(t)
	var v docdb.DBInstance
	resourceName := "aws_docdb_cluster_instance.cluster_instances"
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, docdb.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckClusterInstanceDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccClusterInstanceConfig_az(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckClusterInstanceExists(ctx, resourceName, &v),
					resource.TestCheckResourceAttrSet(resourceName, "availability_zone"),
				),
			},

			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"apply_immediately",
					"identifier_prefix",
				},
			},
		},
	})
}

func TestAccDocDBClusterInstance_namePrefix(t *testing.T) {
	ctx := acctest.Context(t)
	var v docdb.DBInstance
	resourceName := "aws_docdb_cluster_instance.test"
	rNamePrefix := acctest.ResourcePrefix
	rName := sdkacctest.RandomWithPrefix(rNamePrefix)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, docdb.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckClusterInstanceDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccClusterInstanceConfig_namePrefix(rName, rNamePrefix),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckClusterInstanceExists(ctx, resourceName, &v),
					resource.TestCheckResourceAttr(resourceName, "db_subnet_group_name", rName),
					resource.TestMatchResourceAttr(resourceName, "identifier", regexache.MustCompile(fmt.Sprintf("^%s", rNamePrefix))),
				),
			},

			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"apply_immediately",
					"identifier_prefix",
				},
			},
		},
	})
}

func TestAccDocDBClusterInstance_generatedName(t *testing.T) {
	ctx := acctest.Context(t)
	var v docdb.DBInstance
	resourceName := "aws_docdb_cluster_instance.test"
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, docdb.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckClusterInstanceDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccClusterInstanceConfig_generatedName(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckClusterInstanceExists(ctx, resourceName, &v),
					resource.TestMatchResourceAttr(resourceName, "identifier", regexache.MustCompile("^tf-")),
				),
			},

			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"apply_immediately",
					"identifier_prefix",
				},
			},
		},
	})
}

func TestAccDocDBClusterInstance_kmsKey(t *testing.T) {
	ctx := acctest.Context(t)
	var v docdb.DBInstance
	resourceName := "aws_docdb_cluster_instance.cluster_instances"
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, docdb.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckClusterInstanceDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccClusterInstanceConfig_kmsKey(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckClusterInstanceExists(ctx, resourceName, &v),
					resource.TestCheckResourceAttrPair(resourceName, "kms_key_id", "aws_kms_key.foo", "arn"),
				),
			},

			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"apply_immediately",
					"identifier_prefix",
				},
			},
		},
	})
}

// https://github.com/hashicorp/terraform/issues/5350
func TestAccDocDBClusterInstance_disappears(t *testing.T) {
	ctx := acctest.Context(t)
	var v docdb.DBInstance
	resourceName := "aws_docdb_cluster_instance.test"
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, docdb.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckClusterInstanceDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccClusterInstanceConfig_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckClusterInstanceExists(ctx, resourceName, &v),
					acctest.CheckResourceDisappears(ctx, acctest.Provider, tfdocdb.ResourceClusterInstance(), resourceName),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccDocDBClusterInstance_copyTagsToSnapshot(t *testing.T) {
	ctx := acctest.Context(t)
	var v docdb.DBInstance
	resourceName := "aws_docdb_cluster_instance.cluster_instances"
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, docdb.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckClusterInstanceDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccClusterInstanceConfig_copyTagsToSnapshot(rName, true),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckClusterInstanceExists(ctx, resourceName, &v),
					resource.TestCheckResourceAttr(resourceName, "copy_tags_to_snapshot", "true"),
				),
			},
			{
				Config: testAccClusterInstanceConfig_copyTagsToSnapshot(rName, false),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckClusterInstanceExists(ctx, resourceName, &v),
					resource.TestCheckResourceAttr(resourceName, "copy_tags_to_snapshot", "false"),
				),
			},

			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"apply_immediately",
					"identifier_prefix",
				},
			},
		},
	})
}

func testAccCheckClusterInstanceDestroy(ctx context.Context) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		conn := acctest.Provider.Meta().(*conns.AWSClient).DocDBConn(ctx)

		for _, rs := range s.RootModule().Resources {
			if rs.Type != "aws_docdb_cluster_instance" {
				continue
			}

			_, err := tfdocdb.FindDBInstanceByID(ctx, conn, rs.Primary.ID)

			if tfresource.NotFound(err) {
				continue
			}

			if err != nil {
				return err
			}

			return fmt.Errorf("DocumentDB Cluster Instance %s still exists", rs.Primary.ID)
		}

		return nil
	}
}

func testAccCheckClusterInstanceExists(ctx context.Context, n string, v *docdb.DBInstance) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		conn := acctest.Provider.Meta().(*conns.AWSClient).DocDBConn(ctx)

		output, err := tfdocdb.FindDBInstanceByID(ctx, conn, rs.Primary.ID)

		if err != nil {
			return err
		}

		*v = *output

		return nil
	}
}

func testAccClusterInstanceConfig_base(rName string) string {
	return fmt.Sprintf(`
resource "aws_docdb_cluster" "test" {
  cluster_identifier  = %[1]q
  master_password     = "avoid-plaintext-passwords"
  master_username     = "tfacctest"
  skip_final_snapshot = true
}

data "aws_docdb_orderable_db_instance" "test" {
  engine                     = aws_docdb_cluster.test.engine
  preferred_instance_classes = ["db.t3.medium", "db.4tg.medium", "db.r5.large", "db.r6g.large"]
}
`, rName)
}

func testAccClusterInstanceConfig_basic(rName string) string {
	return acctest.ConfigCompose(testAccClusterInstanceConfig_base(rName), fmt.Sprintf(`
resource "aws_docdb_cluster_instance" "test" {
  identifier         = %[1]q
  cluster_identifier = aws_docdb_cluster.test.id
  instance_class     = data.aws_docdb_orderable_db_instance.test.instance_class
}
`, rName))
}

func testAccClusterInstanceConfig_modified(rName string) string {
	return acctest.ConfigCompose(testAccClusterInstanceConfig_base(rName), fmt.Sprintf(`
resource "aws_docdb_cluster_instance" "test" {
  identifier                 = %[1]q
  cluster_identifier         = aws_docdb_cluster.test.id
  instance_class             = data.aws_docdb_orderable_db_instance.test.instance_class
  auto_minor_version_upgrade = false
  promotion_tier             = 3
}
`, rName))
}

func testAccClusterInstanceConfig_performanceInsights(rName string) string {
	return acctest.ConfigCompose(acctest.ConfigAvailableAZsNoOptIn(), fmt.Sprintf(`
resource "aws_kms_key" "test" {
  description = "Terraform acc test %[1]s"

  policy = <<POLICY
{
  "Version": "2012-10-17",
  "Id": "kms-tf-1",
  "Statement": [
    {
  	  "Sid": "Enable IAM User Permissions",
	  "Effect": "Allow",
	  "Principal": {
	    "AWS": "*"
	  },
	  "Action": "kms:*",
	  "Resource": "*"
	}
  ]
}
POLICY
}
resource "aws_docdb_cluster" "default" {
  cluster_identifier  = %[1]q
  availability_zones  = [data.aws_availability_zones.available.names[0], data.aws_availability_zones.available.names[1], data.aws_availability_zones.available.names[2]]
  master_username     = "foo"
  master_password     = "mustbeeightcharaters"
  skip_final_snapshot = true
}

data "aws_docdb_orderable_db_instance" "test" {
  engine                     = "docdb"
  preferred_instance_classes = ["db.t3.medium", "db.r4.large", "db.r5.large", "db.r5.xlarge"]
}

resource "aws_docdb_cluster_instance" "test" {
  identifier                      = %[1]q
  cluster_identifier              = aws_docdb_cluster.default.id
  instance_class                  = data.aws_docdb_orderable_db_instance.test.instance_class
  promotion_tier                  = "3"
  enable_performance_insights     = true
  performance_insights_kms_key_id = aws_kms_key.test.arn
}
`, rName))
}

func testAccClusterInstanceConfig_az(rName string) string {
	return acctest.ConfigCompose(acctest.ConfigAvailableAZsNoOptIn(), fmt.Sprintf(`
resource "aws_docdb_cluster" "default" {
  cluster_identifier  = %[1]q
  availability_zones  = [data.aws_availability_zones.available.names[0], data.aws_availability_zones.available.names[1], data.aws_availability_zones.available.names[2]]
  master_username     = "foo"
  master_password     = "mustbeeightcharaters"
  skip_final_snapshot = true
}

data "aws_docdb_orderable_db_instance" "test" {
  engine                     = "docdb"
  preferred_instance_classes = ["db.t3.medium", "db.r4.large", "db.r5.large", "db.r5.xlarge"]
}

resource "aws_docdb_cluster_instance" "cluster_instances" {
  identifier         = %[1]q
  cluster_identifier = aws_docdb_cluster.default.id
  instance_class     = data.aws_docdb_orderable_db_instance.test.instance_class
  promotion_tier     = "3"
  availability_zone  = data.aws_availability_zones.available.names[0]
}
`, rName))
}

func testAccClusterInstanceConfig_namePrefix(rName, rNamePrefix string) string {
	return acctest.ConfigCompose(acctest.ConfigAvailableAZsNoOptIn(), fmt.Sprintf(`
data "aws_docdb_orderable_db_instance" "test" {
  engine                     = "docdb"
  preferred_instance_classes = ["db.t3.medium", "db.r4.large", "db.r5.large", "db.r5.xlarge"]
}

resource "aws_docdb_cluster_instance" "test" {
  identifier_prefix  = %[2]q
  cluster_identifier = aws_docdb_cluster.test.id
  instance_class     = data.aws_docdb_orderable_db_instance.test.instance_class
}

resource "aws_docdb_cluster" "test" {
  cluster_identifier   = %[1]q
  master_username      = "root"
  master_password      = "password"
  db_subnet_group_name = aws_docdb_subnet_group.test.name
  skip_final_snapshot  = true
}

resource "aws_vpc" "test" {
  cidr_block = "10.0.0.0/16"

  tags = {
    Name = %[1]q
  }
}

resource "aws_subnet" "a" {
  vpc_id            = aws_vpc.test.id
  cidr_block        = "10.0.0.0/24"
  availability_zone = data.aws_availability_zones.available.names[0]

  tags = {
    Name = "%[1]s-a"
  }
}

resource "aws_subnet" "b" {
  vpc_id            = aws_vpc.test.id
  cidr_block        = "10.0.1.0/24"
  availability_zone = data.aws_availability_zones.available.names[1]

  tags = {
    Name = "%[1]s-b"
  }
}

resource "aws_docdb_subnet_group" "test" {
  name       = %[1]q
  subnet_ids = [aws_subnet.a.id, aws_subnet.b.id]
}
`, rName, rNamePrefix))
}

func testAccClusterInstanceConfig_generatedName(rName string) string {
	return acctest.ConfigCompose(acctest.ConfigAvailableAZsNoOptIn(), fmt.Sprintf(`
data "aws_docdb_orderable_db_instance" "test" {
  engine                     = "docdb"
  preferred_instance_classes = ["db.t3.medium", "db.r4.large", "db.r5.large", "db.r5.xlarge"]
}

resource "aws_docdb_cluster_instance" "test" {
  cluster_identifier = aws_docdb_cluster.test.id
  instance_class     = data.aws_docdb_orderable_db_instance.test.instance_class
}

resource "aws_docdb_cluster" "test" {
  cluster_identifier   = %[1]q
  master_username      = "root"
  master_password      = "password"
  db_subnet_group_name = aws_docdb_subnet_group.test.name
  skip_final_snapshot  = true
}

resource "aws_vpc" "test" {
  cidr_block = "10.0.0.0/16"

  tags = {
    Name = %[1]q
  }
}

resource "aws_subnet" "a" {
  vpc_id            = aws_vpc.test.id
  cidr_block        = "10.0.0.0/24"
  availability_zone = data.aws_availability_zones.available.names[0]

  tags = {
    Name = "%[1]s-a"
  }
}

resource "aws_subnet" "b" {
  vpc_id            = aws_vpc.test.id
  cidr_block        = "10.0.1.0/24"
  availability_zone = data.aws_availability_zones.available.names[1]

  tags = {
    Name = "%[1]s-b"
  }
}

resource "aws_docdb_subnet_group" "test" {
  name       = %[1]q
  subnet_ids = [aws_subnet.a.id, aws_subnet.b.id]
}
`, rName))
}

func testAccClusterInstanceConfig_kmsKey(rName string) string {
	return acctest.ConfigCompose(acctest.ConfigAvailableAZsNoOptIn(), fmt.Sprintf(`
resource "aws_kms_key" "foo" {
  description = "Terraform acc test %[1]s"

  policy = <<POLICY
{
  "Version": "2012-10-17",
  "Id": "kms-tf-1",
  "Statement": [
    {
      "Sid": "Enable IAM User Permissions",
      "Effect": "Allow",
      "Principal": {
        "AWS": "*"
      },
      "Action": "kms:*",
      "Resource": "*"
    }
  ]
}
POLICY
}

resource "aws_docdb_cluster" "default" {
  cluster_identifier  = %[1]q
  availability_zones  = [data.aws_availability_zones.available.names[0], data.aws_availability_zones.available.names[1], data.aws_availability_zones.available.names[2]]
  master_username     = "foo"
  master_password     = "mustbeeightcharaters"
  storage_encrypted   = true
  kms_key_id          = aws_kms_key.foo.arn
  skip_final_snapshot = true
}

data "aws_docdb_orderable_db_instance" "test" {
  engine                     = "docdb"
  preferred_instance_classes = ["db.t3.medium", "db.r4.large", "db.r5.large", "db.r5.xlarge"]
}

resource "aws_docdb_cluster_instance" "cluster_instances" {
  identifier         = %[1]q
  cluster_identifier = aws_docdb_cluster.default.id
  instance_class     = data.aws_docdb_orderable_db_instance.test.instance_class
}
`, rName))
}

func testAccClusterInstanceConfig_copyTagsToSnapshot(rName string, flag bool) string {
	return acctest.ConfigCompose(acctest.ConfigAvailableAZsNoOptIn(), fmt.Sprintf(`
resource "aws_docdb_cluster" "default" {
  cluster_identifier  = %[1]q
  availability_zones  = [data.aws_availability_zones.available.names[0], data.aws_availability_zones.available.names[1], data.aws_availability_zones.available.names[2]]
  master_username     = "foo"
  master_password     = "mustbeeightcharaters"
  skip_final_snapshot = true
}

data "aws_docdb_orderable_db_instance" "test" {
  engine                     = "docdb"
  preferred_instance_classes = ["db.t3.medium", "db.t4g.medium", "db.r4.large", "db.r5.large", "db.r5.xlarge"]
}

resource "aws_docdb_cluster_instance" "cluster_instances" {
  identifier            = %[1]q
  cluster_identifier    = aws_docdb_cluster.default.id
  copy_tags_to_snapshot = %[2]t
  instance_class        = data.aws_docdb_orderable_db_instance.test.instance_class
  promotion_tier        = "3"
}
`, rName, flag))
}
