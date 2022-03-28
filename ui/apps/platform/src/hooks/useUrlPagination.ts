import { useCallback } from 'react';
import useUrlParameter from './useUrlParameter';

function useUrlPagination(): [number, (n: number) => void] {
    const [page, setPageString] = useUrlParameter<string>('page', '1');
    const setPage = useCallback((num: number) => setPageString(`${num}`), [setPageString]);
    return [Number(page), setPage];
}

export default useUrlPagination;
