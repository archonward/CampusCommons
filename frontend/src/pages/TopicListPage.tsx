import React, { useState, useEffect } from 'react';
import { Topic } from '../types';
import { useNavigate } from 'react-router-dom';

const TopicListPage: React.FC = () => {
  const navigate = useNavigate();

  const [topics, setTopics] = useState<Topic[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  const handleLogout = () => {
	localStorage.removeItem('currentUser');	// clear user from localstorage
	navigate('/login');
  };


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
        setError(err.message || 'Failed to load topics.');
      } finally {
        setLoading(false);
      }
    };

    fetchTopics();
  }, []);

  return (
    <div style={{
      display: 'flex',
      justifyContent: 'center',
      alignItems: 'flex-start',
      minHeight: '100vh',
      backgroundColor: '#f5f7fa',
      margin: 0,
      fontFamily: 'Arial, sans-serif',
      padding: '2rem 1rem'
    }}>
      <div style={{
        backgroundColor: 'white',
        padding: '2rem',
        borderRadius: '8px',
        boxShadow: '0 4px 12px rgba(0,0,0,0.1)',
        width: '100%',
        maxWidth: '600px',
	position: 'relative'
      }}>

	<button
          onClick={handleLogout}
          style={{
            position: 'absolute',
            top: '1rem',
            right: '1rem',
            background: 'none',
            border: 'none',
            color: '#666',
            fontSize: '0.85rem',
            cursor: 'pointer',
            padding: '0.25rem 0.5rem',
            borderRadius: '4px'
          }}
          onMouseEnter={(e) => e.currentTarget.style.color = '#d32f2f'}
          onMouseLeave={(e) => e.currentTarget.style.color = '#666'}
        >
          Logout
        </button>
        <h2 style={{ textAlign: 'center', margin: '0 0 1.5rem 0', color: '#333' }}>
          Topics
        </h2>

        {error && (
          <div style={{
            backgroundColor: '#ffebee',
            color: '#c62828',
            padding: '0.75rem',
            borderRadius: '4px',
            marginBottom: '1.5rem',
            textAlign: 'center'
          }}>
            {error}
          </div>
        )}

        {loading ? (
          <p style={{ textAlign: 'center', color: '#666' }}>Loading topics...</p>
        ) : (
          <div>
            <button
              onClick={() => navigate('/topics/new')}
              style={{
                backgroundColor: '#388e3c',
                color: 'white',
                border: 'none',
                borderRadius: '4px',
                padding: '0.6rem 1rem',
                fontSize: '1rem',
                fontWeight: 'bold',
                cursor: 'pointer',
                marginBottom: '1.5rem',
                width: '100%'
              }}
            >
              Create New Topic
            </button>

            {topics.length === 0 ? (
              <p style={{ textAlign: 'center', color: '#666', fontStyle: 'italic' }}>
                No topics yet. You can create the first topic using the button above.
              </p>
            ) : (
              <ul style={{ listStyle: 'none', padding: 0, margin: 0 }}>
                {topics.map((topic) => (
                  <li
                    key={topic.id}
                    style={{
                      padding: '1rem',
                      borderBottom: '1px solid #eee',
                      cursor: 'pointer',
                      transition: 'background-color 0.2s'
                    }}
                    onClick={() => navigate(`/topics/${topic.id}`)}
                  >
                    <h3 style={{ margin: '0 0 0.25rem 0', color: '#1976d2' }}>
                      {topic.title}
                    </h3>
                    <p style={{ margin: '0 0 0.5rem 0', color: '#555', fontSize: '0.95rem' }}>
                      {topic.description}
                    </p>
                    <small style={{ color: '#888', fontSize: '0.85rem' }}>
                      Created by user {topic.created_by} • {new Date(topic.created_at).toLocaleString()}
                    </small>
		     <button
     			 onClick={(e) => {
       				 e.stopPropagation(); // prevent triggering navigate
       				 if (window.confirm('delete this topic? This will also delete all posts and comments.')) {
         				 fetch(`http://localhost:8080/topics/${topic.id}`, { method: 'DELETE' })
           					 .then(res => {
             						 if (res.ok) {
               						    setTopics(topics.filter(t => t.id !== topic.id));
             						 } else {
               						    alert('Failed to delete topic');
             						 }
           		    });
       			   }
     		        }}
     			 style={{
       				 marginLeft: '0.5rem',
       				 background: 'none',
       				 border: '1px solid #d32f2f',
       				 color: '#d32f2f',
       				 borderRadius: '3px',
       				 padding: '0.1rem 0.4rem',
       				 fontSize: '0.8rem',
       				 cursor: 'pointer'
     			 }}
   		       >
     			 Delete
   		      </button>
                  </li>
                ))}
              </ul>
            )}
          </div>
        )}
      </div>
    </div>
  );
};

export default TopicListPage;
