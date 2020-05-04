import React, { Component } from "react";
import ReactDOM from "react-dom";

class App extends Component {
  constructor() {
    super();
    this.state = {
      serverResponse: ""
    };
  }

  componentDidMount() {
    fetch("http://localhost:5000/", {
      method: "GET",
      mode: "cors",
      headers: {
        "Content-Type": "application/json"
      }
    })
      .then(res => res.json())
      .then(res => {
        this.setState({ serverResponse: res.text });
      });
  }

  render() {
    const { serverResponse } = this.state;
    return (
      <div>
        <div>This is Client 6003. Hello Gadhu! Server do you read? ...</div>
        <br />
        <div>{serverResponse}</div>
      </div>
    );
  }
}

ReactDOM.render(<App />, document.getElementById("app"));
