const API_BASE_URL = "";

// State
let authToken = localStorage.getItem("todo_auth_token") || "";
let currentUser = JSON.parse(localStorage.getItem("todo_current_user") || "null");
let currentFilter = "all";
let cachedTodos = [];

// DOM Elemanları - Auth
const authSection = document.getElementById("auth-section");
const appContent = document.getElementById("app-content");
const userProfile = document.getElementById("user-profile");
const currentUsername = document.getElementById("current-username");
const currentRoleBadge = document.getElementById("current-role-badge");
const logoutButton = document.getElementById("logout-button");

const tabLogin = document.getElementById("tab-login");
const tabRegister = document.getElementById("tab-register");
const authForm = document.getElementById("auth-form");
const authUsernameInput = document.getElementById("auth-username");
const authPasswordInput = document.getElementById("auth-password");
const authSubmitBtn = document.getElementById("auth-submit-btn");
const authStatus = document.getElementById("auth-status");

// DOM Elemanları - Todos
const todoForm = document.getElementById("todo-form");
const titleInput = document.getElementById("title");
const visibilitySelect = document.getElementById("visibility");
const addButton = document.getElementById("add-button");
const formStatus = document.getElementById("form-status");
const todoList = document.getElementById("todo-list");
const listStatus = document.getElementById("list-status");
const filterButtons = document.querySelectorAll(".filter-btn");

let authMode = "login"; // "login" veya "register"

// Helper: Yetkili fetch
async function authFetch(url, options = {}) {
    const headers = options.headers || {};
    if (authToken) {
        headers["Authorization"] = `Bearer ${authToken}`;
    }
    options.headers = headers;

    const response = await fetch(url, options);
    if (response.status === 401) {
        handleLogout("Oturum süreniz doldu, lütfen tekrar giriş yapın.");
        throw new Error("Oturum süresi doldu.");
    }
    return response;
}

// Auth Yönetimi
function setAuthMode(mode) {
    authMode = mode;
    tabLogin.classList.toggle("active", mode === "login");
    tabRegister.classList.toggle("active", mode === "register");
    authSubmitBtn.textContent = mode === "login" ? "Giriş Yap" : "Kayıt Ol";
    authStatus.textContent = "";
}

tabLogin.addEventListener("click", () => setAuthMode("login"));
tabRegister.addEventListener("click", () => setAuthMode("register"));

authForm.addEventListener("submit", async (e) => {
    e.preventDefault();
    const username = authUsernameInput.value.trim();
    const password = authPasswordInput.value;

    if (!username || !password) {
        authStatus.textContent = "Lütfen kullanıcı adı ve şifre girin.";
        return;
    }

    authSubmitBtn.disabled = true;
    authStatus.textContent = "İşlem yapılıyor...";

    const endpoint = authMode === "login" ? "/api/auth/login" : "/api/auth/register";

    try {
        const response = await fetch(`${API_BASE_URL}${endpoint}`, {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify({ username, password }),
        });

        if (!response.ok) {
            const errText = await response.text();
            throw new Error(errText || "İşlem başarısız.");
        }

        const data = await response.json();
        authToken = data.token;
        currentUser = data.user;

        localStorage.setItem("todo_auth_token", authToken);
        localStorage.setItem("todo_current_user", JSON.stringify(currentUser));

        authForm.reset();
        authStatus.textContent = "";
        updateUIForLoggedIn();
        await loadTodos();
    } catch (err) {
        authStatus.textContent = err.message;
    } finally {
        authSubmitBtn.disabled = false;
    }
});

function updateUIForLoggedIn() {
    if (currentUser) {
        authSection.classList.add("hidden");
        appContent.classList.remove("hidden");
        userProfile.classList.remove("hidden");

        currentUsername.textContent = currentUser.username;
        const isAdmin = currentUser.role === "superadmin";
        currentRoleBadge.textContent = isAdmin ? "Süper Admin" : "Kullanıcı";
        currentRoleBadge.className = `badge ${isAdmin ? "badge-admin" : "badge-user"}`;
    } else {
        authSection.classList.remove("hidden");
        appContent.classList.add("hidden");
        userProfile.classList.add("hidden");
        todoList.replaceChildren();
    }
}

function handleLogout(message = "") {
    authToken = "";
    currentUser = null;
    localStorage.removeItem("todo_auth_token");
    localStorage.removeItem("todo_current_user");
    updateUIForLoggedIn();
    if (message) {
        authStatus.textContent = message;
    }
}

logoutButton.addEventListener("click", () => handleLogout());

// Filtreler
filterButtons.forEach((btn) => {
    btn.addEventListener("click", () => {
        filterButtons.forEach((b) => b.classList.remove("active"));
        btn.classList.add("active");
        currentFilter = btn.dataset.filter;
        renderFilteredTodos();
    });
});

function renderFilteredTodos() {
    let filtered = cachedTodos;
    if (currentFilter === "mine") {
        filtered = cachedTodos.filter((t) => t.user_id === currentUser.id);
    } else if (currentFilter === "public") {
        filtered = cachedTodos.filter((t) => t.visibility === "public");
    }
    renderTodos(filtered);
}

// Todo Listesi Çizimi
function renderTodos(todos) {
    todoList.replaceChildren();

    if (todos.length === 0) {
        listStatus.textContent = "Gösterilecek görev bulunmuyor.";
        return;
    }
    listStatus.textContent = "";

    const isSuperAdmin = currentUser && currentUser.role === "superadmin";

    for (const todo of todos) {
        const item = document.createElement("li");
        item.className = "todo-item";
        item.classList.toggle("completed", todo.done);

        const isOwner = currentUser && todo.user_id === currentUser.id;
        const canModify = isOwner || isSuperAdmin;

        const mainDiv = document.createElement("div");
        mainDiv.className = "todo-main";

        const checkbox = document.createElement("input");
        checkbox.type = "checkbox";
        checkbox.checked = todo.done;
        checkbox.disabled = !canModify;
        if (!canModify) {
            checkbox.title = "Sadece görevin sahibi veya süper admin güncelleyebilir.";
        }

        const contentDiv = document.createElement("div");
        const titleSpan = document.createElement("div");
        titleSpan.className = "todo-title";
        titleSpan.textContent = todo.title;

        const metaDiv = document.createElement("div");
        metaDiv.className = "todo-meta";

        // Görünürlük Rozeti
        const isPublic = todo.visibility === "public";
        const visBadge = document.createElement("span");
        visBadge.className = `badge ${isPublic ? "badge-public" : "badge-private"}`;
        visBadge.textContent = isPublic ? "Herkese Açık" : "Özel";
        metaDiv.appendChild(visBadge);

        // Sahip bilgisi (başkası eklediyse veya admin görüyorsa)
        if (todo.username && (isSuperAdmin || !isOwner)) {
            const creatorSpan = document.createElement("span");
            creatorSpan.className = "todo-creator";
            creatorSpan.textContent = `Yazan: ${todo.username}`;
            metaDiv.appendChild(creatorSpan);
        }

        contentDiv.append(titleSpan, metaDiv);
        mainDiv.append(checkbox, contentDiv);
        item.append(mainDiv);

        if (canModify) {
            const deleteButton = document.createElement("button");
            deleteButton.type = "button";
            deleteButton.className = "delete-button";
            deleteButton.textContent = "Sil";

            checkbox.addEventListener("change", async () => {
                const previousDone = todo.done;
                const newDone = checkbox.checked;
                checkbox.disabled = true;
                listStatus.textContent = "Görev güncelleniyor...";

                try {
                    const response = await authFetch(
                        `${API_BASE_URL}/api/todos/${encodeURIComponent(todo.id)}`,
                        {
                            method: "PUT",
                            headers: { "Content-Type": "application/json" },
                            body: JSON.stringify({ done: newDone }),
                        }
                    );

                    if (!response.ok) {
                        const message = await response.text();
                        throw new Error(message || `HTTP ${response.status}`);
                    }

                    const updatedTodo = await response.json();
                    todo.done = updatedTodo.done;
                    item.classList.toggle("completed", updatedTodo.done);
                    listStatus.textContent = "Görev güncellendi.";
                } catch (error) {
                    checkbox.checked = previousDone;
                    listStatus.textContent = `Görev güncellenemedi: ${error.message}`;
                } finally {
                    checkbox.disabled = false;
                }
            });

            deleteButton.addEventListener("click", async () => {
                deleteButton.disabled = true;
                listStatus.textContent = "Görev siliniyor...";

                try {
                    const response = await authFetch(
                        `${API_BASE_URL}/api/todos/${encodeURIComponent(todo.id)}`,
                        { method: "DELETE" }
                    );

                    if (!response.ok) {
                        const message = await response.text();
                        throw new Error(message || `HTTP ${response.status}`);
                    }

                    cachedTodos = cachedTodos.filter((t) => t.id !== todo.id);
                    renderFilteredTodos();
                    listStatus.textContent = "Görev silindi.";
                } catch (error) {
                    listStatus.textContent = `Görev silinemedi: ${error.message}`;
                    deleteButton.disabled = false;
                }
            });

            item.append(deleteButton);
        } else {
            const readonlyTag = document.createElement("span");
            readonlyTag.className = "readonly-tag";
            readonlyTag.textContent = "Salt Okunur";
            item.append(readonlyTag);
        }

        todoList.append(item);
    }
}

// Todoları Yükleme
async function loadTodos() {
    if (!authToken) return;

    listStatus.textContent = "Görevler yükleniyor...";
    try {
        const response = await authFetch(`${API_BASE_URL}/api/todos`);
        if (!response.ok) {
            throw new Error(`API isteği başarısız: HTTP ${response.status}`);
        }

        cachedTodos = await response.json();
        renderFilteredTodos();
    } catch (error) {
        listStatus.textContent = `Görevler yüklenemedi: ${error.message}`;
    }
}

// Todo Ekleme
todoForm.addEventListener("submit", async (event) => {
    event.preventDefault();

    const title = titleInput.value.trim();
    const visibility = visibilitySelect.value;

    if (title === "") {
        formStatus.textContent = "Görev başlığı boş olamaz.";
        titleInput.focus();
        return;
    }

    addButton.disabled = true;
    titleInput.disabled = true;
    visibilitySelect.disabled = true;
    formStatus.textContent = "Görev ekleniyor...";

    try {
        const response = await authFetch(`${API_BASE_URL}/api/todos`, {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify({ title, visibility }),
        });

        if (!response.ok) {
            const message = await response.text();
            throw new Error(message || `HTTP ${response.status}`);
        }

        todoForm.reset();
        visibilitySelect.value = "private";
        formStatus.textContent = "Görev eklendi.";
        await loadTodos();
    } catch (error) {
        formStatus.textContent = `Görev eklenemedi: ${error.message}`;
    } finally {
        addButton.disabled = false;
        titleInput.disabled = false;
        visibilitySelect.disabled = false;
        titleInput.focus();
    }
});

// Sayfa ilk yüklendiğinde oturum durumunu doğrula
async function initApp() {
    if (authToken) {
        try {
            const response = await fetch(`${API_BASE_URL}/api/auth/me`, {
                headers: { Authorization: `Bearer ${authToken}` },
            });
            if (response.ok) {
                currentUser = await response.json();
                localStorage.setItem("todo_current_user", JSON.stringify(currentUser));
                updateUIForLoggedIn();
                await loadTodos();
                return;
            }
        } catch (e) {
            console.warn("Mevcut oturum doğrulanamadı:", e);
        }
        handleLogout();
    } else {
        updateUIForLoggedIn();
    }
}

initApp();
