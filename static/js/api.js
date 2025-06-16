// Configuration de l'API
const API_BASE_URL = 'http://localhost:8080/api';

// Fonction utilitaire pour les appels API
async function apiCall(endpoint, method = 'GET', data = null) {
    const options = {
        method,
        headers: {
            'Content-Type': 'application/json',
        },
    };

    if (data) {
        options.body = JSON.stringify(data);
    }

    try {
        const response = await fetch(`${API_BASE_URL}${endpoint}`, options);
        const result = await response.json();
        
        if (!response.ok) {
            throw new Error(result.message || 'Une erreur est survenue');
        }
        
        return result;
    } catch (error) {
        console.error('Erreur API:', error);
        throw error;
    }
}

// Fonctions d'authentification
const auth = {
    async login(email, password) {
        return apiCall('/auth/login', 'POST', { email, password });
    },

    async register(username, email, password) {
        return apiCall('/auth/register', 'POST', { username, email, password });
    }
};

// Fonctions pour les threads
const threads = {
    async getAll() {
        return apiCall('/threads');
    },

    async create(title, content) {
        return apiCall('/threads', 'POST', { title, content });
    }
};

// Fonctions pour les messages
const messages = {
    async getByThread(threadId) {
        return apiCall(`/messages?threadId=${threadId}`);
    },

    async create(threadId, content) {
        return apiCall('/messages', 'POST', { threadId, content });
    }
};

// Fonctions pour les likes
const likes = {
    async toggle(threadId) {
        return apiCall('/likes', 'POST', { threadId });
    }
};

// Export des fonctions
window.api = {
    auth,
    threads,
    messages,
    likes
}; 