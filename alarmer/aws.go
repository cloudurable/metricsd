package alarmer

import (
    lg "github.com/cloudurable/simplelog/logging"
    c "github.com/cloudurable/metricsd/common"
    //"github.com/aws/aws-sdk-go/aws"
    "github.com/aws/aws-sdk-go/service/cloudwatch"
    "fmt"
    "github.com/aws/aws-sdk-go/aws"
    //"github.com/aws/aws-sdk-go/aws/client"
)

type AwsCloudMetricAlarmer struct {
    logger lg.Logger
    config *c.Config
    awsContext *c.AwsContext
    conn *cloudwatch.CloudWatch
    diskAlarmArns []string
}

func (this *AwsCloudMetricAlarmer) Verify() bool {
    _, _, err := c.NewAwsSession(this.logger, this.config)
    return err == nil
}

func (this *AwsCloudMetricAlarmer) Alarm(metrics []c.Metric) error {

    if len(metrics) > 0 {
        for _, m := range metrics {
            if !m.Alarm.Empty() {
                dimensions := c.CreateDimensions(this.awsContext, this.config, m.Provider)
                var alarmActions []*string
                if this.diskAlarmArns != nil && len(this.diskAlarmArns) > 0 {
                    alarmActions = aws.StringSlice(this.diskAlarmArns)
                }
                alarmName := fmt.Sprintf("%s.%s.%s", this.config.NameSpace, this.awsContext.EC2InstanceId, m.Name)

                putInput := &cloudwatch.PutMetricAlarmInput{
                    ActionsEnabled:     aws.Bool(true),
                    AlarmActions:       alarmActions,
                    AlarmDescription:   aws.String(m.Alarm.Description),
                    AlarmName:          aws.String(alarmName),
                    ComparisonOperator: aws.String(m.Alarm.Comparison.Name()),
                    Dimensions:         dimensions,
                    EvaluationPeriods:  aws.Int64(m.Alarm.EvaluationPeriods),
                    MetricName:         aws.String(m.Name),
                    Namespace:          aws.String(this.config.NameSpace),
                    Period:             aws.Int64(m.Alarm.PeriodSeconds),
                    Statistic:          aws.String(cloudwatch.StatisticAverage),
                    Threshold:          aws.Float64(m.Alarm.Threshold),
                }

                this.logger.Debugf("Setting Alarm: '%s' [%s]", alarmName, m.Alarm.Description)
                _, err := this.conn.PutMetricAlarm(putInput)
                if err != nil {
                    this.logger.WarnError("Error Sending Alarm to AWS", err)
                }

            }
        }
    }

    return nil
}

func (this *AwsCloudMetricAlarmer) Name() string { return this.logger.Name() }

func NewAwsCloudMetricAlarmer(config *c.Config) *AwsCloudMetricAlarmer {
    logger := c.EnsureLogger(nil, config.Debug, c.ALARMER + c.ALARMER_AWS)
    awsContext, session, err := c.NewAwsSession(logger, config) // this just verifies, no point in continuing if it won't connect
    if err != nil {
        logger.Critical(err)
        return nil
    }

    if config.Debug {
        logger.Debug(c.ToJsonLabeledString(awsContext))
        logger.Debug(config.AwsAlarmerConfig.DiskAlarmArns)
    }

    return &AwsCloudMetricAlarmer{logger, config, awsContext, cloudwatch.New(session), config.AwsAlarmerConfig.DiskAlarmArns}
}
