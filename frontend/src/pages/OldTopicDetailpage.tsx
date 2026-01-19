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
        const topicRes = await fetch(`http://localhost:8080/topics`);
        if (!topicRes.ok) throw new Error('failed to fetch topics');
        const topics: Topic[] = await topicRes.json();
        const foundTopic = topics.find(t => t.id === parseInt(id));
        if (!foundTopic) throw new Error('topic not found');

        const postsRes = await fetch(`http://localhost:8080/topics/${id}/posts`);
        if (!postsRes.ok) throw new Error('failed to fetch posts');
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

  if (loading) return <p>Loading topic...</p>;
  if (error) return <p style={{ color: 'red' }}>{error}</p>;
  if (!topic) return <p>Topic not found.</p>;

  return (
    <div>
      <h2>{topic.title}</h2>
      <p><em>{topic.description}</em></p>
      <small>Created by user {topic.created_by} • {new Date(topic.created_at).toLocaleString()}</small>

      <br /><br />

      <button onClick={handleCreatePost}>Create New Post</button>

      <br /><br />

      <h3>Posts ({posts.length})</h3>

      {posts.length === 0 ? (
        <p>No posts yet. Be the first to post!</p>
      ) : (
        <ul style={{ listStyle: 'none', padding: 0 }}>
          {posts.map(post => (
            <li key={post.id} style={{ marginBottom: '1rem' }}>
              <h4>{post.title}</h4>
              <p>{post.body}</p>
              <small>By user {post.created_by} • {new Date(post.created_at).toLocaleString()}</small>
            </li>
          ))}
        </ul>
      )}
    </div>
  );
};

export default TopicDetailPage;
