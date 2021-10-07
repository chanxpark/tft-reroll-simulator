import React, { Component } from "react";
import './App.css';

import Header from './components/Header/Header';
import NavBar from './components/NavBar/NavBar';
import * as Blocks from './components/Blocks';

class App extends Component {

  constructor() {
    super()

    this.state = {
      selectedBlock: 'DropRates',
      dropRatesInfo: []
    }
  }

  async send() {
    const response = await fetch('/api');

    console.log(response.json());
  }

  getNavPage = (val) => {
    this.setState({ selectedBlock: val })
  }

  renderRollBlock(selectedBlock) {
    console.log(selectedBlock)
    const Block = Blocks[selectedBlock]

    return <Block />
  }


  render() {
    return (
      <div className="App">
        <Header />
        <NavBar returnNavPage={this.getNavPage} />
        {this.renderRollBlock(this.state.selectedBlock)}
      </div >
    );
  }
}

export default App;
