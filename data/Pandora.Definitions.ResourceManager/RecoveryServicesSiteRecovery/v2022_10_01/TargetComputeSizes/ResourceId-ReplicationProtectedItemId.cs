using System.Collections.Generic;
using Pandora.Definitions.Interfaces;


// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.


namespace Pandora.Definitions.ResourceManager.RecoveryServicesSiteRecovery.v2022_10_01.TargetComputeSizes;

internal class ReplicationProtectedItemId : ResourceID
{
    public string? CommonAlias => null;

    public string ID => "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.RecoveryServices/vaults/{resourceName}/replicationFabrics/{fabricName}/replicationProtectionContainers/{protectionContainerName}/replicationProtectedItems/{replicatedProtectedItemName}";

    public List<ResourceIDSegment> Segments => new List<ResourceIDSegment>
    {
        ResourceIDSegment.Static("staticSubscriptions", "subscriptions"),
        ResourceIDSegment.SubscriptionId("subscriptionId"),
        ResourceIDSegment.Static("staticResourceGroups", "resourceGroups"),
        ResourceIDSegment.ResourceGroup("resourceGroupName"),
        ResourceIDSegment.Static("staticProviders", "providers"),
        ResourceIDSegment.ResourceProvider("staticMicrosoftRecoveryServices", "Microsoft.RecoveryServices"),
        ResourceIDSegment.Static("staticVaults", "vaults"),
        ResourceIDSegment.UserSpecified("resourceName"),
        ResourceIDSegment.Static("staticReplicationFabrics", "replicationFabrics"),
        ResourceIDSegment.UserSpecified("fabricName"),
        ResourceIDSegment.Static("staticReplicationProtectionContainers", "replicationProtectionContainers"),
        ResourceIDSegment.UserSpecified("protectionContainerName"),
        ResourceIDSegment.Static("staticReplicationProtectedItems", "replicationProtectedItems"),
        ResourceIDSegment.UserSpecified("replicatedProtectedItemName"),
    };
}
