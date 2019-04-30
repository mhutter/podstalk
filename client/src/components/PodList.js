import React, { Component } from 'react'
import Pod from './Pod';

const pods = [
  {
    "name": "nginx-7db9fccd9b-jnljs",
    "namespace": "podstalk",
    "node": "minikube",
    "phase": "Pending",
    "uid": "a9375ce7-6b78-11e9-9b1f-080027970b01"
  },
  {
    "name": "nginx-7db9fccd9b-zq5zh",
    "namespace": "podstalk",
    "node": "minikube",
    "phase": "Running",
    "uid": "86b37ea4-6b78-11e9-9b1f-080027970b01"
  },
  {
    "name": "nginx-7db9fccd9b-r4c5k",
    "namespace": "podstalk",
    "node": "minikube",
    "phase": "Terminating",
    "uid": "a1de83a7-6b78-11e9-9b1f-080027970b01"
  }
]

class PodList extends Component {
  constructor(props) {
    super(props);
    this.state = { pods };
  }

  onComponentDidMount() {

  }

  render() {
    return (
      <div className="pod_list">
        <h1>Pods in my Namespace</h1>
        {this.state.pods.map(pod => <Pod key={pod.uid} pod={pod} />)}
      </div>
    )
  }
}

export default PodList
