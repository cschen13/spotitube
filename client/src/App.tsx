import * as React from "react";
import "./App.css";
import Home from "./Home/Home";
import styled from "styled-components";

class App extends React.Component {
  public render() {
    return (
      <div className="ui container">
        <Home />
      </div>
    );
  }
}

export default App;
