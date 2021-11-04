import React from 'react';
import PropTypes from 'prop-types';
import { resolveAlerts } from 'services/AlertsService';
import pluralize from 'pluralize';
import Dialog from 'Components/Dialog';

function ResolveConfirmation({
    setDialogue,
    checkedAlertIds,
    setCheckedAlertIds,
    resolvableAlerts,
}) {
    function closeAndClear() {
        setDialogue(null);
        setCheckedAlertIds([]);
    }

    function resolveAlertsAction() {
        const resolveSelection = checkedAlertIds.filter((id) => resolvableAlerts.has(id));
        resolveAlerts(resolveSelection).then(closeAndClear, closeAndClear);
    }

    function close() {
        setDialogue(null);
    }

    const numSelectedRows = checkedAlertIds.reduce(
        (acc, id) => (resolvableAlerts.has(id) ? acc + 1 : acc),
        0
    );
    return (
        <Dialog
            isOpen
            text={`Are you sure you want to resolve ${numSelectedRows} ${pluralize(
                'violation',
                numSelectedRows
            )}?`}
            onConfirm={resolveAlertsAction}
            onCancel={close}
        />
    );
}

ResolveConfirmation.propTypes = {
    setDialogue: PropTypes.func.isRequired,
    checkedAlertIds: PropTypes.arrayOf(PropTypes.string).isRequired,
    setCheckedAlertIds: PropTypes.func.isRequired,
    resolvableAlerts: PropTypes.shape({
        has: PropTypes.func.isRequired,
    }).isRequired,
};

export default ResolveConfirmation;