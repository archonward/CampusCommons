import React, { useState, useEffect } from 'react';
import { useNavigate, useParams } from 'react-router-dom';
import { Topic, Post } from '../types';

const TopicDetailPage: React.FC = () => {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();

  const [topic, setTopic] = useState<Topic | null>(null);
  const [posts, setPosts] = useState<Post[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    if (!id) {
      setError('Invalid topic ID');
      setLoading(false);
      return;
    }

    const fetchTopicAndPosts = async () => {
      try {
        // Fetch all topics to find the one with matching ID
        const topicRes = await fetch(`http://localhost:8080/topics`);
        if (!topicRes.ok) throw new Error('Failed to fetch topics');
        const topics: Topic[] = await topicRes.json();
        const foundTopic = topics.find(t => t.id === parseInt(id));
        if (!foundTopic) throw new Error('Topic not found');

        // Fetch posts for this topic
        const postsRes = await fetch(`http://localhost:8080/topics/${id}/posts`);
        if (!postsRes.ok) throw new Error('Failed to fetch posts');
        const postsData: Post[] = await postsRes.json();

        setTopic(foundTopic);
        setPosts(postsData);
      } catch (err: any) {
        setError(err.message || 'Failed to load topic and posts.');
      } finally {
        setLoading(false);
      }
    };

    fetchTopicAndPosts();
  }, [id]);

  const handleCreatePost = () => {
    navigate(`/topics/${id}/posts/new`);
  };

  if (loading) {
    return (
      <div style={{
        display: 'flex',
        justifyContent: 'center',
        alignItems: 'center',
        minHeight: '100vh',
        backgroundColor: '#f5f7fa',
        fontFamily: 'Arial, sans-serif'
      }}>
        <p style={{ color: '#666' }}>Loading topic...</p>
      </div>
    );
  }

  if (error) {
    return (
      <div style={{
        display: 'flex',
        justifyContent: 'center',
        alignItems: 'center',
        minHeight: '100vh',
        backgroundColor: '#f5f7fa',
        fontFamily: 'Arial, sans-serif'
      }}>
        <div style={{
          backgroundColor: '#ffebee',
          color: '#c62828',
          padding: '1rem',
          borderRadius: '8px',
          textAlign: 'center',
          maxWidth: '500px'
        }}>
          {error}
        </div>
      </div>
    );
  }

  if (!topic) {
    return (
      <div style={{
        display: 'flex',
        justifyContent: 'center',
        alignItems: 'center',
        minHeight: '100vh',
        backgroundColor: '#f5f7fa',
        fontFamily: 'Arial, sans-serif'
      }}>
        <p style={{ color: '#666' }}>Topic not found.</p>
      </div>
    );
  }

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
        maxWidth: '700px'
      }}>
        {/* Topic Header */}
        <div style={{ marginBottom: '1.5rem' }}>
          <h2 style={{ margin: '0 0 0.5rem 0', color: '#333' }}>{topic.title}</h2>
          <p style={{ margin: '0 0 0.75rem 0', color: '#555', fontStyle: 'italic' }}>
            {topic.description}
          </p>
          <small style={{ color: '#888', fontSize: '0.9rem' }}>
            Created by user {topic.created_by} • {new Date(topic.created_at).toLocaleString()}
          </small>
        </div>

        {/* Create Post Button */}
        <button
          onClick={handleCreatePost}
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
          Create New Post
        </button>

        {/* Posts Section */}
        <h3 style={{ color: '#333', marginBottom: '1rem' }}>
          Posts ({posts.length})
        </h3>

        {posts.length === 0 ? (
          <p style={{ textAlign: 'center', color: '#666', fontStyle: 'italic' }}>
            No posts yet. Be the first to post!
          </p>
        ) : (
          <ul style={{ listStyle: 'none', padding: 0, margin: 0 }}>
            {posts.map(post => (
              <li
                key={post.id}
                style={{
                  padding: '1rem',
                  borderBottom: '1px solid #eee',
                  marginBottom: '0.75rem',
                  backgroundColor: '#fafafa',
                  borderRadius: '4px'
                }}
              >
                <h4	onClick={() => navigate(`/posts/${post.id}`)}
  			style={{ margin: '0 0 0.25rem 0', cursor: 'pointer', color: '#1976d2' }}
		>
  			{post.title}
		</h4>
                <p style={{ margin: '0 0 0.5rem 0', color: '#444', whiteSpace: 'pre-wrap' }}>
                  {post.body}
                </p>
                <small style={{ color: '#888', fontSize: '0.85rem' }}>
                  By user {post.created_by} • {new Date(post.created_at).toLocaleString()}
                </small>
              </li>
            ))}
          </ul>
        )}
      </div>
    </div>
  );
};

export default TopicDetailPage;
