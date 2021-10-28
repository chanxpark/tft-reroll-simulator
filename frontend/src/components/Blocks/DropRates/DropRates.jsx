import React, { useState, useEffect } from 'react'
import { useTable, useSortBy } from 'react-table'
import "./DropRates.scss";

const columns = [
    {
        Header: "Level",
        accessor: "Level",
        sortType: "alphanumeric",
    },
    {
        Header: "1-Cost",
        accessor: "OneCost",
        sortType: "alphanumeric",
    },
    {
        Header: "2-Cost",
        accessor: "TwoCost",
        sortType: "alphanumeric",
    },
    {
        Header: "3-Cost",
        accessor: "ThreeCost",
        sortType: "alphanumeric",
    },
    {
        Header: "4-Cost",
        accessor: "FourCost",
        sortType: "alphanumeric",
    },
    {
        Header: "5-Cost",
        accessor: "FiveCost",
        sortType: "alphanumeric",
    }
]

export function DropRates() {

    const [dropRates, setDropRates] = useState([])

    useEffect(() => {
        fetch(
            `/api/droprates`, {
            method: "GET"
        })
            .then(res => res.json())
            .then(response => {
                setDropRates(response)
            })
            .catch(error => console.log(error))
    }, [])

    const tableInstance = useTable(
        {
            columns,
            data: dropRates
        },
        useSortBy
    )

    const { getTableProps, getTableBodyProps, headerGroups, rows, prepareRow } = tableInstance


    return (
        <div className="DropRatesBlock">
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
        </div>
    )
}