"use client";

import { useState } from "react";
import { apiRequest } from "../lib/api";

export default function CreateTaskModal({ open, onClose, projects, users, onCreated }) {
  const [form, setForm] = useState({
    title: "",
    description: "",
    project_id: "",
    assigned_to: "",
    due_date: "",
  });
  const [error, setError] = useState("");
  const [loading, setLoading] = useState(false);

  if (!open) return null;

  async function handleSubmit(event) {
    event.preventDefault();
    setLoading(true);
    setError("");

    try {
      await apiRequest("/tasks", {
        method: "POST",
        body: JSON.stringify({
          ...form,
          project_id: Number(form.project_id),
          assigned_to: form.assigned_to ? Number(form.assigned_to) : undefined,
        }),
      });
      setForm({ title: "", description: "", project_id: "", assigned_to: "", due_date: "" });
      await onCreated();
    } catch (err) {
      setError(err.message || "Failed to create task");
    } finally {
      setLoading(false);
    }
  }

  return (
    <div className="modal-backdrop">
      <section className="modal-card">
        <div className="panel-header">
          <h2>Create Task</h2>
          <button className="ghost" onClick={onClose}>
            Close
          </button>
        </div>
        <form onSubmit={handleSubmit} className="form-grid">
          <label>
            Title
            <input value={form.title} onChange={(event) => setForm({ ...form, title: event.target.value })} />
          </label>
          <label>
            Description
            <textarea
              rows="4"
              value={form.description}
              onChange={(event) => setForm({ ...form, description: event.target.value })}
            />
          </label>
          <label>
            Project
            <select value={form.project_id} onChange={(event) => setForm({ ...form, project_id: event.target.value })}>
              <option value="">Select project</option>
              {projects.map((project) => (
                <option key={project.id} value={project.id}>
                  {project.name}
                </option>
              ))}
            </select>
          </label>
          <label>
            Assign To
            <select value={form.assigned_to} onChange={(event) => setForm({ ...form, assigned_to: event.target.value })}>
              <option value="">Unassigned</option>
              {users.map((user) => (
                <option key={user.id} value={user.id}>
                  {user.name} ({user.role})
                </option>
              ))}
            </select>
          </label>
          <label>
            Due Date
            <input type="date" value={form.due_date} onChange={(event) => setForm({ ...form, due_date: event.target.value })} />
          </label>
          {error ? <p className="error-text">{error}</p> : null}
          <button type="submit" disabled={loading}>
            {loading ? "Creating..." : "Create Task"}
          </button>
        </form>
      </section>
    </div>
  );
}
