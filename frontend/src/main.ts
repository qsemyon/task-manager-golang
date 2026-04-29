import './style.css';

interface Task {
    id: number;
    title: string;
    done: boolean;
}

const taskInput = document.getElementById('taskInput') as HTMLInputElement;
const addBtn = document.getElementById('addBtn') as HTMLButtonElement;
const taskList = document.getElementById('taskList') as HTMLDivElement;

const createTaskTemplate = (task: Task): string => `
    <div class="task-item ${task.done ? 'completed' : ''}">
        <span class="task-text" data-id="${task.id}" data-done="${task.done}">
            ${task.title}
        </span>
        <button class="delete-btn" data-id="${task.id}">Удалить</button>
    </div>
`;

async function loadTasks(): Promise<void> {
    try {
        const response = await fetch('/tasks');
        if (!response.ok) throw new Error('Ошибка при загрузке');
        
        const tasks: Task[] = await response.json();
        
        if (!tasks || tasks.length === 0) {
            taskList.innerHTML = '<p style="text-align:center; color:#555">Список пуст</p>';
            return;
        }

        taskList.innerHTML = tasks.map(createTaskTemplate).join('');
        attachEventListeners();
    } catch (error) {
        console.error('Ошибка:', error);
        taskList.innerHTML = '<p style="text-align:center; color:red">Ошибка связи с сервером</p>';
    }
}

async function addTask(): Promise<void> {
    const title = taskInput.value.trim();
    if (!title) return;

    try {
        await fetch('/tasks', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ title })
        });
        taskInput.value = '';
        await loadTasks();
    } catch (error) {
        alert('Не удалось добавить задачу');
    }
}

function attachEventListeners(): void {
    taskList.querySelectorAll('.task-text').forEach(el => {
        el.addEventListener('click', (e) => {
            const target = e.currentTarget as HTMLElement;
            const id = Number(target.dataset.id);
            const done = target.dataset.done === 'true';
            toggleTask(id, done);
        });
    });

    taskList.querySelectorAll('.delete-btn').forEach(el => {
        el.addEventListener('click', (e) => {
            const target = e.currentTarget as HTMLElement;
            const id = Number(target.dataset.id);
            deleteTask(id);
        });
    });
}

async function toggleTask(id: number, currentStatus: boolean): Promise<void> {
    try {
        await fetch('/tasks', {
            method: 'PUT',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ id, done: !currentStatus })
        });
        await loadTasks();
    } catch (error) {
        console.error('Ошибка обновления');
    }
}

async function deleteTask(id: number): Promise<void> {
    try {
        await fetch(`/tasks?id=${id}`, { method: 'DELETE' });
        await loadTasks();
    } catch (error) {
        alert('Ошибка удаления');
    }
}

addBtn.addEventListener('click', addTask);
taskInput.addEventListener('keypress', (e: KeyboardEvent) => {
    if (e.key === 'Enter') addTask();
});

document.addEventListener('DOMContentLoaded', loadTasks);