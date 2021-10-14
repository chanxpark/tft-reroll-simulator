import React, { Component } from "react";
import "./NavBar.scss";

class NavBar extends Component {

    constructor() {
        super()

        this.state = {
            clicked: false,
            item: false
        }
    }

    handleClick = (e) => {
        this.setState({ clicked: !this.state.clicked })
        if (e.target.parentElement.id === "DropRatesSelect") {
            this.props.returnNavPage("DropRates")
        } else if (e.target.parentElement.id === "RollByLevelSelect") {
            this.props.returnNavPage("RollByLevel")
        } else if (e.target.parentElement.id === "RollByChampionSelect") {
            this.props.returnNavPage("RollByChampion")
        }
    }

    render() {
        return (
            <nav className="Navbar">
                <div className="menu-icon"
                    style={{ backgroundColor: this.state.hoverColor }}
                    onClick={this.handleClick}
                >
                    <i className={this.state.clicked ? 'fas fa-times' : 'fas fa-bars'}></i>
                </div>
                <ul className={this.state.clicked ? 'nav-menu active' : 'nav-menu'}>
                    <li id="DropRatesSelect" className={this.state.clicked ? 'nav-links active' : 'nav-links'}
                        onClick={this.handleClick}
                    >
                        <a>Drop Rates</a>
                    </li>
                    <li id="RollByLevelSelect" className={this.state.clicked ? 'nav-links active' : 'nav-links'}
                        onClick={this.handleClick}
                    >
                        <a>Roll By Level</a>
                    </li>
                </ul>
            </nav >
        )
    }
}

export default NavBar;