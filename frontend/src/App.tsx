import React from 'react';
import logo from './logo.svg';
import './App.css';
import { BrowserRouter as Router, Routes, Route } from 'react-router-dom';
import LoginPage from './pages/LoginPage';
import TopicListPage from './pages/TopicListPage';
import NewTopicPage from './pages/NewTopicPage';
import TopicDetailPage from './pages/TopicDetailPage';
import NewPostPage from './pages/NewPostPage';

function App() {
  return (
    <Router>
      <Routes>
        <Route path="/" element={<LoginPage />} />
        <Route path="/login" element={<LoginPage />} />
	<Route path="/topics" element={<TopicListPage />} />
	<Route path="/topics/new" element={<NewTopicPage />} /> 
	<Route path="/topics/:id" element={<TopicDetailPage />} />
	<Route path="/topics/:id/posts/new" element={<NewPostPage />} />
      </Routes>
    </Router>
  );
}


export default App;
