package common

import (
	l "github.com/cloudurable/simplelog/logging"
	"github.com/aws/aws-sdk-go/aws"
	awsCredentials "github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/credentials/ec2rolecreds"
	"github.com/aws/aws-sdk-go/aws/ec2metadata"
	awsSession "github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"os"
    "github.com/aws/aws-sdk-go/service/cloudwatch"
)

type AwsContext struct {
    Region             string
    ServerRole         string
    EC2InstanceId      string
    EC2InstanceNameTag string
    IpAddress          string
}

func GetAwsContext(logger l.Logger, config *Config) *AwsContext {
    metaDataClient, session := getMetadataClient(logger, config)
    return prepareAwsContext(logger, config, metaDataClient, session)
}

func NewAwsSession(logger l.Logger, config *Config) (*AwsContext, *awsSession.Session, error) {

	metaDataClient, session := getMetadataClient(logger, config)
    awsContext := prepareAwsContext(logger, config, metaDataClient, session)
    credentials := getCredentials(logger, metaDataClient)
    var awsConfig *aws.Config

	if credentials != nil {
		awsConfig = &aws.Config{
			Credentials: credentials,
			Region:      aws.String(awsContext.Region),
			MaxRetries:  aws.Int(3),
		}
	} else {
        awsConfig = &aws.Config{
			Region:     aws.String(awsContext.Region),
			MaxRetries: aws.Int(3),
		}
	}

    metricsAwsSession, err := awsSession.NewSession(awsConfig)
    return awsContext, metricsAwsSession, err
}

func CreateDimensions(awsContext *AwsContext, config *Config, provider string) []*cloudwatch.Dimension {

    dimensions := []*cloudwatch.Dimension{}

    dimensions = appendDimension(dimensions, "InstanceId", awsContext.EC2InstanceId)
    dimensions = appendDimension(dimensions, "InstanceName", awsContext.EC2InstanceNameTag)
    dimensions = appendDimension(dimensions, "IpAddress", awsContext.IpAddress)
    dimensions = appendDimension(dimensions, "Role", awsContext.ServerRole)
    dimensions = appendDimension(dimensions, "Environment", config.Env)
    dimensions = appendDimension(dimensions, "Provider", provider)

    return dimensions
}

func appendDimension(dimensions []*cloudwatch.Dimension, name string, value string) []*cloudwatch.Dimension {
    if value != EMPTY {
        dim := &cloudwatch.Dimension{
            Name:  aws.String(name),
            Value: aws.String(value),
        }
        dimensions = append(dimensions, dim)
    }

    return dimensions
}

func getMetadataClient(logger l.Logger, config *Config) (*ec2metadata.EC2Metadata, *awsSession.Session) {
    if config.AwsConfig.Local {
        logger.Println("Config set to local")
        return nil, nil
    }

    logger.Debug("Config NOT set to local. Using meta-data client to find local")
    session, err := awsSession.NewSession(&aws.Config{})
    if err != nil {
        logger.Error("Unable to contruct meta client : %s %v", err.Error(), err)
        os.Exit(102)
    }
    return ec2metadata.New(session), session
}

func prepareAwsContext(logger l.Logger, config *Config, client *ec2metadata.EC2Metadata, session *awsSession.Session) *AwsContext {

    awsContext := AwsContext{
        config.AwsConfig.Region,
        config.AwsConfig.ServerRole,
        config.AwsConfig.EC2InstanceId,
        config.AwsConfig.EC2InstanceNameTag,
        config.AwsConfig.IpAddress,
    }

    if client == nil {
        logger.Info("Client missing using config to set region")
        if awsContext.Region == EMPTY {
            logger.Info("AwsRegion missing using default region us-west-2")
            awsContext.Region = "us-west-2"
        }
    } else {
        var err error
        awsContext.Region, err = client.Region()
        if err != nil {
            logger.Error("Unable to get region from aws meta client : %s %v", err.Error(), err)
            os.Exit(103)
        }

        awsContext.EC2InstanceId, err = client.GetMetadata("instance-id")
        if err != nil {
            logger.Error("Unable to get instance id from aws meta client : %s %v", err.Error(), err)
            os.Exit(104)
        }

        awsContext.IpAddress = findLocalIp(logger, client)
        awsContext.EC2InstanceNameTag = findInstanceName(logger, awsContext.EC2InstanceId, awsContext.Region, session)
    }

    return &awsContext
}

func findLocalIp(logger l.Logger, metaClient *ec2metadata.EC2Metadata) string {
    ip, err := metaClient.GetMetadata("local-ipv4")

    if err != nil {
        logger.Error("Unable to get private ip address from aws meta client : %s %v", err.Error(), err)
        os.Exit(106)
    }

    return ip
}

func getCredentials(logger l.Logger, client *ec2metadata.EC2Metadata) *awsCredentials.Credentials {

	if client == nil {
        logger.Info("Client missing credentials not looked up")
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

func findAZ(logger l.Logger, metaClient *ec2metadata.EC2Metadata) string {

	az, err := metaClient.GetMetadata("placement/availability-zone")

	if err != nil {
        logger.Errorf("Unable to get az from aws meta client : %s %v", err.Error(), err)
		os.Exit(105)
	}

	return az
}

func findInstanceName(logger l.Logger, instanceId string, region string, session *awsSession.Session) string {

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
        logger.Errorf("Unable to get instance name tag DescribeInstances failed : %s %v", err.Error(), err)
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
        logger.Errorf("Unable to get find name tag ")
		return name

	} else {
        logger.Errorf("Unable to get find name tag ")
		return name
	}
}
