import React, { useState, useEffect } from 'react';
import { Topic } from '../types';
import { useNavigate } from 'react-router-dom';


const TopicListPage: React.FC = () => {
  const navigate = useNavigate();

  const [topics, setTopics] = useState<Topic[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    const fetchTopics = async () => {
      try {
        const response = await fetch('http://localhost:8080/topics');
        if (!response.ok) {
          throw new Error(`No response from backend server: ${response.status}`);
        }
        const data: Topic[] = await response.json();
        setTopics(data);
      } catch (err: any) {
        setError(err.message || 'fail to load topics.');
      } finally {
        setLoading(false);
      }
    };

    fetchTopics();
  }, []);

  return (
    <div>
      <h2>Topics</h2>

      {error && <p style={{ color: 'red' }}>{error}</p>}
      {loading && <p>Loading topics...</p>}

      {!loading && !error && (
        <div>
	  <button onClick={() => navigate('/topics/new')}>
  		Create New Topic
	  </button>
          <br />
          <br />
          {topics.length === 0 ? (
            <p>No topics yet. You can create the first topic by using the button. </p>
          ) : (
            <ul style={{ listStyle: 'none', padding: 0 }}>
              {topics.map((topic) => (
                <li key={topic.id} style={{ marginBottom: '1rem' }}>
                  <h3>{topic.title}</h3>
                  <p>{topic.description}</p>
                  <small>Created by user {topic.created_by} • {new Date(topic.created_at).toLocaleString()}</small>
                </li>
              ))}
            </ul>
          )}
        </div>
      )}
    </div>
  );
};

export default TopicListPage;
