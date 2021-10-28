import React, { Component } from "react";
import './App.css';

import Header from './components/Header/Header';
import NavBar from './components/NavBar/NavBar';
import * as Blocks from './components/Blocks';

import ReactNotification from 'react-notifications-component'
import 'react-notifications-component/dist/theme.css'

class App extends Component {

  constructor() {
    super()

    this.state = {
      selectedBlock: 'DropRates',
      dropRatesInfo: []
    }
  }

  getNavPage = (val) => {
    this.setState({ selectedBlock: val })
  }

  renderRollBlock(selectedBlock) {
    const Block = Blocks[selectedBlock]

    return <Block />
  }

  render() {
    return (
      <div className="App">
        <ReactNotification />
        <Header />
        <NavBar returnNavPage={this.getNavPage} />
        {this.renderRollBlock(this.state.selectedBlock)}
      </div >
    );
  }
}

export default App;
