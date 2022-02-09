import React, { ReactElement, useEffect, useState } from 'react';
import { Alert, Flex, FlexItem, Spinner, Title, Divider, Button } from '@patternfly/react-core';
import { useFormikContext } from 'formik';

import { DryRunAlert, checkDryRun, startDryRun } from 'services/PoliciesService';
import { Policy } from 'types/policy.proto';
import { getAxiosErrorMessage } from 'utils/responseErrorUtils';

import { getServerPolicy } from '../../policies.utils';
import PolicyDetailContent from '../../Detail/PolicyDetailContent';
import PreviewViolations from './PreviewViolations';

import './ReviewPolicyForm.css';

type ReviewPolicyFormProps = {
    isBadRequest: boolean;
    policyErrorMessage: string;
    setIsBadRequest: (isBadRequest: boolean) => void;
    setIsValidOnServer: (isValidOnServer: boolean) => void;
    setPolicyErrorMessage: (message: string) => void;
};

function ReviewPolicyForm({
    isBadRequest,
    policyErrorMessage,
    setIsBadRequest,
    setIsValidOnServer,
    setPolicyErrorMessage,
}: ReviewPolicyFormProps): ReactElement {
    const { values } = useFormikContext<Policy>();

    const [showPolicyResults, setShowPolicyResults] = useState(true);
    const [isRunningDryRun, setIsRunningDryRun] = useState(false);
    const [jobIdOfDryRun, setJobIdOfDryRun] = useState('');
    const [counterToCheckDryRun, setCounterToCheckDryRun] = useState(0);
    const [checkDryRunErrorMessage, setCheckDryRunErrorMessage] = useState('');
    const [alertsFromDryRun, setAlertsFromDryRun] = useState<DryRunAlert[]>([]);

    // Start "dry run" job for preview of violations.
    useEffect(() => {
        setIsValidOnServer(false);
        setIsRunningDryRun(true);
        setPolicyErrorMessage('');
        setIsBadRequest(false);
        setCheckDryRunErrorMessage('');
        setAlertsFromDryRun([]);

        startDryRun(getServerPolicy(values))
            .then(({ data: { jobId } }) => {
                /*
                 * TODO after policiesSagas.js has been deleted:
                 * Replace ({ data: { jobId } }) with (jobId) above.
                 */
                setIsValidOnServer(true);
                setJobIdOfDryRun(jobId);
            })
            .catch((error) => {
                setIsRunningDryRun(false);
                setPolicyErrorMessage(getAxiosErrorMessage(error));
                if (error.response?.status === 400) {
                    setIsBadRequest(true);
                }
            });
    }, [setIsBadRequest, setIsValidOnServer, setPolicyErrorMessage, values]);

    // Poll "dry run" job for preview of violations.
    useEffect(() => {
        if (jobIdOfDryRun) {
            checkDryRun(jobIdOfDryRun)
                .then(({ data: { pending, result } }) => {
                    /*
                     * TODO after policiesSagas.js has been deleted:
                     * Replace ({ data: { pending, result } }) with ({ pending, result }) above.
                     */
                    if (pending) {
                        // To make another request, increment counterToCheckDryRun which is in useEffect dependencies.
                        setCounterToCheckDryRun((counter) => counter + 1);
                    } else {
                        setIsRunningDryRun(false);
                        setJobIdOfDryRun('');
                        setCounterToCheckDryRun(0);
                        setAlertsFromDryRun(result.alerts);
                    }
                })
                .catch((error) => {
                    setIsRunningDryRun(false);
                    setCheckDryRunErrorMessage(getAxiosErrorMessage(error));
                    setJobIdOfDryRun('');
                    setCounterToCheckDryRun(0);
                });
        }
    }, [jobIdOfDryRun, counterToCheckDryRun]);

    /*
     * flex_1 so columns have equal width.
     * alignSelfStretch so columns have equal height for border.
     */

    /* eslint-disable no-nested-ternary */
    return (
        <Flex>
            <Flex
                flex={{ default: 'flex_1' }}
                direction={{ default: 'column' }}
                alignSelf={{ default: 'alignSelfStretch' }}
                className="review-policy"
            >
                <Flex>
                    <FlexItem flex={{ default: 'flex_1' }}>
                        <Title headingLevel="h2">Review policy</Title>
                        <div>Review policy settings and violations.</div>
                    </FlexItem>
                    <FlexItem className="pf-u-pr-md" alignSelf={{ default: 'alignSelfCenter' }}>
                        <Button
                            variant="secondary"
                            onClick={() => setShowPolicyResults(!showPolicyResults)}
                        >
                            Policy results
                        </Button>
                    </FlexItem>
                </Flex>
                {policyErrorMessage && (
                    <Alert
                        title={isBadRequest ? 'Policy is invalid' : 'Policy request failure'}
                        variant="danger"
                        isInline
                    >
                        {policyErrorMessage}
                    </Alert>
                )}
                <Divider component="div" />
                <PolicyDetailContent policy={getServerPolicy(values)} isReview />
            </Flex>
            {showPolicyResults && (
                <>
                    <Divider component="div" isVertical />
                    <Flex
                        flex={{ default: 'flex_1' }}
                        direction={{ default: 'column' }}
                        alignSelf={{ default: 'alignSelfStretch' }}
                        className="preview-violations"
                    >
                        <Title headingLevel="h2">Preview violations</Title>
                        <div className="pf-u-mb-md pf-u-mt-sm">
                            The policy settings you have selected will generate violations for the
                            following deployments. Before you save the policy, verify that the
                            violations seem accurate.
                        </div>
                        {isRunningDryRun ? (
                            <Flex justifyContent={{ default: 'justifyContentCenter' }}>
                                <FlexItem>
                                    <Spinner isSVG />
                                </FlexItem>
                            </Flex>
                        ) : checkDryRunErrorMessage ? (
                            <Alert title="Violations request failure" variant="warning" isInline>
                                {checkDryRunErrorMessage}
                            </Alert>
                        ) : (
                            <PreviewViolations alertsFromDryRun={alertsFromDryRun} />
                        )}
                    </Flex>
                </>
            )}
        </Flex>
    );
    /* eslint-enable no-nested-ternary */
}

export default ReviewPolicyForm;