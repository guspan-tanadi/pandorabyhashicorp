using System.Collections.Generic;
using Pandora.Definitions.Interfaces;

namespace Pandora.Definitions.ResourceManager.ServiceFabricManagedCluster.v2021_05_01.ManagedClusterVersion;

internal class ManagedClusterVersionId : ResourceID
{
    public string? CommonAlias => null;

    public string ID => "/subscriptions/{subscriptionId}/providers/Microsoft.ServiceFabric/locations/{location}/managedClusterVersions/{clusterVersion}";

    public List<ResourceIDSegment> Segments => new List<ResourceIDSegment>
    {
                new()
                {
                    Name = "staticSubscriptions",
                    Type = ResourceIDSegmentType.Static,
                    FixedValue = "subscriptions"
                },

                new()
                {
                    Name = "subscriptionId",
                    Type = ResourceIDSegmentType.SubscriptionId
                },

                new()
                {
                    Name = "staticProviders",
                    Type = ResourceIDSegmentType.Static,
                    FixedValue = "providers"
                },

                new()
                {
                    Name = "staticMicrosoftServiceFabric",
                    Type = ResourceIDSegmentType.ResourceProvider,
                    FixedValue = "Microsoft.ServiceFabric"
                },

                new()
                {
                    Name = "staticLocations",
                    Type = ResourceIDSegmentType.Static,
                    FixedValue = "locations"
                },

                new()
                {
                    Name = "location",
                    Type = ResourceIDSegmentType.UserSpecified
                },

                new()
                {
                    Name = "staticManagedClusterVersions",
                    Type = ResourceIDSegmentType.Static,
                    FixedValue = "managedClusterVersions"
                },

                new()
                {
                    Name = "clusterVersion",
                    Type = ResourceIDSegmentType.UserSpecified
                },

    };
}
