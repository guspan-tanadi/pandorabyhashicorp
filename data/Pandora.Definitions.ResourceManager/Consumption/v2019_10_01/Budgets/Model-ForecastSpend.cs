using System;
using System.Collections.Generic;
using System.Text.Json.Serialization;
using Pandora.Definitions.Attributes;
using Pandora.Definitions.Attributes.Validation;
using Pandora.Definitions.CustomTypes;


// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.


namespace Pandora.Definitions.ResourceManager.Consumption.v2019_10_01.Budgets;


internal class ForecastSpendModel
{
    [JsonPropertyName("amount")]
    public float? Amount { get; set; }

    [JsonPropertyName("unit")]
    public string? Unit { get; set; }
}
