// Fonction pour effectuer un appel API
async function apiCall(endpoint, method = 'GET', data = null) {
    console.log('[DEBUG] apiCall - Début de la requête:', {
        endpoint,
        method,
        data,
        url: window.location.href
    });

    const token = localStorage.getItem('jwt_token');
    console.log('[DEBUG] apiCall - Token JWT:', token ? 'Présent' : 'Absent');

    const headers = {
        'Content-Type': 'application/json'
    };

    // Ajouter le token JWT s'il existe
    if (token) {
        headers['Authorization'] = `Bearer ${token}`;
        console.log('[DEBUG] apiCall - Headers avec token:', headers);
    }

    const options = {
        method,
        headers
    };

    if (data) {
        options.body = JSON.stringify(data);
        console.log('[DEBUG] apiCall - Body:', options.body);
    }

    try {
        console.log('[DEBUG] apiCall - Envoi de la requête avec options:', {
            ...options,
            headers: Object.fromEntries(Object.entries(options.headers).map(([k, v]) => [k, k === 'Authorization' ? 'Bearer ***' : v]))
        });
        
        const response = await fetch(endpoint, options);
        console.log('[DEBUG] apiCall - Réponse brute:', {
            status: response.status,
            statusText: response.statusText,
            headers: Object.fromEntries(response.headers.entries()),
            url: response.url
        });

        let result;
        const contentType = response.headers.get('content-type');
        console.log('[DEBUG] apiCall - Content-Type:', contentType);

        if (contentType && contentType.includes('application/json')) {
            result = await response.json();
            console.log('[DEBUG] apiCall - Données JSON reçues:', result);
        } else {
            const text = await response.text();
            console.log('[DEBUG] apiCall - Données texte reçues:', text);
            try {
                result = JSON.parse(text);
            } catch (e) {
                console.error('[DEBUG] apiCall - Erreur de parsing JSON:', e);
                result = { message: text };
            }
        }

        if (!response.ok) {
            console.log('[DEBUG] apiCall - Erreur de réponse:', {
                status: response.status,
                message: result.message,
                data: result
            });

            if (response.status === 401) {
                console.log('[DEBUG] apiCall - Token invalide, redirection vers login');
                localStorage.removeItem('jwt_token');
                window.location.href = '/auth/login';
            }
            throw new Error(result.message || 'Une erreur est survenue');
        }

        return result;
    } catch (error) {
        console.error('[DEBUG] apiCall - Erreur détaillée:', {
            message: error.message,
            stack: error.stack,
            name: error.name,
            endpoint,
            method
        });
        throw error;
    }
}

// Exporter la fonction apiCall
window.apiCall = apiCall;

// Fonction pour charger les discussions
async function loadThreads() {
    try {
        console.log('[DEBUG] Début du chargement des discussions');
        const response = await apiCall('/api/threads');
        console.log('[DEBUG] Discussions chargées avec succès:', JSON.stringify(response, null, 2));
        return response.data || [];
    } catch (error) {
        console.error('[DEBUG] Erreur lors du chargement des discussions:', {
            message: error.message,
            stack: error.stack
        });
        throw error;
    }
}

// Fonction pour charger les statistiques
async function loadStats() {
    try {
        console.log('[DEBUG] Début du chargement des statistiques');
        const response = await apiCall('/api/stats');
        console.log('[DEBUG] Statistiques chargées avec succès:', JSON.stringify(response, null, 2));
        return response.data || {};
    } catch (error) {
        console.error('[DEBUG] Erreur lors du chargement des statistiques:', {
            message: error.message,
            stack: error.stack
        });
        throw error;
    }
}

// Objet API avec toutes les fonctions nécessaires
const api = {
    // Authentification
    auth: {
        login: (username, password) => apiCall('/api/auth/login', 'POST', { username, password }),
        register: (username, email, password) => apiCall('/api/auth/register', 'POST', { username, email, password })
    },

    // Utilisateurs
    users: {
        getProfile: () => apiCall('/api/users/me'),
        updateProfile: (data) => apiCall('/api/users/me', 'PUT', data),
        getStats: () => apiCall('/api/users/stats'),
        getThreads: () => apiCall('/api/users/threads'),
        getMessages: () => apiCall('/api/users/messages')
    },

    // Fils de discussion
    threads: {
        getAll: async () => {
            const response = await apiCall('/api/threads');
            return {
                status: response.status,
                data: {
                    threads: Array.isArray(response.data.threads) ? response.data.threads : []
                }
            };
        },
        getById: (id) => apiCall(`/api/threads/${id}`),
        create: (title, description, tags) => {
            const data = {
                title,
                description,
                tags: Array.isArray(tags) ? tags : [tags]
            };
            return apiCall('/api/threads', 'POST', data);
        },
        update: (id, data) => apiCall(`/api/threads/${id}`, 'PUT', data),
        delete: (id) => apiCall(`/api/threads/${id}`, 'DELETE')
    },

    // Messages
    messages: {
        getByThread: (threadId) => apiCall(`/api/threads/${threadId}/messages`),
        create: (threadId, content) => apiCall(`/api/threads/${threadId}/messages`, 'POST', { content }),
        update: (messageId, content) => apiCall(`/api/messages/${messageId}`, 'PUT', { content }),
        delete: (messageId) => apiCall(`/api/messages/${messageId}`, 'DELETE'),
        like: (messageId) => apiCall(`/api/messages/${messageId}/like`, 'POST'),
        dislike: (messageId) => apiCall(`/api/messages/${messageId}/dislike`, 'POST')
    },

    // Likes
    likes: {
        toggle: (messageId, type) => apiCall(`/api/messages/${messageId}/vote`, 'POST', { type })
    }
};

// Exporter l'objet api
window.api = api; 