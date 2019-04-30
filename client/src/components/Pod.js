import React from 'react'
import './Pod.css'

const podPhase = phase => 
  <span className={`pod--phase-${phase}`}>{phase}</span>

const Pod = ({ pod }) => (
  <div className="pod">
    Pod <code>{pod.name}</code> is {podPhase(pod.phase)} on node <code>{pod.node}</code>
  </div>
)

export default Pod
