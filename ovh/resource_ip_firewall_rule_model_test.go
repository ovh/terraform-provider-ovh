package ovh

import (
	"encoding/json"
	"testing"

	ovhtypes "github.com/ovh/terraform-provider-ovh/v2/ovh/types"
)

func TestIpFirewallRuleModelToCreateTcpOptionByProtocol(t *testing.T) {
	testCases := []struct {
		name          string
		protocol      string
		tcpOption     string
		wantTcpOption bool
	}{
		{
			name:          "udp omits tcpOption when computed fields are populated",
			protocol:      "udp",
			tcpOption:     "established",
			wantTcpOption: false,
		},
		{
			name:          "icmp omits tcpOption when computed fields are populated",
			protocol:      "icmp",
			tcpOption:     "established",
			wantTcpOption: false,
		},
		{
			name:          "tcp preserves established tcpOption and fragments",
			protocol:      "tcp",
			tcpOption:     "established",
			wantTcpOption: true,
		},
		{
			name:          "tcp preserves syn tcpOption and fragments",
			protocol:      "tcp",
			tcpOption:     "syn",
			wantTcpOption: true,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			model := IpFirewallRuleModel{
				Protocol:  ovhtypes.NewTfStringValue(testCase.protocol),
				TcpOption: ovhtypes.NewTfStringValue(testCase.tcpOption),
				Fragments: ovhtypes.NewTfBoolValue(false),
			}

			payload, err := json.Marshal(model.ToCreate())
			if err != nil {
				t.Fatalf("marshal create payload: %v", err)
			}

			var fields map[string]json.RawMessage
			if err := json.Unmarshal(payload, &fields); err != nil {
				t.Fatalf("unmarshal create payload: %v", err)
			}

			_, hasTcpOption := fields["tcpOption"]
			if hasTcpOption != testCase.wantTcpOption {
				t.Errorf("tcpOption present = %t, want %t; payload: %s", hasTcpOption, testCase.wantTcpOption, payload)
			}
			if testCase.wantTcpOption && string(fields["tcpOption"]) != `{"fragments":false,"option":"`+testCase.tcpOption+`"}` {
				t.Errorf("unexpected tcpOption payload: %s", fields["tcpOption"])
			}
		})
	}
}
