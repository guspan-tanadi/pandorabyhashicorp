using System;
using System.Collections.Generic;
using System.Text.Json.Serialization;
using Pandora.Definitions.Attributes;
using Pandora.Definitions.Attributes.Validation;
using Pandora.Definitions.CustomTypes;

namespace Pandora.Definitions.ResourceManager.StoragePool.v2021_08_01.DiskPools;


internal class DiskPoolCreateModel
{
    [JsonPropertyName("id")]
    public string? Id { get; set; }

    [JsonPropertyName("location")]
    [Required]
    public CustomTypes.Location Location { get; set; }

    [JsonPropertyName("managedBy")]
    public string? ManagedBy { get; set; }

    [JsonPropertyName("managedByExtended")]
    public List<string>? ManagedByExtended { get; set; }

    [JsonPropertyName("name")]
    public string? Name { get; set; }

    [JsonPropertyName("properties")]
    [Required]
    public DiskPoolCreatePropertiesModel Properties { get; set; }

    [JsonPropertyName("sku")]
    [Required]
    public SkuModel Sku { get; set; }

    [JsonPropertyName("tags")]
    public CustomTypes.Tags? Tags { get; set; }

    [JsonPropertyName("type")]
    public string? Type { get; set; }
}
