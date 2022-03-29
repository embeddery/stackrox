import { useCallback, useRef } from 'react';
import { useLocation, useHistory } from 'react-router-dom';
import isEqual from 'lodash/isEqual';
import { getQueryObject, getQueryString } from 'utils/queryStringUtils';

type QueryValue = undefined | string | string[] | qs.ParsedQs | qs.ParsedQs[];

/**
 * Hook to handle reading and writing of a piece of state in the page's URL query parameters.
 *
 * The return value of this hook follows the `useState` convention, returning a 2-length
 * array where the first item is the state value of the URL parameter and the second
 * value is a setter function to change that value.
 *
 * Both the returned state and setter function maintain referential equality across
 * calls as long as the state in the URL does not change.
 *
 * @param keyPrefix The key value of the url parameter to manage
 * @param defaultValue A default value to use when the parameter is not available in the URL
 *
 * @returns [value, setterFn]
 */
function useUrlParameter<T extends QueryValue>(
    keyPrefix: string,
    defaultValue: T
): [T, (newValue: T) => void] {
    const history = useHistory();
    const location = useLocation();
    // We use an internal Ref here so that calling code that depends on the
    // value of returned by this hook can detect updates. e.g. When used in the
    // dependency array of a `useEffect`.
    const internalValue = useRef<T>(defaultValue);
    // memoize the setter function to retain referential equality as long
    // as the URL parameters do not change
    const setValue = useCallback(
        (newValue: T) => {
            const previousQuery = getQueryObject(location.search) || {};
            // Merge the parameter managed by this hook with the rest of the URL parameters
            const newQueryString = getQueryString({
                ...previousQuery,
                [keyPrefix]: newValue,
            });
            history.replace({
                search: newQueryString,
            });
        },
        [keyPrefix, history, location.search]
    );

    const nextValue =
        getQueryObject<{ [s: string]: T }>(location.search)[keyPrefix] || defaultValue;

    // If the search filter has changed, replace the object reference.
    if (!isEqual(internalValue.current, nextValue)) {
        internalValue.current = nextValue;
    }

    return [internalValue.current, setValue];
}

export default useUrlParameter;
