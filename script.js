const taskInput = document.getElementById('taskInput');
const addBtn = document.getElementById('addBtn');
const taskList = document.getElementById('taskList');

async function loadTasks() {
    const res = await fetch('/tasks');
    const tasks = await res.json();
    
    taskList.innerHTML = tasks && tasks.length 
        ? tasks.map(t => `
            <div class="task-item ${t.done ? 'completed' : ''}">
                <span class="task-text" onclick="toggleTask(${t.id}, ${t.done})">${t.title}</span>
                <button class="delete-btn" onclick="deleteTask(${t.id})">Удалить</button>
            </div>
        `).join('') 
        : '<p style="text-align:center; color:#555">Список пуст</p>';
}

async function addTask() {
    const title = taskInput.value.trim();
    if (!title) return;

    await fetch('/tasks', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ title })
    });

    taskInput.value = '';
    loadTasks();
}

async function toggleTask(id, currentStatus) {
    await fetch('/tasks', {
        method: 'PUT',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ id: id, done: !currentStatus })
    });
    loadTasks();
}

async function deleteTask(id) {
    await fetch(`/tasks?id=${id}`, { method: 'DELETE' });
    loadTasks();
}

addBtn.addEventListener('click', addTask);
taskInput.addEventListener('keypress', (e) => { if (e.key === 'Enter') addTask(); });

loadTasks();