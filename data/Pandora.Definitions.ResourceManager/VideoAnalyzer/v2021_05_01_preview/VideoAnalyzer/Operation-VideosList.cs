using Pandora.Definitions.Attributes;
using Pandora.Definitions.CustomTypes;
using Pandora.Definitions.Interfaces;
using Pandora.Definitions.Operations;
using System;
using System.Collections.Generic;
using System.Net;

namespace Pandora.Definitions.ResourceManager.VideoAnalyzer.v2021_05_01_preview.VideoAnalyzer;

internal class VideosListOperation : Operations.ListOperation
{
    public override string? FieldContainingPaginationDetails() => "@nextLink";

    public override ResourceID? ResourceId() => new VideoAnalyzerId();

    public override Type NestedItemType() => typeof(VideoEntityModel);

    public override Type? OptionsObject() => typeof(VideosListOperation.VideosListOptions);

    public override string? UriSuffix() => "/videos";

    internal class VideosListOptions
    {
        [QueryStringName("$top")]
        [Optional]
        public int Top { get; set; }
    }
}
