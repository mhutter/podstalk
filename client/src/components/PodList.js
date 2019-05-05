/* global WebSocket */
import React, { Component } from 'react'
import Pod from './Pod'

const wsURL = () => {
  if (process.env.NODE_ENV !== 'production') {
    return 'ws://localhost:3001/api/ws'
  }

  const url = new URL('./api/ws', window.location.href)
  url.protocol = window.location.protocol === 'http:' ? 'ws:' : 'wss:'
  return url.href
}

class PodList extends Component {
  constructor (props) {
    super(props)
    this.state = { pods: {} }
    this.connect = this.connect.bind(this)
    this.handleEvent = this.handleEvent.bind(this)
  }

  handleEvent ({ data }) {
    const {type, pod} = JSON.parse(data)
    const pods = this.state.pods || {}

    switch (type) {
      case 'ADDED':
      case 'MODIFIED':
        pods[pod.uid] = pod
        break

      case 'DELETED':
        delete pods[pod.uid]
        break

      default:
        console.error('Unknown event type:', type)
    }

    this.setState({ pods })
  }

  connect () {
    const ws = new WebSocket(wsURL())

    // After opening a new connection, we'll get ADDED events for all
    // existing pods. However, we may have missed some DELETED events so
    // make sure we have no orphaned pods in our state.
    ws.onopen = () => {
      this.setState({ pods: {} })
    }

    // Reconnect after 2 seconds. "onclose" is also fired on errors,
    // even during initial connection.
    ws.onclose = () => {
      setTimeout(this.connect, 2000)
    }
    ws.onmessage = this.handleEvent
  }

  componentDidMount () {
    this.connect()
  }

  render () {
    const pods = Object.values(this.state.pods)
    pods.sort((a, b) => a.name > b.name ? 1 : -1)

    return (
      <div className='pod_list'>
        <h1>Pods in my neighborhood</h1>
        {pods.map(pod => <Pod key={pod.uid} pod={pod} />)}
      </div>
    )
  }
}

export default PodList
