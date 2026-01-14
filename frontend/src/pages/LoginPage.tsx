import React, { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { login } from '../services/api';

const LoginPage: React.FC = () => {
  const [username, setUsername] = useState('');
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const navigate = useNavigate();

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (username == "" || username == " ") {
      setError('Username is required');
      return;
    }

	console.log("page did not refresh. Login button no bug.");

    	setLoading(true);
    	setError(null);

    	try {
      		await login(username.trim());
      		// For now, just go to topics
      		navigate('/topics');
    	} catch (err: any) {
      		setError(err.message || 'Login failed');
		console.error(err);
    	} finally {
      		setLoading(false);
    	}
  };

  return (
    <div>
      <h2>Login page for CampusCommons webApp</h2>
      <p>Type in your username.</p>

      {error && <p style={{ color: 'red' }}>{error}</p>}
	
     {/* React knows to call the handleSubmit function and the API call */} 
      <form onSubmit={handleSubmit}>	
        <div>
          <label htmlFor="username">Username:</label>
          <br />
          <input
            id="username"
            type="text"
            value={username}
            onChange={(e) => setUsername(e.target.value)}
            disabled={loading}
          />
        </div>
        <br />
        <button type="submit" disabled={loading}>
          {loading ? 'Logging in...' : 'Log In'}
        </button>
      </form>
    </div>
  );
};

export default LoginPage;
