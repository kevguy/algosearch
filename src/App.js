import React from 'react';
import { BrowserRouter as Router, Route, Switch } from 'react-router-dom';

// Page imports
import Home from './pages/Home';
import Address from './pages/Address';
import AddressTransaction from './pages/AddressTransaction';
import Block from './pages/Block';
import Blocks from './pages/Blocks';
import Transaction from './pages/Transaction';
import Transactions from './pages/Transactions';
import Analytics from './pages/Analytics';
import RichList from './pages/RichList';
import Dev from './pages/Dev';

function App() {
  return (
    <div className="App">
      <Router>
        <Route component={scrollRestoration} />
        <Switch>
          <Route path="/" exact component={Home} />
          <Route path="/address/:address" component={props => <Address {...props} key={Math.ceil(Math.random() * 10)}/>} />
          <Route path="/addresstx/:address" component={AddressTransaction} />
          <Route path="/blocks" exact component={Blocks} />
          <Route path="/block/:blocknum" component={props => <Block {...props} key={Math.ceil(Math.random() * 10)}/>} />
          <Route path="/transactions" exact component={Transactions} />
          <Route path="/tx/:txid" component={Transaction} />
          <Route path="/analytics" exact component={Analytics} />
          <Route path="/richlist" exact component={RichList} />
          <Route path="/dev" exact component={Dev} />
          <Route component={Home} />
        </Switch>
      </Router>
    </div>
  );
}

class scrollRestoration extends React.Component {
  componentDidUpdate(prevProps) {
    if (
      this.props.location.pathname !== prevProps.location.pathname
    ) {
      window.scrollTo(0, 0);
    }
  }

  render() {
    return null;
  }
}

export default App;
