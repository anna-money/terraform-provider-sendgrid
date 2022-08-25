package sendgrid

import (
	"errors"
	"fmt"
)

var (
	// ErrCreateRateLimit error displayed when we reach the API calls rate limit.
	ErrCreateRateLimit = errors.New("expected instance to be created but we were rate limited")

	// ErrInvalidImportFormat error displayed when the string passed to import a template version
	// doesn't have the good format.
	ErrInvalidImportFormat = errors.New("invalid import. Supported import format: {{templateID}}/{{templateVersionID}}")

	// ErrSubUserNotFound error displayed when the subUser can not be found.
	ErrSubUserNotFound = errors.New("subUser wasn't found")

	// ErrNoNewVersionFoundForTemplate error displayed when no recent version can be found for a given template.
	ErrNoNewVersionFoundForTemplate = errors.New("no recent version found for template_id")

	// ErrSetTemplateName error displayed when the provider can't set the template name.
	ErrSetTemplateName = errors.New("could not set template name")

	// ErrSetTemplateGeneration error displayed when the provider can't set the template generation.
	ErrSetTemplateGeneration = errors.New("could not set template generation")

	// ErrSetTemplateUpdatedAt error displayed when the provider can't set the template updated_at attribute.
	ErrSetTemplateUpdatedAt = errors.New("could not set template version updated_at attribute")

	// ErrSetTemplateVersionName error displayed when the provider can't set the template version name.
	ErrSetTemplateVersionName = errors.New("could not set template version name")

	// ErrSetTemplateVersionUpdatedAt error displayed when the provider
	// can't set the template version updated_at attribute.
	ErrSetTemplateVersionUpdatedAt = errors.New("could not set template version updated_at attribute")

	// ErrSetTemplateVersionActive error displayed when the provider
	// can't set the template version active attribute.
	ErrSetTemplateVersionActive = errors.New("could not set template version active attribute")

	// ErrSetTemplateVersionHTMLContent error displayed when the provider
	// can't set the template version hmtl_content attribute.
	ErrSetTemplateVersionHTMLContent = errors.New("could not set template version html_content attribute")

	// ErrSetTemplateVersionPlainContent error displayed when the provider
	// can't set the template version plain_content attribute.
	ErrSetTemplateVersionPlainContent = errors.New("could not set template version plain_content attribute")

	// ErrSetTemplateVersionGenPlainContent error displayed when the provider
	// can't set the template version generate_plain_content attribute.
	ErrSetTemplateVersionGenPlainContent = errors.New("could not set template version generate_plain_content attribute")

	// ErrSetTemplateVersionSubject error displayed when the provider
	// can't set the template version subject attribute.
	ErrSetTemplateVersionSubject = errors.New("could not set template version subject attribute")

	// ErrSetTemplateVersionEditor error displayed when the provider
	// can't set the template version editor attribute.
	ErrSetTemplateVersionEditor = errors.New("could not set template version editor attribute")

	// ErrSetTemplateVersionThumbnailURL error displayed when the provider can't set the template version thumbnail URL.
	ErrSetTemplateVersionThumbnailURL = errors.New("could not set template version thumbnail URL")

	// ErrSetUnsubscribeGroupName error displayed when the provider can't set the unsubscribe group name.
	ErrSetUnsubscribeGroupName = errors.New("could not set unsubscribe group name")

	// ErrSetUnsubscribeGroupDesc error displayed when the provider can't set the unsubscribe group description.
	ErrSetUnsubscribeGroupDesc = errors.New("could not set unsubscribe group description")

	// ErrSetUnsubscribeGroupIsDefault error displayed when the provider can't set the unsubscribe group is_default flag.
	ErrSetUnsubscribeGroupIsDefault = errors.New("could not set unsubscribe group is_default flag")

	// ErrSetUnsubscribeGroupUnsuscribes error displayed when the provider
	// can't set the unsubscribe group unsubscribes attribute.
	ErrSetUnsubscribeGroupUnsuscribes = errors.New("could not set unsubscribe group unsubscribes attribute")
)

func subUserNotFound(name string) error {
	return fmt.Errorf("%w: %s", ErrSubUserNotFound, name)
}
