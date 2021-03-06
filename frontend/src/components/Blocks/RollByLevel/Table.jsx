import { useTable, useSortBy } from 'react-table'

const columns = [
    {
        Header: "Champion",
        accessor: "Champion",
        sortType: "alphanumeric",
    },
    {
        Header: "Cost",
        accessor: "Cost",
        sortType: "alphanumeric",
    },
    {
        Header: "Appearances",
        accessor: "Appearances",
        sortType: "alphanumeric",
    },
];

// {Champion: "Senna", Cost: 1, Appearances: 57}
export function CreateTable(props) {

    const tableInstance = useTable(
        {
            columns,
            data: props.data[0]
        },
        useSortBy
    )

    const { getTableProps, getTableBodyProps, headerGroups, rows, prepareRow } = tableInstance

    return (
        <table {...getTableProps()}>
            <thead>
                {headerGroups.map((headerGroup) => (
                    <tr {...headerGroup.getHeaderGroupProps()}>
                        {headerGroup.headers.map((column) => (
                            <th {...column.getHeaderProps(column.getSortByToggleProps())}>
                                {column.render('Header')}
                                <span> </span><span className={column.isSorted ? (column.isSortedDesc ? 'fas fa-sort-down' : 'fas fa-sort-up') : 'fas fa-sort'}></span>
                            </th>
                        ))}
                    </tr>
                ))}
            </thead>
            <tbody {...getTableBodyProps()}>
                {
                    rows.map(row => {
                        prepareRow(row)
                        return (
                            <tr {...row.getRowProps()}>
                                {row.cells.map((cell) => {
                                    return <td {...cell.getCellProps()}>{cell.render('Cell')}</td>
                                })}
                            </tr>
                        )
                    })
                }
            </tbody>
        </table>
    )
}