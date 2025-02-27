package resources

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/configservice"
	"github.com/aws/aws-sdk-go/service/configservice/configserviceiface"
	"github.com/gruntwork-io/cloud-nuke/config"
	"github.com/gruntwork-io/cloud-nuke/telemetry"
	"github.com/stretchr/testify/require"
	"regexp"
	"testing"
)

type mockedConfigServiceRecorders struct {
	configserviceiface.ConfigServiceAPI
	DescribeConfigurationRecordersOutput configservice.DescribeConfigurationRecordersOutput
	DeleteConfigurationRecorderOutput    configservice.DeleteConfigurationRecorderOutput
}

func (m mockedConfigServiceRecorders) DescribeConfigurationRecorders(input *configservice.DescribeConfigurationRecordersInput) (*configservice.DescribeConfigurationRecordersOutput, error) {
	return &m.DescribeConfigurationRecordersOutput, nil
}

func (m mockedConfigServiceRecorders) DeleteConfigurationRecorder(input *configservice.DeleteConfigurationRecorderInput) (*configservice.DeleteConfigurationRecorderOutput, error) {
	return &m.DeleteConfigurationRecorderOutput, nil
}

func TestConfigServiceRecorder_GetAll(t *testing.T) {
	telemetry.InitTelemetry("cloud-nuke", "")
	t.Parallel()

	testName1 := "test-recorder-1"
	testName2 := "test-recorder-2"
	csr := ConfigServiceRecorders{
		Client: mockedConfigServiceRecorders{
			DescribeConfigurationRecordersOutput: configservice.DescribeConfigurationRecordersOutput{
				ConfigurationRecorders: []*configservice.ConfigurationRecorder{
					{Name: aws.String(testName1)},
					{Name: aws.String(testName2)},
				},
			},
		},
	}

	tests := map[string]struct {
		configObj config.ResourceType
		expected  []string
	}{
		"emptyFilter": {
			configObj: config.ResourceType{},
			expected:  []string{testName1, testName2},
		},
		"nameExclusionFilter": {
			configObj: config.ResourceType{
				ExcludeRule: config.FilterRule{
					NamesRegExp: []config.Expression{{
						RE: *regexp.MustCompile(testName1),
					}}},
			},
			expected: []string{testName2},
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			names, err := csr.getAll(config.Config{
				ConfigServiceRecorder: tc.configObj,
			})

			require.NoError(t, err)
			require.Equal(t, tc.expected, aws.StringValueSlice(names))
		})
	}
}

func TestConfigServiceRecorder_NukeAll(t *testing.T) {
	telemetry.InitTelemetry("cloud-nuke", "")
	t.Parallel()

	csr := ConfigServiceRecorders{
		Client: mockedConfigServiceRecorders{
			DeleteConfigurationRecorderOutput: configservice.DeleteConfigurationRecorderOutput{},
		},
	}

	err := csr.nukeAll([]string{"test"})
	require.NoError(t, err)
}
