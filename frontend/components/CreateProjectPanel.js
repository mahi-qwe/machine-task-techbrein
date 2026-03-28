"use client";

import { useState } from "react";
import { apiRequest } from "../lib/api";

const initialForm = {
  name: "",
  description: "",
};

export default function CreateProjectPanel({ onCreated }) {
  const [form, setForm] = useState(initialForm);
  const [error, setError] = useState("");
  const [loading, setLoading] = useState(false);

  async function handleSubmit(event) {
    event.preventDefault();
    setLoading(true);
    setError("");

    try {
      await apiRequest("/projects", {
        method: "POST",
        body: JSON.stringify(form),
      });
      setForm(initialForm);
      await onCreated();
    } catch (err) {
      setError(err.message || "Failed to create project");
    } finally {
      setLoading(false);
    }
  }

  return (
    <section className="panel">
      <div className="panel-header">
        <div>
          <h2>Create Project</h2>
          <p className="muted section-copy">Add a project here, then use it immediately in the task creation modal.</p>
        </div>
      </div>

      <form onSubmit={handleSubmit} className="form-grid compact-form">
        <label>
          Project Name
          <input value={form.name} onChange={(event) => setForm({ ...form, name: event.target.value })} />
        </label>
        <label className="wide-field">
          Description
          <textarea
            rows="3"
            value={form.description}
            onChange={(event) => setForm({ ...form, description: event.target.value })}
          />
        </label>
        {error ? <p className="error-text">{error}</p> : null}
        <button type="submit" disabled={loading}>
          {loading ? "Creating..." : "Create Project"}
        </button>
      </form>
    </section>
  );
}
