"use client";

import { useEffect, useMemo, useState } from "react";
import { useRouter } from "next/navigation";
import CreateProjectPanel from "../../components/CreateProjectPanel";
import CreateTaskModal from "../../components/CreateTaskModal";
import CreateUserPanel from "../../components/CreateUserPanel";
import { apiRequest, clearToken, getToken } from "../../lib/api";

export default function DashboardPage() {
  const router = useRouter();
  const [projects, setProjects] = useState([]);
  const [tasks, setTasks] = useState([]);
  const [users, setUsers] = useState([]);
  const [statusFilter, setStatusFilter] = useState("");
  const [projectFilter, setProjectFilter] = useState("");
  const [showCreateTask, setShowCreateTask] = useState(false);
  const [assigningTaskId, setAssigningTaskId] = useState(null);
  const [error, setError] = useState("");

  const usersById = useMemo(() => {
    return new Map(users.map((user) => [String(user.id), user]));
  }, [users]);

  const filteredTasks = useMemo(() => {
    return tasks.filter((task) => {
      if (statusFilter && task.status !== statusFilter) return false;
      if (projectFilter && String(task.project_id) !== projectFilter) return false;
      return true;
    });
  }, [tasks, statusFilter, projectFilter]);

  async function loadData() {
    try {
      const [projectResponse, taskResponse, userResponse] = await Promise.all([
        apiRequest("/projects"),
        apiRequest("/tasks"),
        apiRequest("/users"),
      ]);

      setProjects(projectResponse.items || []);
      setTasks(taskResponse.items || []);
      setUsers(userResponse.items || []);
      setError("");
    } catch (err) {
      if (err.status === 401) {
        clearToken();
        router.push("/");
        return;
      }
      setError(err.message || "Failed to load dashboard");
    }
  }

  useEffect(() => {
    if (!getToken()) {
      router.push("/");
      return;
    }
    loadData();
  }, []);

  async function handleAssignTask(taskId, assignedTo) {
    if (!assignedTo) return;

    try {
      setAssigningTaskId(taskId);
      const updatedTask = await apiRequest(`/tasks/${taskId}/assign`, {
        method: "PATCH",
        body: JSON.stringify({ assigned_to: Number(assignedTo) }),
      });
      setTasks((currentTasks) =>
        currentTasks.map((task) =>
          task.id === taskId
            ? {
                ...task,
                assigned_to: updatedTask.assigned_to,
              }
            : task
        )
      );
      setError("");
    } catch (err) {
      setError(err.message || "Failed to assign task");
    } finally {
      setAssigningTaskId(null);
    }
  }

  return (
    <main className="dashboard-shell">
      <header className="hero">
        <div>
          <p className="eyebrow">Internal Workspace</p>
          <h1>Projects and Tasks</h1>
          <p className="muted">A minimal client for login, project browsing, task filtering, and quick task creation.</p>
        </div>
        <div className="hero-actions">
          <button onClick={() => setShowCreateTask(true)}>Create Task</button>
          <button
            className="ghost"
            onClick={() => {
              clearToken();
              router.push("/");
            }}
          >
            Logout
          </button>
        </div>
      </header>

      {error ? <p className="error-banner">{error}</p> : null}

      <CreateUserPanel onCreated={loadData} />
      <CreateProjectPanel onCreated={loadData} />

      <section className="panel">
        <div className="panel-header">
          <h2>Projects</h2>
          <span>{projects.length} total</span>
        </div>
        <div className="project-grid">
          {projects.map((project) => (
            <article key={project.id} className="project-card">
              <h3>{project.name}</h3>
              <p>{project.description}</p>
              <small>Created by user #{project.created_by}</small>
            </article>
          ))}
        </div>
      </section>

      <section className="panel">
        <div className="panel-header">
          <h2>Tasks</h2>
          <div className="filters">
            <select value={statusFilter} onChange={(event) => setStatusFilter(event.target.value)}>
              <option value="">All statuses</option>
              <option value="todo">To Do</option>
              <option value="in_progress">In Progress</option>
              <option value="done">Done</option>
            </select>
            <select value={projectFilter} onChange={(event) => setProjectFilter(event.target.value)}>
              <option value="">All projects</option>
              {projects.map((project) => (
                <option key={project.id} value={project.id}>
                  {project.name}
                </option>
              ))}
            </select>
          </div>
        </div>
        <div className="task-list">
          {filteredTasks.map((task) => (
            <article key={task.id} className="task-card">
              <div className="task-title-row">
                <h3>{task.title}</h3>
                <span className={`status-pill status-${task.status}`}>{task.status.replace("_", " ")}</span>
              </div>
              <p>{task.description}</p>
              <small>Project #{task.project_id}</small>
              <small>
                Assigned to:{" "}
                {task.assigned_to
                  ? usersById.get(String(task.assigned_to))?.name || `user #${task.assigned_to}`
                  : "Unassigned"}
              </small>
              <small>Due: {task.due_date ? new Date(task.due_date).toLocaleDateString() : "Not set"}</small>
              <label className="assign-control">
                Assign user
                <select
                  value={task.assigned_to || ""}
                  onChange={(event) => handleAssignTask(task.id, event.target.value)}
                  disabled={assigningTaskId === task.id}
                >
                  <option value="">Select user</option>
                  {users.map((user) => (
                    <option key={user.id} value={user.id}>
                      {user.name} ({user.role})
                    </option>
                  ))}
                </select>
              </label>
            </article>
          ))}
        </div>
      </section>

      <CreateTaskModal
        open={showCreateTask}
        onClose={() => setShowCreateTask(false)}
        projects={projects}
        users={users}
        onCreated={async () => {
          await loadData();
          setShowCreateTask(false);
        }}
      />
    </main>
  );
}
