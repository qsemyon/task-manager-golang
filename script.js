// Константы для элементов управления
const taskInput = document.getElementById('taskInput');
const addBtn    = document.getElementById('addBtn');
const taskList  = document.getElementById('taskList');

const createTaskTemplate = (task) => `
    <div class="task-item ${task.done ? 'completed' : ''}">
        <span class="task-text" onclick="toggleTask(${task.id}, ${task.done})">
            ${task.title}
        </span>
        <button class="delete-btn" onclick="deleteTask(${task.id})">Удалить</button>
    </div>
`;

async function loadTasks() {
    try {
        const response = await fetch('/tasks');
        if (!response.ok) throw new Error('Ошибка при загрузке данных');
        
        const tasks = await response.json();
        
        if (!tasks || tasks.length === 0) {
            taskList.innerHTML = '<p style="text-align:center; color:#555">Список пуст</p>';
            return;
        }

        taskList.innerHTML = tasks.map(createTaskTemplate).join('');
    } catch (error) {
        console.error('Ошибка:', error);
        taskList.innerHTML = '<p style="text-align:center; color:red">Не удалось связаться с сервером</p>';
    }
}

async function addTask() {
    const title = taskInput.value.trim();
    if (!title) return;

    try {
        const response = await fetch('/tasks', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ title })
        });

        if (response.ok) {
            taskInput.value = '';
            await loadTasks();
        }
    } catch (error) {
        alert('Не удалось добавить задачу');
    }
}

async function toggleTask(id, currentStatus) {
    try {
        await fetch('/tasks', {
            method: 'PUT',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ id, done: !currentStatus })
        });
        await loadTasks();
    } catch (error) {
        console.error('Не удалось обновить задачу');
    }
}

async function deleteTask(id) {
    try {
        await fetch(`/tasks?id=${id}`, { method: 'DELETE' });
        await loadTasks();
    } catch (error) {
        alert('Ошибка при удалении');
    }
}

addBtn.addEventListener('click', addTask);

taskInput.addEventListener('keypress', (event) => {
    if (event.key === 'Enter') {
        addTask();
    }
});

document.addEventListener('DOMContentLoaded', loadTasks);