import React from 'react';
import logo from './logo.svg';
import './App.css';
import { BrowserRouter as Router, Routes, Route } from 'react-router-dom';
import LoginPage from './pages/LoginPage';
import TopicListPage from './pages/TopicListPage'; // ← new import

function App() {
  return (
    <Router>
      <Routes>
        <Route path="/" element={<LoginPage />} />
        <Route path="/login" element={<LoginPage />} />
	<Route path="/topics" element={<TopicListPage />} />
      </Routes>
    </Router>
  );
}


export default App;
