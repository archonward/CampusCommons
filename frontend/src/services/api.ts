import type { User } from "../types";


const API_BASE_URL = 'http://localhost:8080';

// DO NOT TOUCH ANYTHING HERE!!!! 

// this part is for POST /login, quite buggy but i got it to work after fixing what the backend responds with.
export const login = async (username: string): Promise<User> => {
  const response = await fetch(`${API_BASE_URL}/login`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',	// need this line else get a 400 error 
    },
    body: JSON.stringify({ username }),	// need this line also, to send over the username as a String to backend
  });

// for easier debugging, DO NOT REMOVE the lines
// console.log("response object:", response);
// console.log("status:", response.status, "ok:", response.ok);
// console.log("url:", response.url);


  if (!response.ok) {
    const errorText = await response.text();
    throw new Error(`Login failed: ${response.status} ${errorText}`);
  }

  return response.json();
};
