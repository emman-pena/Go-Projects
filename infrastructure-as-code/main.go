package main

/**
Infrastructure as Code (IaC) refers to the practice of managing and provisioning
computing infrastructure through machine-readable configuration files rather
than through manual processes.

go get -u github.com/aws/aws-sdk-go
managing cloud infrastructure (e.g., on AWS), you can use the AWS SDK for Go

Google Cloud SDK for Go: https://github.com/googleapis/google-cloud-go
Azure SDK for Go: https://github.com/Azure/azure-sdk-for-go
*/

/**
fmt: Used for formatted I/O operations (printing output).

log: Used for logging errors and other important information.

github.com/aws/aws-sdk-go/aws: This is part of the AWS SDK for Go that
provides the core AWS SDK functionality.

github.com/aws/aws-sdk-go/aws/session: This package helps in managing AWS
sessions, which are needed to make requests to AWS services.

github.com/aws/aws-sdk-go/service/ec2: This package contains the API methods
for managing EC2 instances.
*/
import (
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

func main() {
	// Create a new session in the "us-west-2" region.

	/**
	  session.NewSession: This function is used to create a new AWS session.
	  A session is needed to interact with AWS services.
	  aws.Config{Region: aws.String("us-west-2")}: This configuration specifies
	  the AWS region where the resources will be created
	  (in this case, "us-west-2").
	  If there’s an error while creating the session, the program logs the error
	  and terminates with log.Fatalf.
	*/
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-west-2")},
	)
	if err != nil {
		log.Fatalf("Unable to create session: %v", err)
	}

	// Create an EC2 service client.
	/**
	  ec2.New(sess): This creates a new EC2 service client using the session
	  (sess) we just created. The EC2 client (svc) allows us to interact
	  with EC2 resources such as instances, volumes, and security groups.
	*/
	svc := ec2.New(sess)

	// Run the EC2 instance.
	/**
	  svc.RunInstances: This is the function that actually creates a new EC2
	  instance.

	  &ec2.RunInstancesInput{...}: This is where we define the input parameters
	  for the instance we want to create:

	  ImageId: The ID of the Amazon Machine Image (AMI) that will be used for
	  the instance. In this case, it’s an AMI for Amazon Linux 2
	  (ami-0c55b159cbfafe1f0).

	  InstanceType: Defines the type of the EC2 instance (e.g., t2.micro).
	  This is a small instance, commonly used for testing or light workloads.

	  MinCount and MaxCount: These define how many instances to create.
	  Both are set to 1, meaning we’ll create exactly one EC2 instance.
	*/
	runResult, err := svc.RunInstances(&ec2.RunInstancesInput{
		ImageId:      aws.String("ami-0c55b159cbfafe1f0"), // Example AMI ID for Amazon Linux 2
		InstanceType: aws.String("t2.micro"),
		MinCount:     aws.Int64(1),
		MaxCount:     aws.Int64(1),
	})

	if err != nil {
		log.Fatalf("Could not create instance: %v", err)
	}

	fmt.Printf("Created instance %s\n", *runResult.Instances[0].InstanceId)
}

/**
Advanced Usage: Managing Infrastructure with Go
To create more advanced infrastructure, you can define and use more complex
APIs from the SDKs, such as:

Managing resources like S3 buckets, DynamoDB tables, RDS instances (AWS).
Automating Kubernetes deployments or managing Google Cloud services.
Defining networking rules, IAM roles, and permissions for secure access control.

Go for Automation in IaC
Go can also be used to create automation scripts that help with provisioning,
monitoring, and managing infrastructure resources. You might want to:

Set up regular tasks (e.g., creating backups, scaling resources).
Integrate Go with continuous integration/continuous deployment (CI/CD)
pipelines.
*/
