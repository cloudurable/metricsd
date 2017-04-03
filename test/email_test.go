package test

import (
	c "github.com/cloudurable/metricsd/common"
	"testing"
)

func TestEmail(test *testing.T) {

	//config := c.Config{
	//	Debug: true,
	//}

	c.SendEmail([]string{"scottfauerbach@gmail.com"}, "This is the subject", "This is a test")
}
