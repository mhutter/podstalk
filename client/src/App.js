import React from 'react';
import './App.css';
import logo from './appuioli.png'

function App() {
  return (
    <div className="app">
      <header className="app__header">
        <img src={logo} className="app__logo" alt="APPUiOli" />
      </header>
      <main>Pod List goes here</main>
    </div>
  );
}

export default App;
