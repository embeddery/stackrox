import { useCallback, useRef } from 'react';
import { useLocation, useHistory } from 'react-router-dom';
import isEqual from 'lodash/isEqual';
import { getQueryObject, getQueryString } from 'utils/queryStringUtils';

type QueryValue = undefined | string | string[] | qs.ParsedQs | qs.ParsedQs[];

function useUrlParameter<T extends QueryValue>(
    keyPrefix: string,
    defaultValue: T
): [T, (newStateValue: T) => void] {
    const history = useHistory();
    const location = useLocation();
    // We use an internal Ref here so that calling code that depends on the
    // value of returned by this hook can detect updates. e.g. When used in the
    // dependency array of a `useEffect`.
    const internalValue = useRef<T>(defaultValue);
    const setValue = useCallback(
        (newValue: T) => {
            const previousQuery = getQueryObject(location.search) || {};
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
