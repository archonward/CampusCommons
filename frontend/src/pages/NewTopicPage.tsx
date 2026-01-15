import React, { useState } from "react";
import { useNavigate } from "react-router-dom";

type CurrentUser = { id: number; username?: string } | null;

export default function NewTopicPage() {
  const navigate = useNavigate();

  const [form, setForm] = useState({ title: "", description: "" });
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState("");

  const currentUser: CurrentUser = JSON.parse(		//if user logged in, currentUser will have id & username
    localStorage.getItem("currentUser") || "null"	// if not, it will just be null
  );

  const handleChange = (
    event: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement>
  ) => {
    setForm((prev) => ({ ...prev, [event.target.name]: event.target.value })); // triple dot copies the old form
  };		// update only 1 field, dont change the other field. 

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();		// this line prevents a full page reload
    setError("");

    if (!currentUser) return setError("must be logged in.");
    if (!form.title.trim()) return setError("Title is required.");

    setLoading(true);
    try {
      const res = await fetch("http://localhost:8080/topics", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({			// converts the object into JSON text
          title: form.title.trim(),
          description: form.description.trim(),
          created_by: currentUser.id,
        }),
      });

      if (!res.ok) throw new Error(await res.text());
	// throwing an error here, we go directly to catch block
	// no need to navigate back to /topics when there is error

      navigate("/topics");

    } catch (err) {
      setError(
        err instanceof Error ? err.message : "failed to create topic."
      );
	console.error("Error:", err);
    } finally {
      setLoading(false);
    }
  };

  return (
    <div>
      <h2>Create Topic</h2>

      {error && <p style={{ color: "red" }}>{error}</p>}

      <form onSubmit={handleSubmit}>
        <div>
          <input
            name="title"
            placeholder="Title"
            value={form.title}
            onChange={handleChange}
            disabled={loading}
          />
        </div>

        <div>
          <textarea
            name="description"
            placeholder="Description (optional)"
            value={form.description}
            onChange={handleChange}
            disabled={loading}
            rows={4}
          />
        </div>

        <button type="submit" disabled={loading}>
          {loading ? "Creating..." : "Create"}
        </button>

        <button
          type="button"
          onClick={() => navigate("/topics")}
          disabled={loading}
          style={{ marginLeft: 8 }}
        >
          Cancel
        </button>
      </form>
    </div>
  );
}

