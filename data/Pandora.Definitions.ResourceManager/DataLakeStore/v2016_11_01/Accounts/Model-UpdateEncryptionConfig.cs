using System;
using System.Collections.Generic;
using System.Text.Json.Serialization;
using Pandora.Definitions.Attributes;
using Pandora.Definitions.Attributes.Validation;
using Pandora.Definitions.CustomTypes;

namespace Pandora.Definitions.ResourceManager.DataLakeStore.v2016_11_01.Accounts;


internal class UpdateEncryptionConfigModel
{
    [JsonPropertyName("keyVaultMetaInfo")]
    public UpdateKeyVaultMetaInfoModel? KeyVaultMetaInfo { get; set; }
}
