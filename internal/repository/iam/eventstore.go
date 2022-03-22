package iam

import (
	"github.com/caos/zitadel/internal/eventstore"
)

func RegisterEventMappers(es *eventstore.Eventstore) {
	es.RegisterFilterEventMapper(SetupStartedEventType, SetupStepMapper).
		RegisterFilterEventMapper(SetupDoneEventType, SetupStepMapper).
		RegisterFilterEventMapper(GlobalOrgSetEventType, GlobalOrgSetMapper).
		RegisterFilterEventMapper(ProjectSetEventType, ProjectSetMapper).
		RegisterFilterEventMapper(DefaultLanguageSetEventType, DefaultLanguageSetMapper).
		RegisterFilterEventMapper(SecretGeneratorAddedEventType, SecretGeneratorAddedEventMapper).
		RegisterFilterEventMapper(SecretGeneratorChangedEventType, SecretGeneratorChangedEventMapper).
		RegisterFilterEventMapper(SecretGeneratorRemovedEventType, SecretGeneratorRemovedEventMapper).
		RegisterFilterEventMapper(SMTPConfigAddedEventType, SMTPConfigAddedEventMapper).
		RegisterFilterEventMapper(SMTPConfigChangedEventType, SMTPConfigChangedEventMapper).
		RegisterFilterEventMapper(SMTPConfigPasswordChangedEventType, SMTPConfigPasswordChangedEventMapper).
		RegisterFilterEventMapper(UniqueConstraintsMigratedEventType, MigrateUniqueConstraintEventMapper).
		RegisterFilterEventMapper(SMSConfigTwilioAddedEventType, SMSConfigTwilioAddedEventMapper).
		RegisterFilterEventMapper(SMSConfigTwilioChangedEventType, SMSConfigTwilioChangedEventMapper).
		RegisterFilterEventMapper(SMSConfigTwilioTokenChangedEventType, SMSConfigTwilioTokenChangedEventMapper).
		RegisterFilterEventMapper(SMSConfigActivatedEventType, SMSConfigActivatedEventMapper).
		RegisterFilterEventMapper(SMSConfigDeactivatedEventType, SMSConfigDeactivatedEventMapper).
		RegisterFilterEventMapper(SMSConfigRemovedEventType, SMSConfigRemovedEventMapper).
		RegisterFilterEventMapper(DebugNotificationProviderFileAddedEventType, DebugNotificationProviderFileAddedEventMapper).
		RegisterFilterEventMapper(DebugNotificationProviderFileChangedEventType, DebugNotificationProviderFileChangedEventMapper).
		RegisterFilterEventMapper(DebugNotificationProviderFileRemovedEventType, DebugNotificationProviderFileRemovedEventMapper).
		RegisterFilterEventMapper(DebugNotificationProviderLogAddedEventType, DebugNotificationProviderLogAddedEventMapper).
		RegisterFilterEventMapper(DebugNotificationProviderLogChangedEventType, DebugNotificationProviderLogChangedEventMapper).
		RegisterFilterEventMapper(DebugNotificationProviderLogRemovedEventType, DebugNotificationProviderLogRemovedEventMapper).
		RegisterFilterEventMapper(OIDCSettingsAddedEventType, OIDCSettingsAddedEventMapper).
		RegisterFilterEventMapper(OIDCSettingsChangedEventType, OIDCSettingsChangedEventMapper).
		RegisterFilterEventMapper(LabelPolicyAddedEventType, LabelPolicyAddedEventMapper).
		RegisterFilterEventMapper(LabelPolicyChangedEventType, LabelPolicyChangedEventMapper).
		RegisterFilterEventMapper(LabelPolicyActivatedEventType, LabelPolicyActivatedEventMapper).
		RegisterFilterEventMapper(LabelPolicyLogoAddedEventType, LabelPolicyLogoAddedEventMapper).
		RegisterFilterEventMapper(LabelPolicyLogoRemovedEventType, LabelPolicyLogoRemovedEventMapper).
		RegisterFilterEventMapper(LabelPolicyIconAddedEventType, LabelPolicyIconAddedEventMapper).
		RegisterFilterEventMapper(LabelPolicyIconRemovedEventType, LabelPolicyIconRemovedEventMapper).
		RegisterFilterEventMapper(LabelPolicyLogoDarkAddedEventType, LabelPolicyLogoDarkAddedEventMapper).
		RegisterFilterEventMapper(LabelPolicyLogoDarkRemovedEventType, LabelPolicyLogoDarkRemovedEventMapper).
		RegisterFilterEventMapper(LabelPolicyIconDarkAddedEventType, LabelPolicyIconDarkAddedEventMapper).
		RegisterFilterEventMapper(LabelPolicyIconDarkRemovedEventType, LabelPolicyIconDarkRemovedEventMapper).
		RegisterFilterEventMapper(LabelPolicyFontAddedEventType, LabelPolicyFontAddedEventMapper).
		RegisterFilterEventMapper(LabelPolicyFontRemovedEventType, LabelPolicyFontRemovedEventMapper).
		RegisterFilterEventMapper(LabelPolicyAssetsRemovedEventType, LabelPolicyAssetsRemovedEventMapper).
		RegisterFilterEventMapper(LoginPolicyAddedEventType, LoginPolicyAddedEventMapper).
		RegisterFilterEventMapper(LoginPolicyChangedEventType, LoginPolicyChangedEventMapper).
		RegisterFilterEventMapper(OrgIAMPolicyAddedEventType, OrgIAMPolicyAddedEventMapper).
		RegisterFilterEventMapper(OrgIAMPolicyChangedEventType, OrgIAMPolicyChangedEventMapper).
		RegisterFilterEventMapper(PasswordAgePolicyAddedEventType, PasswordAgePolicyAddedEventMapper).
		RegisterFilterEventMapper(PasswordAgePolicyChangedEventType, PasswordAgePolicyChangedEventMapper).
		RegisterFilterEventMapper(PasswordComplexityPolicyAddedEventType, PasswordComplexityPolicyAddedEventMapper).
		RegisterFilterEventMapper(PasswordComplexityPolicyChangedEventType, PasswordComplexityPolicyChangedEventMapper).
		RegisterFilterEventMapper(LockoutPolicyAddedEventType, LockoutPolicyAddedEventMapper).
		RegisterFilterEventMapper(LockoutPolicyChangedEventType, LockoutPolicyChangedEventMapper).
		RegisterFilterEventMapper(PrivacyPolicyAddedEventType, PrivacyPolicyAddedEventMapper).
		RegisterFilterEventMapper(PrivacyPolicyChangedEventType, PrivacyPolicyChangedEventMapper).
		RegisterFilterEventMapper(MemberAddedEventType, MemberAddedEventMapper).
		RegisterFilterEventMapper(MemberChangedEventType, MemberChangedEventMapper).
		RegisterFilterEventMapper(MemberRemovedEventType, MemberRemovedEventMapper).
		RegisterFilterEventMapper(MemberCascadeRemovedEventType, MemberCascadeRemovedEventMapper).
		RegisterFilterEventMapper(IDPConfigAddedEventType, IDPConfigAddedEventMapper).
		RegisterFilterEventMapper(IDPConfigChangedEventType, IDPConfigChangedEventMapper).
		RegisterFilterEventMapper(IDPConfigRemovedEventType, IDPConfigRemovedEventMapper).
		RegisterFilterEventMapper(IDPConfigDeactivatedEventType, IDPConfigDeactivatedEventMapper).
		RegisterFilterEventMapper(IDPConfigReactivatedEventType, IDPConfigReactivatedEventMapper).
		RegisterFilterEventMapper(IDPOIDCConfigAddedEventType, IDPOIDCConfigAddedEventMapper).
		RegisterFilterEventMapper(IDPOIDCConfigChangedEventType, IDPOIDCConfigChangedEventMapper).
		RegisterFilterEventMapper(IDPJWTConfigAddedEventType, IDPJWTConfigAddedEventMapper).
		RegisterFilterEventMapper(IDPJWTConfigChangedEventType, IDPJWTConfigChangedEventMapper).
		RegisterFilterEventMapper(LoginPolicyIDPProviderAddedEventType, IdentityProviderAddedEventMapper).
		RegisterFilterEventMapper(LoginPolicyIDPProviderRemovedEventType, IdentityProviderRemovedEventMapper).
		RegisterFilterEventMapper(LoginPolicyIDPProviderCascadeRemovedEventType, IdentityProviderCascadeRemovedEventMapper).
		RegisterFilterEventMapper(LoginPolicySecondFactorAddedEventType, SecondFactorAddedEventMapper).
		RegisterFilterEventMapper(LoginPolicySecondFactorRemovedEventType, SecondFactorRemovedEventMapper).
		RegisterFilterEventMapper(LoginPolicyMultiFactorAddedEventType, MultiFactorAddedEventMapper).
		RegisterFilterEventMapper(LoginPolicyMultiFactorRemovedEventType, MultiFactorRemovedEventMapper).
		RegisterFilterEventMapper(MailTemplateAddedEventType, MailTemplateAddedEventMapper).
		RegisterFilterEventMapper(MailTemplateChangedEventType, MailTemplateChangedEventMapper).
		RegisterFilterEventMapper(MailTextAddedEventType, MailTextAddedEventMapper).
		RegisterFilterEventMapper(MailTextChangedEventType, MailTextChangedEventMapper).
		RegisterFilterEventMapper(CustomTextSetEventType, CustomTextSetEventMapper).
		RegisterFilterEventMapper(CustomTextRemovedEventType, CustomTextRemovedEventMapper).
		RegisterFilterEventMapper(CustomTextTemplateRemovedEventType, CustomTextTemplateRemovedEventMapper).
		RegisterFilterEventMapper(FeaturesSetEventType, FeaturesSetEventMapper).
		RegisterFilterEventMapper(InstanceDomainAddedEventType, DomainAddedEventMapper).
		RegisterFilterEventMapper(InstanceDomainRemovedEventType, DomainRemovedEventMapper)
}
