import React, { useCallback, useState } from 'react'
import championsData from '../../../assets/champions.json'
import "./RollByChampion.scss"

export function RollByChampion() {
    const [champions, setChampions] = useState(championsData)

    function importAll(r) {
        let images = {};
        r.keys().forEach((item, index) => { images[item.replace('./', '')] = r(item); });
        return images
    }

    const images = importAll(require.context('../../../assets/champions', false, /\.png$/))

    const loadChampList = useCallback(() => {
        function createChampWrapper(name, cost, id) {
            return `<div class="champion-wrapper c${cost}">
                        <img class="champion-icon" src="${images[`${id}.png`].default}" alt="${name}"></img>
                    </div>`
        }

        var results = ""
        champions.map((champion) => (results += createChampWrapper(champion.name, champion.cost, champion.championId)))

        return results
    }, [])

    return (
        <div className="champion-list" dangerouslySetInnerHTML={{ __html: loadChampList() }}>
        </div >

    )
}
