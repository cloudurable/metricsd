package repeater

import (
	lg "github.com/cloudurable/simplelog/logging"
	c "github.com/cloudurable/metricsd/common"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
)

type AwsCloudMetricRepeater struct {
	logger lg.Logger
	conn   *cloudwatch.CloudWatch
	config *c.Config
}

const debugFormat = "{\"provider\": \"%s\", \"name\": \"%s\", \"type\": %d, \"value\": %d, \"unit\": \"%s\"}"

func (r AwsCloudMetricRepeater) RepeatForContext() bool { return true }
func (r AwsCloudMetricRepeater) RepeatForNoIdContext() bool { return true }
func (r AwsCloudMetricRepeater) Verify() bool {
    _, err := NewAWSSession(r.config)
    return err == nil
}

func (cw AwsCloudMetricRepeater) ProcessMetrics(context c.MetricContext, metrics []c.Metric) error {

	data := []*cloudwatch.MetricDatum{}
	var err error

	for index, metric := range metrics {

        dimensions := cw.createDimensions(context, metric)

        datum := &cloudwatch.MetricDatum{
            MetricName: aws.String(metric.Name),
            Timestamp:  aws.Time(metric.When),
            Dimensions: dimensions,
        }

        datum.Value = aws.Float64(metric.FloatValue)

		switch metric.Type {
		case c.MT_COUNT:	 datum.Unit = aws.String(cloudwatch.StandardUnitCount)
		case c.MT_PERCENT: 	 datum.Unit = aws.String(cloudwatch.StandardUnitPercent)
		case c.MT_MICROS: 	 datum.Unit = aws.String(cloudwatch.StandardUnitMicroseconds)
		case c.MT_MILLIS: 	 datum.Unit = aws.String(cloudwatch.StandardUnitMilliseconds)
		case c.MT_SIZE_BYTE: datum.Unit = aws.String(cloudwatch.StandardUnitBytes)
		case c.MT_SIZE_MB: 	 datum.Unit = aws.String(cloudwatch.StandardUnitMegabytes)
		case c.MT_SIZE_KB:	 datum.Unit = aws.String(cloudwatch.StandardUnitKilobytes)
		default:	         datum.Unit = aws.String(cloudwatch.StandardUnitNone)
		}

		if cw.config.Debug {
			cw.logger.Debugf(debugFormat, metric.Provider, metric.Name, metric.Type, metric.FloatValue, datum.Unit)
		}

		data = append(data, datum)

		if index%20 == 0 && index != 0 {
			data = []*cloudwatch.MetricDatum{}

			if len(data) > 0 {
				request := &cloudwatch.PutMetricDataInput{
					Namespace:  aws.String(context.GetNameSpace()),
					MetricData: data,
				}
				_, err = cw.conn.PutMetricData(request)
				if err != nil {
					cw.logger.PrintError("Error writing metrics", err)
					cw.logger.Error("Error writing metrics", err, index)
				} else {
					if cw.config.Debug {
						cw.logger.Debug("SENT..........................")
					}
				}
			}
		}
	}

	if len(data) > 0 {
		request := &cloudwatch.PutMetricDataInput{
			Namespace:  aws.String(context.GetNameSpace()),
			MetricData: data,
		}
		_, err = cw.conn.PutMetricData(request)

	}
	return err
}

func (cw AwsCloudMetricRepeater) createDimensions(context c.MetricContext, metric c.Metric) []*cloudwatch.Dimension {

    dimensions := make([]*cloudwatch.Dimension, 0, 6)

    if context.SendId() {
        instanceIdDim := &cloudwatch.Dimension{
            Name:  aws.String("InstanceId"),
            Value: aws.String(cw.config.EC2InstanceId),
        }
        dimensions = append(dimensions, instanceIdDim)

        if cw.config.IpAddress != c.EMPTY {
            ipDim := &cloudwatch.Dimension{
                Name:  aws.String("IpAddress"),
                Value: aws.String(cw.config.IpAddress),
            }
            dimensions = append(dimensions, ipDim)
        }

        if cw.config.EC2InstanceNameTag != c.EMPTY {
            dim := &cloudwatch.Dimension{
                Name:  aws.String("InstanceName"),
                Value: aws.String(cw.config.EC2InstanceNameTag),
            }
            dimensions = append(dimensions, dim)
        }
    }

    if context.GetEnv() != c.EMPTY {
        dim := &cloudwatch.Dimension{
            Name:  aws.String("Environment"),
            Value: aws.String(context.GetEnv()),
        }
        dimensions = append(dimensions, dim)
    }

    if context.GetRole() != c.EMPTY {
        dim := &cloudwatch.Dimension{
            Name:  aws.String("Role"),
            Value: aws.String(context.GetRole()),
        }
        dimensions = append(dimensions, dim)
    }

    if metric.Provider != c.EMPTY {
        dim := &cloudwatch.Dimension{
            Name:  aws.String("Provider"),
            Value: aws.String(metric.Provider),
        }
        dimensions = append(dimensions, dim)
    }

    return dimensions
}

func NewAwsCloudMetricRepeater(config *c.Config) *AwsCloudMetricRepeater {
    logger := c.EnsureLogger(nil, config.Debug, "aws-repeater")
	session, err := NewAWSSession(config) // this just verifies, no point in continuing if it won't connect
    if err != nil {
        logger.Critical(err)
        return nil
    }
	return &AwsCloudMetricRepeater{logger, cloudwatch.New(session), config}
}
