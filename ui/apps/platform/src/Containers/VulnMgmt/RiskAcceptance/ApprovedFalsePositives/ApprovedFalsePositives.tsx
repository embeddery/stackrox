import React, { ReactElement } from 'react';

import usePagination from 'hooks/patternfly/usePagination';
import queryService from 'utils/queryService';
import useSearch from 'hooks/useSearch';
import useVulnerabilityRequests from '../useVulnerabilityRequests';
import ApprovedFalsePositivesTable from './ApprovedFalsePositivesTable';

function ApprovedFalsePositives(): ReactElement {
    const { searchFilter, setSearchFilter } = useSearch();

    let modifiedSearchObject = { ...searchFilter };
    modifiedSearchObject = {
        ...modifiedSearchObject,
        'Expired Request': 'false',
        'Requested Vulnerability State': 'FALSE_POSITIVE',
        'Request Status': 'APPROVED',
    };
    const requestID = modifiedSearchObject['Request ID'];
    delete modifiedSearchObject['Request ID'];
    const query = queryService.objectToWhereClause(modifiedSearchObject);

    const { page, perPage, onSetPage, onPerPageSelect } = usePagination();
    const { isLoading, data, refetchQuery } = useVulnerabilityRequests({
        query,
        requestID,
        pagination: {
            limit: perPage,
            offset: (page - 1) * perPage,
            sortOption: {
                field: 'Last Updated',
                reversed: false,
            },
        },
    });

    const rows = data?.vulnerabilityRequests || [];
    const itemCount = data?.vulnerabilityRequestsCount || 0;

    return (
        <ApprovedFalsePositivesTable
            rows={rows}
            updateTable={refetchQuery}
            isLoading={isLoading}
            itemCount={itemCount}
            searchFilter={searchFilter}
            setSearchFilter={setSearchFilter}
            page={page}
            perPage={perPage}
            onSetPage={onSetPage}
            onPerPageSelect={onPerPageSelect}
        />
    );
}

export default ApprovedFalsePositives;