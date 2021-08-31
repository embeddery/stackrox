import React, { ReactElement } from 'react';
import { Checkbox, Form, PageSection, TextInput } from '@patternfly/react-core';
import * as yup from 'yup';

import usePageState from 'Containers/Integrations/hooks/usePageState';
import useIntegrationForm from '../useIntegrationForm';
import { IntegrationFormProps } from '../integrationFormTypes';

import IntegrationFormActions from '../IntegrationFormActions';
import FormCancelButton from '../FormCancelButton';
import FormTestButton from '../FormTestButton';
import FormSaveButton from '../FormSaveButton';
import FormMessage from '../FormMessage';
import FormLabelGroup from '../FormLabelGroup';

export type SplunkIntegration = {
    id?: string;
    name: string;
    splunk: {
        httpEndpoint: string;
        httpToken: string;
        truncate: number;
        insecure: boolean;
        auditLoggingEnabled: boolean;
        sourceTypes: {
            alert: string;
            audit: string;
        };
    };
    uiEndpoint: string;
    type: 'splunk';
    enabled: boolean;
};

export type SplunkIntegrationFormValues = {
    notifier: SplunkIntegration;
    updatePassword: boolean;
};

const validHttpEndpointRegex = /^((https?):\/\/)?([a-zA-Z0-9\-.]\.)?[a-zA-Z0-9\-.]{1,}\.[a-zA-Z]{2,}(\.[a-zA-Z]{2,})?(:[0-9]+)?(\/[a-zA-Z0-9-]*\/?)*$/;

export const validationSchema = yup.object().shape({
    notifier: yup.object().shape({
        name: yup.string().trim().required('Name is required'),
        splunk: yup.object().shape({
            httpEndpoint: yup
                .string()
                .trim()
                .required('HTTP event collector URL is required')
                .matches(validHttpEndpointRegex, 'Must be a valid server address'),
            // httpToken: yup.string().trim().required('Required'),
            httpToken: yup
                .string()
                .test(
                    'httpToken-test',
                    'HTTP token is required',
                    (value, context: yup.TestContext) => {
                        const requireHttpTokenField =
                            // eslint-disable-next-line @typescript-eslint/ban-ts-comment
                            // @ts-ignore
                            context?.from[2]?.value?.updatePassword || false;

                        if (!requireHttpTokenField) {
                            return true;
                        }

                        const trimmedValue = value?.trim();
                        return !!trimmedValue;
                    }
                ),
            truncate: yup.number().required('HEC truncate limit is required'),
            sourceTypes: yup.object().shape({
                alert: yup.string().trim().required('Source type for alert is required'),
                audit: yup.string().trim().required('Source type for audit is required'),
            }),
        }),
    }),
    updatePassword: yup.bool(),
});

export const defaultValues: SplunkIntegrationFormValues = {
    notifier: {
        name: '',
        splunk: {
            httpEndpoint: '',
            httpToken: '',
            truncate: 10000,
            insecure: false,
            auditLoggingEnabled: false,
            sourceTypes: {
                alert: 'stackrox-alert',
                audit: 'stackrox-audit-message',
            },
        },
        uiEndpoint: window.location.origin,
        type: 'splunk',
        enabled: true,
    },
    updatePassword: true,
};

function SplunkIntegrationForm({
    initialValues = null,
    isEditable = false,
}: IntegrationFormProps<SplunkIntegration>): ReactElement {
    const formInitialValues = defaultValues;
    if (initialValues) {
        formInitialValues.notifier = {
            ...formInitialValues.notifier,
            ...initialValues,
        };
        // We want to clear the password because backend returns '******' to represent that there
        // are currently stored credentials
        formInitialValues.notifier.splunk.httpToken = '';
    }
    const {
        values,
        touched,
        errors,
        dirty,
        isValid,
        setFieldValue,
        handleBlur,
        isSubmitting,
        isTesting,
        onSave,
        onTest,
        onCancel,
        message,
    } = useIntegrationForm<SplunkIntegrationFormValues, typeof validationSchema>({
        initialValues: formInitialValues,
        validationSchema,
    });
    const { isCreating } = usePageState();

    function onChange(value, event) {
        return setFieldValue(event.target.id, value);
    }

    return (
        <>
            <PageSection variant="light" isFilled hasOverflowScroll>
                {message && <FormMessage message={message} />}
                <Form isWidthLimited>
                    <FormLabelGroup
                        isRequired
                        label="Integration name"
                        fieldId="notifier.name"
                        touched={touched}
                        errors={errors}
                    >
                        <TextInput
                            isRequired
                            type="text"
                            id="notifier.name"
                            value={values.notifier.name}
                            onChange={onChange}
                            onBlur={handleBlur}
                            isDisabled={!isEditable}
                        />
                    </FormLabelGroup>
                    <FormLabelGroup
                        isRequired
                        label="HTTP event collector URL"
                        fieldId="notifier.splunk.httpEndpoint"
                        touched={touched}
                        errors={errors}
                    >
                        <TextInput
                            isRequired
                            type="text"
                            id="notifier.splunk.httpEndpoint"
                            value={values.notifier.splunk.httpEndpoint}
                            onChange={onChange}
                            onBlur={handleBlur}
                            isDisabled={!isEditable}
                        />
                    </FormLabelGroup>
                    {!isCreating && (
                        <FormLabelGroup
                            label=""
                            fieldId="updatePassword"
                            helperText="Leave this off to use the currently stored credentials."
                            errors={errors}
                        >
                            <Checkbox
                                label="Update token"
                                id="updatePassword"
                                isChecked={values.updatePassword}
                                onChange={onChange}
                                onBlur={handleBlur}
                                isDisabled={!isEditable}
                            />
                        </FormLabelGroup>
                    )}
                    <FormLabelGroup
                        isRequired={values.updatePassword}
                        label="HTTP event collector token"
                        fieldId="notifier.splunk.httpToken"
                        touched={touched}
                        errors={errors}
                    >
                        <TextInput
                            isRequired={values.updatePassword}
                            type="password"
                            id="notifier.splunk.httpToken"
                            name="notifier.splunk.httpToken"
                            value={values.notifier.splunk.httpToken}
                            onChange={onChange}
                            onBlur={handleBlur}
                            isDisabled={!isEditable || !values.updatePassword}
                            placeholder={
                                values.updatePassword
                                    ? ''
                                    : 'Currently-stored password will be used.'
                            }
                        />
                    </FormLabelGroup>
                    <FormLabelGroup
                        isRequired
                        label="HEC truncate limit"
                        fieldId="notifier.splunk.truncate"
                        helperText="Message length limit in bytes (characters)"
                        touched={touched}
                        errors={errors}
                    >
                        <TextInput
                            isRequired
                            type="number"
                            id="notifier.splunk.truncate"
                            value={values.notifier.splunk.truncate}
                            onChange={onChange}
                            onBlur={handleBlur}
                            isDisabled={!isEditable}
                        />
                    </FormLabelGroup>
                    <FormLabelGroup label="" fieldId="notifier.splunk.insecure" errors={errors}>
                        <Checkbox
                            label="Disable TLS certificate validation (insecure)"
                            id="notifier.splunk.insecure"
                            isChecked={values.notifier.splunk.insecure}
                            onChange={onChange}
                            onBlur={handleBlur}
                            isDisabled={!isEditable}
                        />
                    </FormLabelGroup>
                    <FormLabelGroup
                        label=""
                        fieldId="notifier.splunk.auditLoggingEnabled"
                        errors={errors}
                    >
                        <Checkbox
                            label="Enable audit logging"
                            id="notifier.splunk.auditLoggingEnabled"
                            isChecked={values.notifier.splunk.auditLoggingEnabled}
                            onChange={onChange}
                            onBlur={handleBlur}
                            isDisabled={!isEditable}
                        />
                    </FormLabelGroup>
                    <FormLabelGroup
                        isRequired
                        label="Source type for alert"
                        fieldId="notifier.splunk.sourceTypes.alert"
                        touched={touched}
                        errors={errors}
                    >
                        <TextInput
                            isRequired
                            type="text"
                            id="notifier.splunk.sourceTypes.alert"
                            value={values.notifier.splunk.sourceTypes.alert}
                            onChange={onChange}
                            onBlur={handleBlur}
                            isDisabled={!isEditable}
                        />
                    </FormLabelGroup>
                    <FormLabelGroup
                        isRequired
                        label="Source type for audit"
                        fieldId="notifier.splunk.sourceTypes.audit"
                        touched={touched}
                        errors={errors}
                    >
                        <TextInput
                            isRequired
                            type="text"
                            id="notifier.splunk.sourceTypes.audit"
                            value={values.notifier.splunk.sourceTypes.audit}
                            onChange={onChange}
                            onBlur={handleBlur}
                            isDisabled={!isEditable}
                        />
                    </FormLabelGroup>
                </Form>
            </PageSection>
            {isEditable && (
                <IntegrationFormActions>
                    <FormSaveButton
                        onSave={onSave}
                        isSubmitting={isSubmitting}
                        isTesting={isTesting}
                        isDisabled={!dirty || !isValid}
                    >
                        Save
                    </FormSaveButton>
                    <FormTestButton
                        onTest={onTest}
                        isSubmitting={isSubmitting}
                        isTesting={isTesting}
                        isValid={isValid}
                    >
                        Test
                    </FormTestButton>
                    <FormCancelButton onCancel={onCancel}>Cancel</FormCancelButton>
                </IntegrationFormActions>
            )}
        </>
    );
}

export default SplunkIntegrationForm;