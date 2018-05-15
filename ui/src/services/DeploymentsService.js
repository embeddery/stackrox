import axios from 'axios';
import queryString from 'query-string';
import { normalize } from 'normalizr';

import { deployment as deploymentSchema } from './schemas';

const deploymentsUrl = '/v1/deployments';

/**
 * Fetches list of registered deployments.
 *
 * @returns {Promise<Object[], Error>} fulfilled with array of deployments (as defined in .proto)
 */
export function fetchDeployments(filters) {
    const params = queryString.stringify({
        ...filters
    });
    return axios
        .get(`${deploymentsUrl}?${params}`)
        .then(response => ({ response: normalize(response.data.deployments, [deploymentSchema]) }));
}

/**
 * Fetches a deployment by its ID.
 *
 * @returns {Promise<Object, Error>} fulfilled with a deployment object (as defined in .proto)
 */
export function fetchDeployment(id) {
    if (!id) throw new Error('Deployment ID must be specified');
    return axios
        .get(`${deploymentsUrl}/${id}`)
        .then(response => ({ response: normalize(response.data, deploymentSchema) }));
}
