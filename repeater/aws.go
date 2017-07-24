package repeater

import (
	lg "github.com/cloudurable/simplelog/logging"
    c "github.com/cloudurable/metricsd/common"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
)

type AwsCloudMetricRepeater struct {
	logger     lg.Logger
    config     *c.Config
    awsContext *c.AwsContext
	conn       *cloudwatch.CloudWatch
}

const debugFormat = "{\"provider\": \"%s\", \"name\": \"%s\", \"type\": %d, \"value\": %d, \"unit\": \"%s\"}"

func (this AwsCloudMetricRepeater) Verify() bool {
    _, _, err := c.NewAwsSession(this.logger, this.config)
    return err == nil
}

func (this AwsCloudMetricRepeater) Repeat(metrics []c.Metric) error {

	data := []*cloudwatch.MetricDatum{}
	var err error

	for index, metric := range metrics {

        dimensions := c.CreateDimensions(this.awsContext, this.config, metric.Provider)

        datum := &cloudwatch.MetricDatum{
            MetricName: aws.String(metric.QualifiedName()),
            Timestamp:  aws.Time(metric.When),
            Dimensions: dimensions,
        }

        datum.Value = aws.Float64(metric.FloatValue)
        datum.Unit = aws.String(metric.Type.Name())

		if this.config.Debug {
			this.logger.Debugf(debugFormat, metric.Provider, metric.QualifiedName(), metric.Type, metric.FloatValue, datum.Unit)
		}

		data = append(data, datum)

		if index%20 == 0 && index != 0 {
			data = []*cloudwatch.MetricDatum{}

			if len(data) > 0 {
				request := &cloudwatch.PutMetricDataInput{
					Namespace:  aws.String(this.config.NameSpace),
					MetricData: data,
				}
				_, err = this.conn.PutMetricData(request)
				if err != nil {
					this.logger.PrintError("Error writing metrics", err)
					this.logger.Error("Error writing metrics", err, index)
				} else {
					if this.config.Debug {
						this.logger.Debug("SENT..........................")
					}
				}
			}
		}
	}

	if len(data) > 0 {
		request := &cloudwatch.PutMetricDataInput{
			Namespace:  aws.String(this.config.NameSpace),
			MetricData: data,
		}
		_, err = this.conn.PutMetricData(request)

	}
	return err
}


func (this AwsCloudMetricRepeater) Name() string { return this.logger.Name() } // already there so just use it

func NewAwsCloudMetricRepeater(config *c.Config) *AwsCloudMetricRepeater {
    logger := c.GetLogger(config.Debug, c.REPEATER + c.REPEATER_AWS)
	awsConfig, session, err := c.NewAwsSession(logger, config) // this just verifies, no point in continuing if it won't connect
    if err != nil {
        logger.Critical(err)
        return nil
    }
	return &AwsCloudMetricRepeater{logger, config, awsConfig, cloudwatch.New(session)}
}
