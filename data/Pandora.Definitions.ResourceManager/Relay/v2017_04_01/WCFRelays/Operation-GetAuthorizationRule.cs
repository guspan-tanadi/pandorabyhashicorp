using Pandora.Definitions.Attributes;
using Pandora.Definitions.CustomTypes;
using Pandora.Definitions.Interfaces;
using Pandora.Definitions.Operations;
using System;
using System.Collections.Generic;
using System.Net;

namespace Pandora.Definitions.ResourceManager.Relay.v2017_04_01.WCFRelays;

internal class GetAuthorizationRuleOperation : Operations.GetOperation
{
    public override ResourceID? ResourceId() => new WcfRelayAuthorizationRuleId();

    public override Type? ResponseObject() => typeof(AuthorizationRuleModel);


}
