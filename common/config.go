package common

import (
    l "github.com/cloudurable/simplelog/logging"
    "github.com/hashicorp/hcl"
    "io/ioutil"
    "time"
    "strings"
)

type Config struct {
    Env               string        `hcl:"env"`
    Debug             bool          `hcl:"debug"`
    TimePeriodSeconds time.Duration `hcl:"interval_seconds"`
    ReadConfigSeconds time.Duration `hcl:"interval_read_config_seconds"`
    NameSpace         string        `hcl:"namespace"`

    Repeaters []string `hcl:"repeaters"`
    Gatherers []string `hcl:"gatherers"`
    Alarmers []string  `hcl:"alarmers"`

    AwsConfig       AwsConfig   `hcl:"aws"`
    KafkaConfig     KafkaConfig `hcl:"kafka"`
    SmtpConfig      SmtpConfig  `hcl:"smtp"`

    AwsAlarmerConfig AwsAlarmerConfig `hcl:"aws_alarmer"`
    EmailAlarmerConfig EmailAlarmerConfig `hcl:"email_alarmer"`

    DiskConfig      DiskGathererConfig  `hcl:"disk"`
    CpuConfig       CpuGathererConfig   `hcl:"cpu"`
    FreeConfig      FreeGathererConfig  `hcl:"free"`
    CassandraConfig CassandraGathererConfig `hcl:"cassandra"`
}

func ConfigEquals(c1 *Config, c2 *Config) bool {
    return c1.Env == c2.Env &&
        c1.Debug == c2.Debug &&
        c1.TimePeriodSeconds == c2.TimePeriodSeconds &&
        c1.ReadConfigSeconds == c2.ReadConfigSeconds &&
        c1.NameSpace == c2.NameSpace &&

        StringArraysEqual(c1.Repeaters, c2.Repeaters) &&
        StringArraysEqual(c1.Gatherers, c2.Gatherers) &&
        StringArraysEqual(c1.Alarmers, c2.Alarmers) &&

        AwsConfigEquals(&c1.AwsConfig, &c2.AwsConfig) &&
        SmtpConfigEquals(&c1.SmtpConfig, &c2.SmtpConfig) &&

        EmailAlarmerConfigEquals(&c1.EmailAlarmerConfig, &c2.EmailAlarmerConfig) &&
        DiskGathererConfigEquals(&c1.DiskConfig, &c2.DiskConfig) &&
        CpuGathererConfigEquals(&c1.CpuConfig, &c2.CpuConfig) &&
        FreeGathererConfigEquals(&c1.FreeConfig, &c2.FreeConfig) &&
        CassandraGathererConfigEquals(&c1.CassandraConfig, &c2.CassandraConfig)
}

type AwsConfig struct {
    Local              bool   `hcl:"local"`
    Region             string `hcl:"region"`
    ServerRole         string `hcl:"server_role"`
    EC2InstanceId      string `hcl:"ec2_instance_id"`
    EC2InstanceNameTag string `hcl:"ec2_instance_name"`
    IpAddress          string `hcl:"ip_address"`
}

func AwsConfigEquals(c1 *AwsConfig, c2 *AwsConfig) bool {
    return c1.Local == c2.Local &&
        c1.Region == c2.Region &&
        c1.ServerRole == c2.ServerRole &&
        c1.EC2InstanceId == c2.EC2InstanceId &&
        c1.EC2InstanceNameTag == c2.EC2InstanceNameTag &&
        c1.IpAddress == c2.IpAddress
}

type KafkaConfig struct {
    ClusterAddresses []string `hcl:"cluster_addresses"`
    Topic              string `hcl:"topic"`
    Partitions         int `hcl:"partitions"`
}

func KafkaConfigEquals(c1 *KafkaConfig, c2 *KafkaConfig) bool {
    return StringArraysEqual(c1.ClusterAddresses, c2.ClusterAddresses) &&
        c1.Topic == c2.Topic &&
        c1.Partitions == c2.Partitions
}

type SmtpConfig struct {
    Host        string `hcl:"host"`
    Port        int    `hcl:"port"`
    Username    string `hcl:"username"`
    Password    string `hcl:"password"`
    FromAddress string `hcl:"from_address"`
    IgnoreCert  bool   `hcl:"ignore_cert"`
}

func SmtpConfigEquals(c1 *SmtpConfig, c2 *SmtpConfig) bool {
    return c1.Host == c2.Host &&
        c1.Port == c2.Port &&
        c1.Username == c2.Username &&
        c1.Password == c2.Password &&
        c1.FromAddress == c2.FromAddress &&
        c1.IgnoreCert == c2.IgnoreCert
}

type AwsAlarmerConfig struct {
    DiskAlarmArns   []string `hcl:"disk_alarm_arns"`
}

func AwsAlarmerConfigEquals(c1 *AwsAlarmerConfig, c2 *AwsAlarmerConfig) bool {
    return StringArraysEqual(c1.DiskAlarmArns, c2.DiskAlarmArns)
}

type EmailAlarmerConfig struct {
    ResendIntervalSeconds int      `hcl:"resend_interval_seconds"`
    DiskAlarmTos          []string `hcl:"disk_alarm_tos"`
}

func EmailAlarmerConfigEquals(c1 *EmailAlarmerConfig, c2 *EmailAlarmerConfig) bool {
    return c1.ResendIntervalSeconds == c2.ResendIntervalSeconds &&
        StringArraysEqual(c1.DiskAlarmTos, c2.DiskAlarmTos)
}

type DiskGathererConfig struct {
    Command          string   `hcl:"command"`
    FileSystems      []string `hcl:"file_systems"`
    Fields           []string `hcl:"fields"`
    AlarmThreshold   int      `hcl:"alarm_threshold"`
}

func DiskGathererConfigEquals(c1 *DiskGathererConfig, c2 *DiskGathererConfig) bool {
    return c1.Command == c2.Command &&
        StringArraysEqual(c1.FileSystems, c2.FileSystems) &&
        StringArraysEqual(c1.Fields, c2.Fields) &&
        c1.AlarmThreshold == c2.AlarmThreshold
}

type CpuGathererConfig struct {
    ProcStat    string `hcl:"proc_stat"`
    ReportZeros bool   `hcl:"report_zeros"`
}

func CpuGathererConfigEquals(c1 *CpuGathererConfig, c2 *CpuGathererConfig) bool {
    return c1.ProcStat == c2.ProcStat && c1.ReportZeros == c2.ReportZeros
}

type FreeGathererConfig struct {
    Command string `hcl:"command"`
}

func FreeGathererConfigEquals(c1 *FreeGathererConfig, c2 *FreeGathererConfig) bool {
    return c1.Command == c2.Command
}

type CassandraGathererConfig struct {
    CqlshCommand      string   `hcl:"cqlsh_command"`
    NodetoolCommand   string   `hcl:"nodetool_command"`
    NodetoolFunctions []string `hcl:"nodetool_functions"`
}

func CassandraGathererConfigEquals(c1 *CassandraGathererConfig, c2 *CassandraGathererConfig) bool {
    return c1.CqlshCommand == c2.CqlshCommand && c1.NodetoolCommand == c2.NodetoolCommand && StringArraysEqual(c1.NodetoolFunctions, c2.NodetoolFunctions)
}

func ReadConfigString(label string, configured string, dflt string, logger l.Logger) string {

    if configured == EMPTY {
        logger.Debugf("%s initialized using default value of '%s'", label, dflt)
        return dflt
    }

    logger.Debugf("%s initialized using config value of '%s'", label, configured)
    return configured
}

func ReadConfigStringArray(label string, configured []string, dflt []string, logger l.Logger, lowerIt bool) []string {

    if configured != nil && len(configured) > 0 {
        logger.Debugf("%s initialized using config value of '%s'", label, ArrayToString(configured))
        if lowerIt {
            for x := 0; x < len(configured); x++ {
                configured[x] = strings.ToLower(configured[x])
            }
        }
        return configured
    }

    logger.Debugf("%s initialized using default value of '%s'", label, ArrayToString(dflt))
    return dflt
}

func LoadConfig(filename string, logger l.Logger) (*Config, error) {

    if logger == nil {
        logger = l.NewSimpleLogger("config")
    }

    logger.Printf("Loading config %s", filename)

    var err error

    configBytes, err := ioutil.ReadFile(filename)
    if err != nil {
        return nil, err
    }
    return LoadConfigFromString(string(configBytes), logger)
}

func LoadConfigFromString(data string, logger l.Logger) (*Config, error) {

    if logger == nil {
        logger = l.NewSimpleLogger("config")
    }

    config := &Config{}
    logger.Println("Loading log...")

    err := hcl.Decode(&config, data)
    if err != nil {
        return nil, err
    }

    if config.TimePeriodSeconds <= 0 {
        config.TimePeriodSeconds = 30
    }

    if config.ReadConfigSeconds <= 0 {
        config.ReadConfigSeconds = 60
    }

    if config.NameSpace == "" {
        config.NameSpace = "Linux System"
    }

    return config, nil

}
