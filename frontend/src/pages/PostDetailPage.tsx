import React, { useState, useEffect } from 'react';
import { useNavigate, useParams } from 'react-router-dom';
import { Post, Comment } from '../types';

const PostDetailPage: React.FC = () => {
  const { postId } = useParams<{ postId: string }>();
  const navigate = useNavigate();

  const [post, setPost] = useState<Post | null>(null);
  const [comments, setComments] = useState<Comment[]>([]);
  const [newComment, setNewComment] = useState('');
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [submitting, setSubmitting] = useState(false);

  const currentUser = JSON.parse(localStorage.getItem('currentUser') || 'null');

  useEffect(() => {
    if (!postId) {
      setError('Invalid post ID');
      setLoading(false);
      return;
    }

  const fetchPostAndComments = async () => {
    try {
      // Fetch single post
      const postRes = await fetch(`http://localhost:8080/posts/${postId}`);
      if (!postRes.ok) throw new Error('Post not found');
      const postData: Post = await postRes.json();

      // Fetch comments
      const commentsRes = await fetch(`http://localhost:8080/posts/${postId}/comments`);
      if (!commentsRes.ok) throw new Error('Failed to fetch comments');
      const commentsData: Comment[] = await commentsRes.json();

      setPost(postData);
      setComments(commentsData);
    } catch (err: any) {
        setError(err.message || 'Failed to load post and comments.');
    } finally {
        setLoading(false);
    }
    };

    fetchPostAndComments();
  }, [postId]);

  const handleAddComment = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!postId || !currentUser || !newComment.trim()) return;

    setSubmitting(true);
    try {
      const response = await fetch(`http://localhost:8080/posts/${postId}/comments`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          body: newComment.trim(),
          created_by: currentUser.id,
        }),
      });

      if (!response.ok) throw new Error('Failed to add comment');

      const comment: Comment = await response.json();
      setComments([...comments, comment]);
      setNewComment('');
    } catch (err: any) {
      alert(err.message || 'Failed to post comment.');
    } finally {
      setSubmitting(false);
    }
  };

  if (loading) return <p style={{ textAlign: 'center', padding: '2rem' }}>Loading post...</p>;
  if (error) return <p style={{ color: 'red', textAlign: 'center', padding: '2rem' }}>{error}</p>;
  if (!post) return null;

  return (
    <div style={{ maxWidth: '700px', margin: '2rem auto', padding: '0 1rem' }}>
      {/* Back + Delete row */}
      <div style={{ display: 'flex', alignItems: 'center', marginBottom: '1rem' }}>
        <button onClick={() => navigate(-1)}>
          ← Back
        </button>

        {post && (
          <>
            <button
              onClick={() => {
                if (window.confirm('Delete this post and all its comments?')) {
                  fetch(`http://localhost:8080/posts/${post.id}`, { method: 'DELETE' })
                    .then(res => {
                      if (res.ok) {
                        navigate(`/topics/${post.topic_id}`); // go back to topic
                      } else {
                        alert('Failed to delete post');
                      }
                    });
                }
              }}
              style={{
                marginLeft: '1rem',
                background: 'none',
                border: '1px solid #d32f2f',
                color: '#d32f2f',
                borderRadius: '3px',
                padding: '0.2rem 0.5rem',
                fontSize: '0.9rem',
                cursor: 'pointer'
              }}
            >
              Delete Post
            </button>

	    <button
		onClick={() => navigate(`/posts/${post.id}/edit`)}
		style={{ marginLeft: '0.5rem' }}
		>
 		   Edit
	    </button>
          </>
        )}
      </div>

      <h2>{post.title}</h2>
      <p>{post.body}</p>
      <small>
        By user {post.created_by} • {new Date(post.created_at).toLocaleString()}
      </small>

      <hr style={{ margin: '2rem 0' }} />

      <h3>Comments ({comments.length})</h3>

      {comments.length === 0 ? (
        <p>No comments yet.</p>
      ) : (
        <ul style={{ listStyle: 'none', padding: 0 }}>
          {comments.map(comment => (
            <li
              key={comment.id}
              style={{
                marginBottom: '1rem',
                paddingBottom: '1rem',
                borderBottom: '1px solid #eee'
              }}
            >
              <p>{comment.body}</p>
              <small>
                By user {comment.created_by} • {new Date(comment.created_at).toLocaleString()}
              </small>
            </li>
          ))}
        </ul>
      )}

      {currentUser ? (
        <form onSubmit={handleAddComment} style={{ marginTop: '2rem' }}>
          <h4>Add a Comment</h4>
          <textarea
            value={newComment}
            onChange={(e) => setNewComment(e.target.value)}
            placeholder="Write your comment..."
            rows={3}
            style={{ width: '100%', padding: '0.5rem', marginBottom: '0.5rem' }}
            disabled={submitting}
          />
          <button type="submit" disabled={!newComment.trim() || submitting}>
            {submitting ? 'Posting...' : 'Post Comment'}
          </button>
        </form>
      ) : (
        <p>You must be logged in to comment.</p>
      )}
    </div>
  );
};

export default PostDetailPage;

