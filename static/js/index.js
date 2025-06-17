// Vérification de l'état de l'authentification au chargement
document.addEventListener('DOMContentLoaded', function() {
    console.log('[DEBUG] DOMContentLoaded - Initializing page');
    checkAuthStatus();
    loadThreads();
    loadStats();
});

// Fonction pour vérifier l'état de l'authentification
function checkAuthStatus() {
    console.log('[DEBUG] checkAuthStatus - Checking authentication status');
    const token = localStorage.getItem('jwt_token');
    if (token) {
        try {
            // Décoder le token JWT
            const base64Url = token.split('.')[1];
            const base64 = base64Url.replace(/-/g, '+').replace(/_/g, '/');
            const jsonPayload = decodeURIComponent(atob(base64).split('').map(function(c) {
                return '%' + ('00' + c.charCodeAt(0).toString(16)).slice(-2);
            }).join(''));

            const payload = JSON.parse(jsonPayload);
            console.log('[DEBUG] checkAuthStatus - User authenticated:', payload.username);
            document.getElementById('authButtons').style.display = 'none';
            document.getElementById('userMenu').style.display = 'flex';
            document.getElementById('createThreadBtn').style.display = 'block';
            document.getElementById('username').textContent = payload.username;
        } catch (error) {
            console.error('[DEBUG] checkAuthStatus - Token decode error:', error);
            localStorage.removeItem('jwt_token');
            document.getElementById('authButtons').style.display = 'flex';
            document.getElementById('userMenu').style.display = 'none';
            document.getElementById('createThreadBtn').style.display = 'none';
        }
    } else {
        console.log('[DEBUG] checkAuthStatus - No token found');
        document.getElementById('authButtons').style.display = 'flex';
        document.getElementById('userMenu').style.display = 'none';
        document.getElementById('createThreadBtn').style.display = 'none';
    }
}

// Fonction pour charger les statistiques
async function loadStats() {
    console.log('[DEBUG] loadStats - Loading statistics');
    try {
        const response = await apiCall('/api/stats');
        console.log('[DEBUG] loadStats - API Response:', response);
        const stats = response.data || {};
        document.getElementById('usersCount').textContent = stats.user_count || 0;
        document.getElementById('threadsCount').textContent = stats.thread_count || 0;
        document.getElementById('messagesCount').textContent = stats.message_count || 0;
    } catch (error) {
        console.error('[DEBUG] loadStats - Error:', error);
        document.getElementById('usersCount').textContent = '0';
        document.getElementById('threadsCount').textContent = '0';
        document.getElementById('messagesCount').textContent = '0';
    }
}

// Fonction pour charger les discussions
async function loadThreads() {
    console.log('[DEBUG] loadThreads - Starting to load threads');
    try {
        const response = await api.threads.getAll();
        console.log('[DEBUG] loadThreads - API Response:', response);
        
        const container = document.getElementById('threadsContainer');
        container.innerHTML = '';

        const threads = response.data.threads || [];
        console.log('[DEBUG] loadThreads - Threads data:', threads);
        
        if (threads.length === 0) {
            container.innerHTML = '<p class="no-threads">Aucune discussion disponible pour le moment.</p>';
            return;
        }

        threads.forEach(thread => {
            console.log('[DEBUG] loadThreads - Processing thread:', thread);
            const threadElement = createThreadElement(thread);
            container.appendChild(threadElement);
        });
    } catch (error) {
        console.error('[DEBUG] loadThreads - Error:', error);
        const container = document.getElementById('threadsContainer');
        container.innerHTML = '<p class="error-message">Erreur lors du chargement des discussions. Veuillez réessayer plus tard.</p>';
    }
}

// Fonction pour créer un élément de discussion
function createThreadElement(thread) {
    console.log('[DEBUG] createThreadElement - Creating element for thread:', thread);
    const threadElement = document.createElement('div');
    threadElement.className = 'thread-card';
    threadElement.innerHTML = `
        <div class="thread-details">
            <h1><a href="/threads/show/${thread.id}">${thread.title}</a></h1>
            <div class="thread-meta">
                <span class="author">Par ${thread.author ? thread.author.username : 'Anonyme'}</span>
                <span class="date">Le ${new Date(thread.created_at).toLocaleDateString()}</span>
            </div>
            <div class="thread-content">${thread.description}</div>
            <div class="thread-tags">
                ${thread.tags ? thread.tags.split(',').map(tag => `<span class="tag">${tag.trim()}</span>`).join('') : ''}
            </div>
        </div>
    `;
    console.log('[DEBUG] createThreadElement - Created element:', threadElement.innerHTML);
    return threadElement;
}

// Fonction pour afficher le modal de création de discussion
function showCreateThread() {
    console.log('[DEBUG] showCreateThread - Opening modal');
    document.getElementById('createThreadModal').style.display = 'block';
}

// Fonction pour fermer le modal
function closeModal(modalId) {
    console.log('[DEBUG] closeModal - Closing modal:', modalId);
    document.getElementById(modalId).style.display = 'none';
}

// Gérer la soumission du formulaire de création de discussion
document.getElementById('createThreadForm').addEventListener('submit', async function(e) {
    e.preventDefault();
    console.log('[DEBUG] createThreadForm - Form submitted');

    const title = document.getElementById('threadTitle').value;
    const description = document.getElementById('threadContent').value;
    const category = document.getElementById('threadCategory').value;

    try {
        console.log('[DEBUG] createThreadForm - Creating thread:', { title, description, category });
        const response = await api.threads.create(title, description, [category]);
        console.log('[DEBUG] createThreadForm - API Response:', response);
        
        if (response.status === 'success') {
            closeModal('createThreadModal');
            loadThreads();
            this.reset();
        } else {
            alert('Erreur lors de la création de la discussion');
        }
    } catch (error) {
        console.error('[DEBUG] createThreadForm - Error:', error);
        alert('Erreur lors de la création de la discussion');
    }
});

// Fermer le modal si on clique en dehors
window.onclick = function(event) {
    if (event.target.className === 'modal') {
        console.log('[DEBUG] window.onclick - Closing modal by clicking outside');
        event.target.style.display = 'none';
    }
}

// Fonction pour la recherche
function searchThreads() {
    const query = document.getElementById('searchInput').value;
    const category = document.getElementById('categoryFilter').value;
    const sort = document.getElementById('sortFilter').value;
    
    console.log('[DEBUG] searchThreads - Searching with params:', { query, category, sort });
    window.location.href = `/search?q=${encodeURIComponent(query)}&category=${category}&sort=${sort}`;
}

// Fonction pour la déconnexion
function logout() {
    console.log('[DEBUG] logout - User logging out');
    localStorage.removeItem('jwt_token');
    window.location.reload();
} 