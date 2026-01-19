import React, { useState } from 'react';
import { useNavigate, useParams } from 'react-router-dom';

const NewPostPage: React.FC = () => {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();

  const [title, setTitle] = useState('');
  const [body, setBody] = useState('');
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const currentUser = JSON.parse(localStorage.getItem('currentUser') || 'null');

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();

    if (!id) {
      setError('Topic ID is missing.');
      return;
    }

    if (!currentUser) {
      setError('You must be logged in to create a post.');
      return;
    }

    if (!title.trim()) {
      setError('Title is required.');
      return;
    }

    if (!body.trim()) {
      setError('Post body is required.');
      return;
    }

    setLoading(true);
    setError(null);

    try {
      const response = await fetch(`http://localhost:8080/topics/${id}/posts`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          title: title.trim(),
          body: body.trim(),
          created_by: currentUser.id,
        }),
      });

      if (!response.ok) {
        const errorText = await response.text();
        throw new Error(`Failed to create post: ${response.status} ${errorText}`);
      }

      navigate(`/topics/${id}`);
    } catch (err: any) {
      console.error('Create post error:', err);
      setError(err.message || 'Failed to create post. Please try again.');
    } finally {
      setLoading(false);
    }
  };

  const handleCancel = () => {
    navigate(id ? `/topics/${id}` : '/topics');
  };

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
        maxWidth: '600px'
      }}>
        <h2 style={{ textAlign: 'center', margin: '0 0 1.5rem 0', color: '#333' }}>
          Create New Post
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

        <form onSubmit={handleSubmit}>
          <div style={{ marginBottom: '1.25rem' }}>
            <label htmlFor="title" style={{ display: 'block', marginBottom: '0.5rem', fontWeight: 'bold', color: '#444' }}>
              Title:
            </label>
            <input
              id="title"
              type="text"
              value={title}
              onChange={(e) => setTitle(e.target.value)}
              disabled={loading}
              style={{
                width: '100%',
                padding: '0.75rem',
                fontSize: '1rem',
                border: '1px solid #ccc',
                borderRadius: '4px',
                boxSizing: 'border-box'
              }}
            />
          </div>

          <div style={{ marginBottom: '1.5rem' }}>
            <label htmlFor="body" style={{ display: 'block', marginBottom: '0.5rem', fontWeight: 'bold', color: '#444' }}>
              Body:
            </label>
            <textarea
              id="body"
              value={body}
              onChange={(e) => setBody(e.target.value)}
              disabled={loading}
              rows={6}
              style={{
                width: '100%',
                padding: '0.75rem',
                fontSize: '1rem',
                border: '1px solid #ccc',
                borderRadius: '4px',
                resize: 'vertical',
                boxSizing: 'border-box',
                fontFamily: 'inherit'
              }}
            />
          </div>

          <div style={{ display: 'flex', gap: '0.75rem' }}>
            <button
              type="submit"
              disabled={loading}
              style={{
                flex: 1,
                padding: '0.75rem',
                backgroundColor: loading ? '#cccccc' : '#1976d2',
                color: 'white',
                border: 'none',
                borderRadius: '4px',
                fontSize: '1rem',
                fontWeight: 'bold',
                cursor: loading ? 'not-allowed' : 'pointer'
              }}
            >
              {loading ? 'Creating...' : 'Create Post'}
            </button>

            <button
              type="button"
              onClick={handleCancel}
              disabled={loading}
              style={{
                flex: 1,
                padding: '0.75rem',
                backgroundColor: '#f5f5f5',
                color: '#333',
                border: '1px solid #ddd',
                borderRadius: '4px',
                fontSize: '1rem',
                cursor: loading ? 'not-allowed' : 'pointer'
              }}
            >
              Cancel
            </button>
          </div>
        </form>
      </div>
    </div>
  );
};

export default NewPostPage;
