package models

import "encoding/xml"

// XBResourceList represents the DownloadXML data for a resource
type XBResourceList struct {
	XMLName    xml.Name `xml:"XBResourceList"`
	Text       string   `xml:",chardata"`
	XBResource []struct {
		Text      string `xml:",chardata"`
		Protocol  string `xml:"protocol"`
		TypeSIPgw struct {
			Text             string `xml:",chardata"`
			PortAddress      string `xml:"portAddress"`
			ServiceState     string `xml:"serviceState"`
			Direction        string `xml:"direction"`
			NAT              string `xml:"NAT"`
			AllowDirectMedia string `xml:"allowDirectMedia"`
			SipProfileIndex  string `xml:"sipProfileIndex"`
			OptionPoll       string `xml:"optionPoll"`
			AuthorizedRPS    string `xml:"authorizedRPS"`
			UnauthorizedRPS  string `xml:"unauthorizedRPS"`
		} `xml:"typeSIPgw"`
		Name        string `xml:"name"`
		CompanyName string `xml:"companyName"`
		TrunkId     string `xml:"trunkId"`
		SgId        string `xml:"sgId"`
		Capacity    string `xml:"capacity"`
		CpsLimit    string `xml:"cpsLimit"`
		Node        struct {
			Text         string `xml:",chardata"`
			Fqdn         string `xml:"fqdn"`
			Netmask      string `xml:"netmask"`
			Capacity     string `xml:"capacity"`
			CpsLimit     string `xml:"cpsLimit"`
			CacProfileId string `xml:"cacProfileId"`
		} `xml:"node"`
		Rtid     string `xml:"rtid"`
		Ingress1 struct {
			Text    string `xml:",chardata"`
			Match   string `xml:"match"`
			Action1 string `xml:"action1"`
			Digits1 string `xml:"digits1"`
			Action2 string `xml:"action2"`
			Digits2 string `xml:"digits2"`
		} `xml:"ingress1"`
		Ingress2 struct {
			Text    string `xml:",chardata"`
			Match   string `xml:"match"`
			Action1 string `xml:"action1"`
			Digits1 string `xml:"digits1"`
			Action2 string `xml:"action2"`
			Digits2 string `xml:"digits2"`
		} `xml:"ingress2"`
		Egress1 struct {
			Text    string `xml:",chardata"`
			Match   string `xml:"match"`
			Action1 string `xml:"action1"`
			Digits1 string `xml:"digits1"`
			Action2 string `xml:"action2"`
			Digits2 string `xml:"digits2"`
		} `xml:"egress1"`
		Egress2 struct {
			Text    string `xml:",chardata"`
			Match   string `xml:"match"`
			Action1 string `xml:"action1"`
			Digits1 string `xml:"digits1"`
			Action2 string `xml:"action2"`
			Digits2 string `xml:"digits2"`
		} `xml:"egress2"`
		OutboundANI string `xml:"outboundANI"`
		TechPrefix  string `xml:"techPrefix"`
		RnIngress1  struct {
			Text    string `xml:",chardata"`
			Match   string `xml:"match"`
			Action1 string `xml:"action1"`
			Digits1 string `xml:"digits1"`
			Action2 string `xml:"action2"`
			Digits2 string `xml:"digits2"`
		} `xml:"rnIngress1"`
		RnIngress2 struct {
			Text    string `xml:",chardata"`
			Match   string `xml:"match"`
			Action1 string `xml:"action1"`
			Digits1 string `xml:"digits1"`
			Action2 string `xml:"action2"`
			Digits2 string `xml:"digits2"`
		} `xml:"rnIngress2"`
		RnEgress1 struct {
			Text    string `xml:",chardata"`
			Match   string `xml:"match"`
			Action1 string `xml:"action1"`
			Digits1 string `xml:"digits1"`
			Action2 string `xml:"action2"`
			Digits2 string `xml:"digits2"`
		} `xml:"rnEgress1"`
		RnEgress2 struct {
			Text    string `xml:",chardata"`
			Match   string `xml:"match"`
			Action1 string `xml:"action1"`
			Digits1 string `xml:"digits1"`
			Action2 string `xml:"action2"`
			Digits2 string `xml:"digits2"`
		} `xml:"rnEgress2"`
		CodecPolicy            string `xml:"codecPolicy"`
		GroupPolicy            string `xml:"groupPolicy"`
		Dtid                   string `xml:"dtid"`
		T38                    string `xml:"t38"`
		Rfc2833                string `xml:"rfc2833"`
		PayloadType            string `xml:"payloadType"`
		Tos                    string `xml:"tos"`
		SvcPortIndex           string `xml:"svcPortIndex"`
		RadiusAuthGrpIndex     string `xml:"radiusAuthGrpIndex"`
		RadiusAcctGrpIndex     string `xml:"radiusAcctGrpIndex"`
		LnpGrpIndex            string `xml:"lnpGrpIndex"`
		TeleblockGrpIndex      string `xml:"teleblockGrpIndex"`
		CnamGrpIndex           string `xml:"cnamGrpIndex"`
		ErsGrpIndex            string `xml:"ersGrpIndex"`
		MaxCallDuration        string `xml:"maxCallDuration"`
		MinCallDuration        string `xml:"minCallDuration"`
		NoAnswerTimeout        string `xml:"noAnswerTimeout"`
		NoRingTimeout          string `xml:"noRingTimeout"`
		CauseCodeProfile       string `xml:"causeCodeProfile"`
		StopRouteProfile       string `xml:"stopRouteProfile"`
		PaiAction              string `xml:"paiAction"`
		PaiString              string `xml:"paiString"`
		InheritedGenericHeader string `xml:"inheritedGenericHeader"`
		OutSMCProfileId        string `xml:"outSMCProfileId"`
	} `xml:"XBResource"`
}
