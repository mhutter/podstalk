import React, { Component } from 'react'
import Pod from './Pod';

class PodList extends Component {
  constructor(props) {
    super(props);
    this.state = { pods: [] };
  }

  componentDidMount() {
    fetch('/api/pods')
      .then(res => res.json())
      .then(pods => {
        pods.sort((a, b) => a.name > b.name ? 1 : -1)
        this.setState({ pods })
      })
      .catch(console.error)
  }

  render() {
    return (
      <div className="pod_list">
        <h1>Pods in my neighborhood</h1>
        {this.state.pods.map(pod => <Pod key={pod.uid} pod={pod} />)}
      </div>
    )
  }
}

export default PodList
