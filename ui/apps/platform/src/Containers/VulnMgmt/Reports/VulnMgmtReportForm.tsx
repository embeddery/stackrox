/* eslint-disable no-void */
import React, { useState, ReactElement } from 'react';
import { useHistory } from 'react-router-dom';
import {
    ActionList,
    ActionListItem,
    Button,
    ButtonVariant,
    Divider,
    Form,
    Grid,
    GridItem,
    PageSection,
    PageSectionVariants,
    Select,
    SelectOption,
    SelectVariant,
    Text,
    TextArea,
    TextInput,
    TextVariants,
    Title,
} from '@patternfly/react-core';
import { useFormik } from 'formik';

import SelectSingle from 'Components/SelectSingle';
import FormMessage, { FormResponseMessage } from 'Components/PatternFly/FormMessage';
import RepeatScheduleDropdown from 'Components/PatternFly/RepeatScheduleDropdown';
import DayPickerDropdown from 'Components/PatternFly/DayPickerDropdown';
import FormLabelGroup from 'Components/PatternFly/FormLabelGroup';
import useMultiSelect from 'hooks/useMultiSelect';
import { saveReport } from 'services/ReportsService';
import { ReportConfiguration } from 'types/report.proto';
import { getAxiosErrorMessage } from 'utils/responseErrorUtils';
import NotifierSelection from './Form/NotifierSelection';
import ResourceScopeSelection from './Form/ResourceScopeSelection';
import { getMappedFixability, getFixabilityConstantFromMap } from './VulnMgmtReport.utils';

export type VulnMgmtReportFormProps = {
    initialValues: ReportConfiguration;
    isEditable?: boolean;
};

function VulnMgmtReportForm({
    initialValues,
    isEditable = true,
}: VulnMgmtReportFormProps): ReactElement {
    const history = useHistory();
    const [message, setMessage] = useState<FormResponseMessage>(null);
    const formik = useFormik<ReportConfiguration>({
        initialValues,
        onSubmit: (formValues) => {
            const response = onSave(formValues);
            return response;
        },
    });

    const {
        values,
        touched,
        errors,
        dirty,
        isValid,
        setFieldValue,
        handleBlur,
        submitForm,
        isSubmitting,
    } = formik;

    const mappedFixabilityValues = getMappedFixability(values.vulnReportFilters.fixability);

    const {
        isOpen: isFixabilitySelectOpen,
        onToggle: onToggleFixabilitySelect,
        onSelect: onSelectFixability,
    } = useMultiSelect(handleFixabilitySelect, mappedFixabilityValues);
    const {
        isOpen: isSeveritySelectOpen,
        onToggle: onToggleSeveritySelect,
        onSelect: onSelectSeverity,
    } = useMultiSelect(handleSeveritySelect, values.vulnReportFilters.severities, false);

    async function onSave(data) {
        let responseData;
        try {
            // eslint-disable-next-line @typescript-eslint/no-unused-vars
            responseData = await saveReport(data);

            history.goBack();
        } catch (error) {
            setMessage({ message: getAxiosErrorMessage(error), isError: true });
        }
    }

    function cancelEdit() {
        history.goBack();
    }

    function onChange(value, event) {
        return setFieldValue(event.target.id, value);
    }

    function handleFixabilitySelect(selection) {
        const fixabilityConstant = getFixabilityConstantFromMap(selection);
        void setFieldValue('vulnReportFilters.fixability', fixabilityConstant);
    }

    function handleSeveritySelect(selection) {
        void setFieldValue('vulnReportFilters.severities', selection);
    }

    function onSinceLastReportChange(_id, selection) {
        void setFieldValue('vulnReportFilters.sinceLastReport', selection === 'true');
    }

    function onScheduledRepeatChange(_id, selection) {
        // zero out the days selected list if changing interval type
        if (selection !== values.schedule.intervalType) {
            void setFieldValue('schedule.interval.days', []);
        }

        void setFieldValue('schedule.intervalType', selection);
    }

    function onScheduledDaysChange(id, selection) {
        void setFieldValue(id, selection);
    }

    return (
        <>
            <PageSection variant={PageSectionVariants.light} isFilled hasOverflowScroll>
                <FormMessage message={message} />
                <Form>
                    <Grid hasGutter>
                        <GridItem span={8}>
                            <Grid hasGutter>
                                <GridItem span={6}>
                                    <FormLabelGroup
                                        label="Report name"
                                        isRequired
                                        fieldId="name"
                                        touched={touched}
                                        errors={errors}
                                    >
                                        <TextInput
                                            isRequired
                                            type="text"
                                            id="name"
                                            value={values.name}
                                            onChange={onChange}
                                            onBlur={handleBlur}
                                            isDisabled={!isEditable}
                                        />
                                    </FormLabelGroup>
                                </GridItem>
                                <GridItem span={3}>
                                    <RepeatScheduleDropdown
                                        label="Repeat report…"
                                        isRequired
                                        fieldId="schedule.intervalType"
                                        value={values.schedule.intervalType}
                                        handleSelect={onScheduledRepeatChange}
                                    />
                                </GridItem>
                                <GridItem span={3}>
                                    <DayPickerDropdown
                                        label="On…"
                                        isRequired
                                        fieldId={
                                            values.schedule.intervalType === 'WEEKLY'
                                                ? 'schedule.daysOfWeek.days'
                                                : 'schedule.daysOfMonth.days'
                                        }
                                        value={
                                            values.schedule.intervalType === 'WEEKLY'
                                                ? values?.schedule?.daysOfWeek?.days || []
                                                : values?.schedule?.daysOfMonth?.days || []
                                        }
                                        handleSelect={onScheduledDaysChange}
                                        intervalType={values.schedule.intervalType}
                                    />
                                </GridItem>
                                <GridItem span={12}>
                                    <FormLabelGroup
                                        label="Description"
                                        fieldId="description"
                                        touched={touched}
                                        errors={errors}
                                    >
                                        <TextArea
                                            type="text"
                                            id="description"
                                            value={values.description}
                                            onChange={onChange}
                                            onBlur={handleBlur}
                                            isDisabled={!isEditable}
                                        />
                                    </FormLabelGroup>
                                </GridItem>
                                <GridItem span={4}>
                                    <FormLabelGroup
                                        isRequired
                                        label="CVE fixability type"
                                        fieldId="vulnReportFilters.fixability"
                                        touched={touched}
                                        errors={errors}
                                    >
                                        <Select
                                            variant={SelectVariant.checkbox}
                                            aria-label="Select CVE fixibility"
                                            onToggle={onToggleFixabilitySelect}
                                            onSelect={onSelectFixability}
                                            selections={mappedFixabilityValues}
                                            isOpen={isFixabilitySelectOpen}
                                            placeholderText={
                                                mappedFixabilityValues.length > 0
                                                    ? 'Fixable states selected'
                                                    : 'Select CVE fixibility'
                                            }
                                        >
                                            <SelectOption key="FIXABLE" value="FIXABLE">
                                                Fixable
                                            </SelectOption>
                                            <SelectOption key="NOT_FIXABLE" value="NOT_FIXABLE">
                                                Unfixable
                                            </SelectOption>
                                        </Select>
                                    </FormLabelGroup>
                                </GridItem>
                                <GridItem span={4}>
                                    <FormLabelGroup
                                        isRequired
                                        label="Show vulnerabilities"
                                        fieldId="vulnReportFilters.sinceLastReport"
                                        touched={touched}
                                        errors={errors}
                                    >
                                        <SelectSingle
                                            id="vulnReportFilters.sinceLastReport"
                                            value={values.vulnReportFilters.sinceLastReport.toString()}
                                            handleSelect={onSinceLastReportChange}
                                            isDisabled={false}
                                        >
                                            <SelectOption value="true">
                                                since last successful report
                                            </SelectOption>
                                            <SelectOption value="false">
                                                all vulnerabilities
                                            </SelectOption>
                                        </SelectSingle>
                                    </FormLabelGroup>
                                </GridItem>
                                <GridItem span={4}>
                                    <FormLabelGroup
                                        isRequired
                                        label="CVE severities"
                                        fieldId="vulnReportFilters.severities"
                                        touched={touched}
                                        errors={errors}
                                    >
                                        <Select
                                            variant={SelectVariant.checkbox}
                                            aria-label="Select CVE severities"
                                            onToggle={onToggleSeveritySelect}
                                            onSelect={onSelectSeverity}
                                            selections={values.vulnReportFilters.severities}
                                            isOpen={isSeveritySelectOpen}
                                            placeholderText={
                                                values.vulnReportFilters.severities.length > 0
                                                    ? 'Severities selected'
                                                    : 'Select CVE severities'
                                            }
                                        >
                                            <SelectOption value="CRITICAL_VULNERABILITY_SEVERITY">
                                                Critical
                                            </SelectOption>
                                            <SelectOption value="IMPORTANT_VULNERABILITY_SEVERITY">
                                                Important
                                            </SelectOption>
                                            <SelectOption value="MODERATE_VULNERABILITY_SEVERITY">
                                                Medium
                                            </SelectOption>
                                            <SelectOption value="LOW_VULNERABILITY_SEVERITY">
                                                Low
                                            </SelectOption>
                                        </Select>
                                    </FormLabelGroup>
                                </GridItem>
                                <GridItem span={12}>
                                    <ResourceScopeSelection
                                        scopeId={values.scopeId}
                                        setFieldValue={setFieldValue}
                                    />
                                </GridItem>
                            </Grid>
                        </GridItem>
                        <GridItem span={4}>
                            <div>
                                <Title headingLevel="h2" className="pf-u-mb-xs">
                                    Notification method and distribution
                                </Title>
                                <Text component={TextVariants.p} className="pf-u-mb-md">
                                    Schedule reports across the organization by defining a
                                    notification method and distribution list for the report
                                </Text>
                                <NotifierSelection
                                    notifierId={values.emailConfig.notifierId}
                                    mailingLists={values.emailConfig.mailingLists}
                                    setFieldValue={setFieldValue}
                                    handleBlur={handleBlur}
                                />
                            </div>
                        </GridItem>
                    </Grid>
                </Form>
            </PageSection>
            <Divider component="div" />
            <PageSection variant={PageSectionVariants.light} style={{ flexGrow: 0 }}>
                <ActionList>
                    <ActionListItem>
                        <Button
                            variant={ButtonVariant.primary}
                            onClick={submitForm}
                            data-testid="create-btn"
                            isDisabled={!dirty || !isValid || isSubmitting}
                            isLoading={isSubmitting}
                        >
                            {values.id ? 'Save' : 'Create'}
                        </Button>
                    </ActionListItem>
                    <ActionListItem>
                        <Button variant={ButtonVariant.link} onClick={cancelEdit}>
                            Cancel
                        </Button>
                    </ActionListItem>
                </ActionList>
            </PageSection>
        </>
    );
}

export default VulnMgmtReportForm;