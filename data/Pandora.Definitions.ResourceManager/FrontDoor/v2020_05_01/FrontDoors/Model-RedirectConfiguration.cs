using System;
using System.Collections.Generic;
using System.Text.Json.Serialization;
using Pandora.Definitions.Attributes;
using Pandora.Definitions.Attributes.Validation;
using Pandora.Definitions.CustomTypes;

namespace Pandora.Definitions.ResourceManager.FrontDoor.v2020_05_01.FrontDoors;

[ValueForType("#Microsoft.Azure.FrontDoor.Models.FrontdoorRedirectConfiguration")]
internal class RedirectConfigurationModel : RouteConfigurationModel
{
    [JsonPropertyName("customFragment")]
    public string? CustomFragment { get; set; }

    [JsonPropertyName("customHost")]
    public string? CustomHost { get; set; }

    [JsonPropertyName("customPath")]
    public string? CustomPath { get; set; }

    [JsonPropertyName("customQueryString")]
    public string? CustomQueryString { get; set; }

    [JsonPropertyName("redirectProtocol")]
    public FrontDoorRedirectProtocolConstant? RedirectProtocol { get; set; }

    [JsonPropertyName("redirectType")]
    public FrontDoorRedirectTypeConstant? RedirectType { get; set; }
}
