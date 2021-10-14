import React, { Component } from "react";
import { CreateTable } from './Table.jsx'
import { store } from 'react-notifications-component';
import "./RollByLevel.scss";

class RollByLevel extends Component {

    constructor() {
        super()

        this.state = {
            level: 1,
            rollNums: 1,
            rollResult: [],
            showTable: false,
            submit: false
        }
        this.handleSubmit = this.handleSubmit.bind(this)
    }

    async handleSubmit(e) {
        e.preventDefault()
        const level = this.state.level
        const rolls = this.state.rollNums
        const response = await fetch('/api/roll?level=' + level + '&rolls=' + rolls);
        const responseJson = await response.json()

        if (response.ok) {
            this.setState({
                rollResult: responseJson,
                submit: true,
                showTable: true
            })
        } else {
            console.log(responseJson["message"])
            store.addNotification({
                title: "Error",
                message: responseJson["message"],
                type: "warning",
                insert: "top",
                container: "top-right",
                dismiss: {
                    duration: 5000,
                    onScreen: true,
                    showIcon: true
                }
            })
        }
    }

    render() {
        return (
            <section className="RollBlock">
                <div className="rollOptions">
                    <form onSubmit={this.handleSubmit}>
                        <label>Level:
                            <select id="level" value={this.state.level} onChange={(e) => this.setState({ level: e.target.value })}>
                                <option value="1">1</option>
                                <option value="2">2</option>
                                <option value="3">3</option>
                                <option value="4">4</option>
                                <option value="5">5</option>
                                <option value="6">6</option>
                                <option value="7">7</option>
                                <option value="8">8</option>
                                <option value="9">9</option>
                            </select>
                        </label>
                        <label>Number of Rolls:
                            <input type="text" id="rollNums" value={this.state.rollNums} onChange={(e) => this.setState({ rollNums: e.target.value })}></input>
                            <input type="Submit" name="Roll"></input>
                        </label>
                    </form>
                </div >
                <div className={this.state.submit ? 'rollResults active' : 'rollResults'}>
                    < CreateTable data={this.state.rollResult} />
                </div>
            </section>
        )
    }
}

export default RollByLevel;