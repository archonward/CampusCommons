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
        throw new Error(`failure in creating post: ${response.status} ${errorText}`);
      }

      // On success, go back to the topic detail page
      navigate(`/topics/${id}`);
    } catch (err: any) {
      console.error('Create post error:', err);
      setError(err.message || 'Failed to create post. Please try again.');
    } finally {
      setLoading(false);
    }
  };

  const handleCancel = () => {
    if (id) {
      navigate(`/topics/${id}`);
    } else {
      navigate('/topics');
    }
  };

  return (
    <div>
      <h2>Create New Post</h2>

      {error && <p style={{ color: 'red' }}>{error}</p>}

      <form onSubmit={handleSubmit}>
        <div>
          <label htmlFor="title">Title:</label>
          <br />
          <input
            id="title"
            type="text"
            value={title}
            onChange={(e) => setTitle(e.target.value)}
            disabled={loading}
            style={{ width: '100%', padding: '0.3rem' }}
          />
        </div>
        <br />

        <div>
          <label htmlFor="body">Body:</label>
          <br />
          <textarea
            id="body"
            value={body}
            onChange={(e) => setBody(e.target.value)}
            disabled={loading}
            rows={6}
            style={{ width: '100%', padding: '0.3rem' }}
          />
        </div>
        <br />

        <button type="submit" disabled={loading}>
          {loading ? 'Creating...' : 'Create Post'}
        </button>
        <button
          type="button"
          onClick={handleCancel}
          disabled={loading}
          style={{ marginLeft: '0.5rem' }}
        >
          Cancel
        </button>
      </form>
    </div>
  );
};

export default NewPostPage;
