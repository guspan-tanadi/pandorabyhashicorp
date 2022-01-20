using System.Collections.Generic;
using Pandora.Definitions.Interfaces;

namespace Pandora.Definitions.ResourceManager.CostManagement.v2021_10_01.Alerts;

internal class ScopeId : ResourceID
{
    public string? CommonAlias => null;

    public string ID => "/{scope}";

    public List<ResourceIDSegment> Segments => new List<ResourceIDSegment>
    {
                new()
                {
                    Name = "scope",
                    Type = ResourceIDSegmentType.Scope
                },

    };
}
