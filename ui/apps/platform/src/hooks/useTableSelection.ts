import React from 'react';

export type UseTableSelection = {
    selected: boolean[];
    allRowsSelected: boolean;
    hasSelections: boolean;
    onSelect: (
        event: React.FormEvent<HTMLInputElement>,
        isSelected: boolean,
        rowId: number
    ) => void;
    onSelectAll: (event: React.FormEvent<HTMLInputElement>, isSelected: boolean) => void;
    onClearAll: () => void;
    getSelectedIds: () => string[];
};

type Base = {
    id: string;
};

function useTableSelection<T extends Base>(data: T[]): UseTableSelection {
    const [allRowsSelected, setAllRowsSelected] = React.useState(false);
    const [selected, setSelected] = React.useState(data.map(() => false));
    const hasSelections = selected.some((sel) => sel === true);

    React.useEffect(() => {
        setSelected(data.map(() => false));
    }, [data]);

    const onClearAll = () => {
        setSelected(data.map(() => false));
        setAllRowsSelected(false);
    };

    const onSelect = (event, isSelected: boolean, rowId: number) => {
        setSelected(
            selected.map((sel: boolean, index: number) => (index === rowId ? isSelected : sel))
        );
        if (!isSelected && allRowsSelected) {
            setAllRowsSelected(false);
        } else if (isSelected && !allRowsSelected) {
            let allSelected = true;
            for (let i = 0; i < selected.length; i += 1) {
                if (i !== rowId) {
                    if (!selected[i]) {
                        allSelected = false;
                    }
                }
            }
            if (allSelected) {
                setAllRowsSelected(true);
            }
        }
    };

    function onSelectAll(event, isSelected: boolean) {
        setAllRowsSelected(isSelected);
        setSelected(selected.map(() => isSelected));
    }

    function getSelectedIds() {
        const ids: string[] = [];
        for (let i = 0; i < selected.length; i += 1) {
            if (selected[i]) {
                ids.push(data[i].id);
            }
        }
        return ids;
    }

    return {
        selected,
        allRowsSelected,
        hasSelections,
        onSelect,
        onSelectAll,
        onClearAll,
        getSelectedIds,
    };
}

export default useTableSelection;