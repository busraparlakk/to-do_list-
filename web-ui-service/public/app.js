const API_BASE_URL = "";

const todoList = document.getElementById("todo-list");
const listStatus = document.getElementById("list-status");

function renderTodos(todos) {
    todoList.replaceChildren();

    for (const todo of todos) {
        const item = document.createElement("li");
        item.className = "todo-item";
        item.classList.toggle("completed", todo.done);

        const label = document.createElement("label");
        label.className = "todo-label";

        const checkbox = document.createElement("input");
        checkbox.type = "checkbox";
        checkbox.checked = todo.done;

        const title = document.createElement("span");
        title.textContent = todo.title;

        const deleteButton = document.createElement("button");
        deleteButton.type = "button";
        deleteButton.className = "delete-button";
        deleteButton.textContent = "Sil";

        checkbox.addEventListener("change", async () => {
            const previousDone = todo.done;
            const newDone = checkbox.checked;

            checkbox.disabled = true;
            deleteButton.disabled = true;
            listStatus.textContent = "Görev güncelleniyor...";

            try {
                const response = await fetch(
                    `${API_BASE_URL}/api/todos/${encodeURIComponent(todo.id)}`,
                    {
                        method: "PUT",
                        headers: {
                            "Content-Type": "application/json",
                        },
                        body: JSON.stringify({ done: newDone }),
                    }
                );

                if (!response.ok) {
                    const message = await response.text();
                    throw new Error(message || `HTTP ${response.status}`);
                }

                const updatedTodo = await response.json();

                todo.done = updatedTodo.done;
                checkbox.checked = updatedTodo.done;
                item.classList.toggle("completed", updatedTodo.done);

                listStatus.textContent = "Görev güncellendi.";
            } catch (error) {
                checkbox.checked = previousDone;
                console.error("Görev güncellenemedi:", error);
                listStatus.textContent =
                    `Görev güncellenemedi: ${error.message}`;
            } finally {
                checkbox.disabled = false;
                deleteButton.disabled = false;
            }
        });

        deleteButton.addEventListener("click", async () => {
            checkbox.disabled = true;
            deleteButton.disabled = true;
            listStatus.textContent = "Görev siliniyor...";

            try {
                const response = await fetch(
                    `${API_BASE_URL}/api/todos/${encodeURIComponent(todo.id)}`,
                    {
                        method: "DELETE",
                    }
                );

                if (!response.ok) {
                    const message = await response.text();
                    throw new Error(message || `HTTP ${response.status}`);
                }

                item.remove();

                listStatus.textContent =
                    todoList.children.length === 0
                        ? "Henüz görev bulunmuyor."
                        : "Görev silindi.";
            } catch (error) {
                console.error("Görev silinemedi:", error);
                listStatus.textContent =
                    `Görev silinemedi: ${error.message}`;
            } finally {
                checkbox.disabled = false;
                deleteButton.disabled = false;
            }
        });

        label.append(checkbox, title);
        item.append(label, deleteButton);
        todoList.append(item);
    }
}

async function loadTodos() {
    listStatus.textContent = "Görevler yükleniyor...";
    todoList.replaceChildren();

    try {
        const response = await fetch(`${API_BASE_URL}/api/todos`);

        if (!response.ok) {
            throw new Error(`API isteği başarısız: HTTP ${response.status}`);
        }

        const todos = await response.json();

        if (!Array.isArray(todos)) {
            throw new Error("API yanıtı bir liste olmalı.");
        }

        renderTodos(todos);

        listStatus.textContent =
            todos.length === 0 ? "Henüz görev bulunmuyor." : "";
    } catch (error) {
        console.error("Görevler yüklenemedi:", error);
        listStatus.textContent =
            "Görevler yüklenemedi. API bağlantısını kontrol et.";
    }
}

loadTodos();
const todoForm = document.getElementById("todo-form");
const titleInput = document.getElementById("title");
const addButton = document.getElementById("add-button");
const formStatus = document.getElementById("form-status");

todoForm.addEventListener("submit", async (event) => {
    event.preventDefault();

    const title = titleInput.value.trim();

    if (title === "") {
        formStatus.textContent = "Görev başlığı boş olamaz.";
        titleInput.focus();
        return;
    }

    addButton.disabled = true;
    titleInput.disabled = true;
    formStatus.textContent = "Görev ekleniyor...";

    try {
        const response = await fetch(`${API_BASE_URL}/api/todos`, {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
            },
            body: JSON.stringify({ title: title }),
        });

        if (!response.ok) {
            const message = await response.text();
            throw new Error(message || `HTTP ${response.status}`);
        }

        todoForm.reset();
        formStatus.textContent = "Görev eklendi.";

        await loadTodos();
    } catch (error) {
        console.error("Görev eklenemedi:", error);
        formStatus.textContent = `Görev eklenemedi: ${error.message}`;
    } finally {
        addButton.disabled = false;
        titleInput.disabled = false;
        titleInput.focus();
    }
});