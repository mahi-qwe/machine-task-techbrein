"use client";

import { useState } from "react";
import { apiRequest } from "../lib/api";

const initialForm = {
  name: "",
  email: "",
  password: "",
  role: "developer",
};

export default function CreateUserPanel({ onCreated }) {
  const [form, setForm] = useState(initialForm);
  const [error, setError] = useState("");
  const [loading, setLoading] = useState(false);

  async function handleSubmit(event) {
    event.preventDefault();
    setLoading(true);
    setError("");

    try {
      await apiRequest("/users", {
        method: "POST",
        body: JSON.stringify(form),
      });
      setForm(initialForm);
      await onCreated();
    } catch (err) {
      setError(err.message || "Failed to create user");
    } finally {
      setLoading(false);
    }
  }

  return (
    <section className="panel">
      <div className="panel-header">
        <div>
          <h2>Users</h2>
          <p className="muted section-copy">Create developers or extra admins directly from the dashboard.</p>
        </div>
      </div>

      <form onSubmit={handleSubmit} className="form-grid compact-form">
        <label>
          Name
          <input value={form.name} onChange={(event) => setForm({ ...form, name: event.target.value })} />
        </label>
        <label>
          Email
          <input
            type="email"
            value={form.email}
            onChange={(event) => setForm({ ...form, email: event.target.value })}
          />
        </label>
        <label>
          Password
          <input
            type="password"
            value={form.password}
            onChange={(event) => setForm({ ...form, password: event.target.value })}
          />
        </label>
        <label>
          Role
          <select value={form.role} onChange={(event) => setForm({ ...form, role: event.target.value })}>
            <option value="developer">Developer</option>
            <option value="admin">Admin</option>
          </select>
        </label>
        {error ? <p className="error-text">{error}</p> : null}
        <button type="submit" disabled={loading}>
          {loading ? "Creating..." : "Create User"}
        </button>
      </form>

      <p className="muted section-copy">New users will immediately appear in the task assignment dropdowns after refresh.</p>
    </section>
  );
}
