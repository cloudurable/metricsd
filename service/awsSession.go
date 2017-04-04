package service

import (
	l "github.com/cloudurable/simplelog/logging"
	c "github.com/cloudurable/metricsd/common"
	"github.com/aws/aws-sdk-go/aws"
	awsCredentials "github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/credentials/ec2rolecreds"
	"github.com/aws/aws-sdk-go/aws/ec2metadata"
	awsSession "github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"os"
)

var awsSessionLogger = l.NewSimpleLogger("aws-session")

func GetAwsConfig(config *c.Config) *AwsConfig {
    metaDataClient, session := getClient(config)
    return prepareAwsConfig(config, metaDataClient, session)
}

func NewAwsSession(config *c.Config) (*awsSession.Session, error) {

	metaDataClient, session := getClient(config)
    configAwsPart := prepareAwsConfig(config, metaDataClient, session)
    credentials := getCredentials(metaDataClient)

	if credentials != nil {
		awsConfig := &aws.Config{
			Credentials: credentials,
			Region:      aws.String(configAwsPart.AWSRegion),
			MaxRetries:  aws.Int(3),
		}
		return awsSession.NewSession(awsConfig)
	} else {

		return awsSession.NewSession(&aws.Config{
			Region:     aws.String(configAwsPart.AWSRegion),
			MaxRetries: aws.Int(3),
		})
	}
}

func getClient(config *c.Config) (*ec2metadata.EC2Metadata, *awsSession.Session) {
	if !config.Local {
		awsSessionLogger.Debug("Config NOT set to local using meta-data client to find local")
		var session = awsSession.New(&aws.Config{})
		return ec2metadata.New(session), session
	} else {
		awsSessionLogger.Println("Config set to local")
		return nil, nil
	}
}

type AwsConfig struct {
	AWSRegion          string
	ServerRole         string
	IpAddress          string
	EC2InstanceId      string
	EC2InstanceNameTag string
}

func prepareAwsConfig(config *c.Config, client *ec2metadata.EC2Metadata, session *awsSession.Session) *AwsConfig {

	awsConfig := AwsConfig{
		config.AWSRegion,
		config.ServerRole,
		config.IpAddress,
		config.EC2InstanceId,
		config.EC2InstanceNameTag,
	}

	if client == nil {
		awsSessionLogger.Info("Client missing using config to set region")
		if awsConfig.AWSRegion == "" {
			awsSessionLogger.Info("AWSRegion missing using default region us-west-2")
			awsConfig.AWSRegion = "us-west-2"
		}
	} else {
		region, err := client.Region()
		if err != nil {
			awsSessionLogger.Error("Unable to get region from aws meta client : %s %v", err.Error(), err)
			os.Exit(3)
		}

		awsConfig.AWSRegion = region
		awsConfig.IpAddress = findLocalIp(client)
		awsConfig.EC2InstanceId, err = client.GetMetadata("instance-id")
		awsConfig.EC2InstanceNameTag = findInstanceName(config.EC2InstanceId, config.AWSRegion, session)
		if err != nil {
			awsSessionLogger.Error("Unable to get instance id from aws meta client : %s %v", err.Error(), err)
			os.Exit(4)
		}
	}

	return &awsConfig
}

func findLocalIp(metaClient *ec2metadata.EC2Metadata) string {
	ip, err := metaClient.GetMetadata("local-ipv4")

	if err != nil {
		awsSessionLogger.Error("Unable to get private ip address from aws meta client : %s %v", err.Error(), err)
		os.Exit(6)
	}

	return ip

}

func getCredentials(client *ec2metadata.EC2Metadata) *awsCredentials.Credentials {

	if client == nil {
		awsSessionLogger.Info("Client missing credentials not looked up")
		return nil
	} else {
		return awsCredentials.NewChainCredentials([]awsCredentials.Provider{
			&awsCredentials.EnvProvider{},
			&ec2rolecreds.EC2RoleProvider{
				Client: client,
			},
		})
	}

}

func findAZ(metaClient *ec2metadata.EC2Metadata) string {

	az, err := metaClient.GetMetadata("placement/availability-zone")

	if err != nil {
		awsSessionLogger.Errorf("Unable to get az from aws meta client : %s %v", err.Error(), err)
		os.Exit(5)
	}

	return az
}

func findInstanceName(instanceId string, region string, session *awsSession.Session) string {

	var name = "NO_NAME"
	var err error

	ec2Service := ec2.New(session, aws.NewConfig().WithRegion(region))

	params := &ec2.DescribeInstancesInput{
		InstanceIds: []*string{
			aws.String(instanceId), // Required
			// More values...
		},
	}

	resp, err := ec2Service.DescribeInstances(params)

	if err != nil {
		awsSessionLogger.Errorf("Unable to get instance name tag DescribeInstances failed : %s %v", err.Error(), err)
		return name
	}

	if len(resp.Reservations) > 0 && len(resp.Reservations[0].Instances) > 0 {
		var instance = resp.Reservations[0].Instances[0]
		if len(instance.Tags) > 0 {

			for _, tag := range instance.Tags {
				if *tag.Key == "Name" {
					return *tag.Value
				}
			}
		}
		awsSessionLogger.Errorf("Unable to get find name tag ")
		return name

	} else {
		awsSessionLogger.Errorf("Unable to get find name tag ")
		return name
	}
}
