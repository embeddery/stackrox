import { useCallback } from 'react';
import useUrlParameter from './useUrlParameter';

type UseUrlPaginationResult = {
    page: number;
    perPage: number;
    setPage: (page: number) => void;
    setPerPage: (perPage: number) => void;
};

function useUrlPagination(defaultPerPage = 20): UseUrlPaginationResult {
    const [page, setPageString] = useUrlParameter<string>('page', '1');
    const [perPage, setPerPageString] = useUrlParameter<string>('perPage', `${defaultPerPage}`);
    const setPage = useCallback((num: number) => setPageString(`${num}`), [setPageString]);
    const setPerPage = useCallback((num: number) => setPerPageString(`${num}`), [setPerPageString]);
    return {
        page: Number(page),
        perPage: Number(perPage),
        setPage,
        setPerPage,
    };
}

export default useUrlPagination;
