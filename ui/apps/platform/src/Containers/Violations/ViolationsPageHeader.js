import React, { useEffect, useState } from 'react';
import PropTypes from 'prop-types';
import { connect } from 'react-redux';
import { createStructuredSelector } from 'reselect';

import { actions as alertActions } from 'reducers/alerts';
import { selectors } from 'reducers';
import { fetchAlerts, fetchAlertCount } from 'services/AlertsService';

import PageHeader from 'Components/PageHeader';
import ReduxSearchInput from 'Containers/Search/ReduxSearchInput';

import { pageSize } from 'Components/TableV2';
import { checkForPermissionErrorMessage } from 'utils/permissionUtils';

function runAfter5Seconds(fn) {
    return new Promise(() => {
        setTimeout(fn, 5000);
    });
}

function ViolationsPageHeader({
    isViewFiltered,
    setIsViewFiltered,
    currentPage,
    sortOption,
    setCurrentPageAlerts,
    setCurrentPageAlertsErrorMessage,
    setAlertCount,
    setSelectedAlertId,
    currentPageAlerts,
    selectedAlertId,
    searchOptions,
    searchModifiers,
    searchSuggestions,
    setSearchOptions,
    setSearchModifiers,
    setSearchSuggestions,
}) {
    // To handle page/count refreshing.
    const [pollEpoch, setPollEpoch] = useState(0);
    const [isFetching, setIsFetching] = useState(false);

    // Update the isViewFiltered and the value of the selectedAlertId based on changes in search options.
    const hasExecutableFilter =
        searchOptions.length && !searchOptions[searchOptions.length - 1].type;
    const hasNoFilter = !searchOptions.length;

    if (hasExecutableFilter && !isViewFiltered) {
        setIsViewFiltered(true);
    } else if (hasNoFilter && isViewFiltered) {
        setIsViewFiltered(false);
    }
    if (
        hasExecutableFilter &&
        selectedAlertId &&
        !currentPageAlerts.find((alert) => alert.id === selectedAlertId)
    ) {
        setSelectedAlertId(null);
    }

    // When any of the deps to this effect change, we want to reload the alerts and count.
    useEffect(() => {
        if (
            !isFetching &&
            (!searchOptions.length || !searchOptions[searchOptions.length - 1].type)
        ) {
            // Get the alerts that match the search request for the current page.
            setCurrentPageAlertsErrorMessage('');
            setIsFetching(true);
            fetchAlerts(searchOptions, sortOption, currentPage, pageSize)
                .then((alerts) => {
                    setCurrentPageAlerts(alerts);
                })
                .catch((error) => {
                    setCurrentPageAlerts([]);
                    const parsedMessage = checkForPermissionErrorMessage(error);
                    setCurrentPageAlertsErrorMessage(parsedMessage);
                })
                .finally(() => {
                    setIsFetching(false);
                });
            // Get the total count of alerts that match the search request.
            fetchAlertCount(searchOptions)
                .then(setAlertCount)
                .catch((error) => {
                    setCurrentPageAlerts([]);
                    const parsedMessage = checkForPermissionErrorMessage(error);
                    setCurrentPageAlertsErrorMessage(parsedMessage);
                });
        }

        // We will update the poll epoch after 5 seconds to force a refresh.
        runAfter5Seconds(() => {
            setPollEpoch(pollEpoch + 1);
        });
    }, [
        searchOptions,
        currentPage,
        sortOption,
        pollEpoch,
        setCurrentPageAlerts,
        setCurrentPageAlertsErrorMessage,
        setAlertCount,
    ]);

    // Render.
    const subHeader = isViewFiltered ? 'Filtered view' : 'Default view';
    const defaultOption = searchModifiers.find((x) => x.value === 'Deployment:');
    return (
        <PageHeader header="Violations" subHeader={subHeader}>
            <ReduxSearchInput
                className="w-full"
                id="alerts"
                searchOptions={searchOptions}
                searchModifiers={searchModifiers}
                searchSuggestions={searchSuggestions}
                setSearchOptions={setSearchOptions}
                setSearchModifiers={setSearchModifiers}
                setSearchSuggestions={setSearchSuggestions}
                defaultOption={defaultOption}
                autoCompleteCategories={['ALERTS']}
            />
        </PageHeader>
    );
}

ViolationsPageHeader.propTypes = {
    isViewFiltered: PropTypes.bool.isRequired,
    setIsViewFiltered: PropTypes.func.isRequired,
    currentPage: PropTypes.number.isRequired,
    sortOption: PropTypes.shape({}).isRequired,
    currentPageAlerts: PropTypes.arrayOf(PropTypes.object),
    setCurrentPageAlerts: PropTypes.func.isRequired,
    setCurrentPageAlertsErrorMessage: PropTypes.func.isRequired,
    setAlertCount: PropTypes.func.isRequired,
    setSelectedAlertId: PropTypes.func.isRequired,
    selectedAlertId: PropTypes.string,

    searchOptions: PropTypes.arrayOf(PropTypes.object).isRequired,
    searchModifiers: PropTypes.arrayOf(PropTypes.object).isRequired,
    searchSuggestions: PropTypes.arrayOf(PropTypes.object).isRequired,
    setSearchOptions: PropTypes.func.isRequired,
    setSearchModifiers: PropTypes.func.isRequired,
    setSearchSuggestions: PropTypes.func.isRequired,
};

ViolationsPageHeader.defaultProps = {
    currentPageAlerts: [],
    selectedAlertId: null,
};

const mapStateToProps = createStructuredSelector({
    searchOptions: selectors.getAlertsSearchOptions,
    searchModifiers: selectors.getAlertsSearchModifiers,
    searchSuggestions: selectors.getAlertsSearchSuggestions,
});

const mapDispatchToProps = {
    setSearchOptions: alertActions.setAlertsSearchOptions,
    setSearchModifiers: alertActions.setAlertsSearchModifiers,
    setSearchSuggestions: alertActions.setAlertsSearchSuggestions,
};

export default connect(mapStateToProps, mapDispatchToProps)(ViolationsPageHeader);