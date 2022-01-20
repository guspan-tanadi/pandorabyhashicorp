using System;
using System.Collections.Generic;
using System.Text.Json.Serialization;
using Pandora.Definitions.Attributes;
using Pandora.Definitions.Attributes.Validation;
using Pandora.Definitions.CustomTypes;

namespace Pandora.Definitions.ResourceManager.VMware.v2020_03_20.Locations;


internal class QuotaModel
{
    [JsonPropertyName("hostsRemaining")]
    public Dictionary<string, int>? HostsRemaining { get; set; }

    [JsonPropertyName("quotaEnabled")]
    public QuotaEnabledConstant? QuotaEnabled { get; set; }
}
