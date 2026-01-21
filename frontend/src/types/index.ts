export interface User {
  id: number;
  username: string;
}

// Topic object from /topics
export interface Topic {
  id: number;
  title: string;
  description: string;
  created_by: number;
  created_at: string; // ISO 8601
}

export interface Post {
  id: number;
  topic_id: number;
  title: string;
  body: string;
  created_by: number;
  created_at: string;
}

export interface Comment {
  id: number;
  post_id: number;
  body: string;
  created_by: number;
  created_at: string; // ISO 8601
}
