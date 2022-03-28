import useUrlParameter from './useUrlParameter';

function useSearch(keyPrefix = 'search') {
    const [searchFilter, setSearchFilter] = useUrlParameter(keyPrefix, {});
    return { searchFilter, setSearchFilter };
}

export default useSearch;
