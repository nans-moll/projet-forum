// Fonction de connexion
async function handleLogin(event) {
    if (event && event.preventDefault) {
        event.preventDefault();
    }
    console.log('[DEBUG] login - Début de la connexion');

    const form = event.target;
    const username = form.username.value;
    const password = form.password.value;

    console.log('[DEBUG] login - Données du formulaire:', {
        username,
        password: '***' // Ne pas logger le mot de passe en clair
    });

    try {
        console.log('[DEBUG] login - Appel de l\'API de connexion');
        const response = await apiCall('/api/auth/login', 'POST', {
            username,
            password
        });

        console.log('[DEBUG] login - Réponse reçue:', response);

        if (response.status === 'success') {
            console.log('[DEBUG] login - Connexion réussie, stockage du token');
            localStorage.setItem('jwt_token', response.data.token);
            
            console.log('[DEBUG] login - Redirection vers la page d\'accueil');
            window.location.href = '/';
        }
    } catch (error) {
        console.error('[DEBUG] login - Erreur:', {
            message: error.message,
            stack: error.stack
        });
        showMessage(error.message);
    }
}

// Fonction pour l'inscription
async function register(username, email, password) {
    try {
        showLoading(true);
        const response = await apiCall('/api/auth/register', 'POST', { username, email, password });
        
        if (response.status === 'success') {
            showMessage('Inscription réussie ! Vous pouvez maintenant vous connecter.', 'success');
            setTimeout(() => {
                window.location.href = '/login';
            }, 2000);
        } else {
            showMessage('Erreur lors de l\'inscription');
        }
    } catch (error) {
        showMessage(error.message || 'Erreur lors de l\'inscription');
    } finally {
        showLoading(false);
    }
}

// Fonction pour vérifier l'état de l'authentification
function checkAuth() {
    console.log('[DEBUG] checkAuth - Vérification de l\'authentification');
    const token = localStorage.getItem('jwt_token');
    console.log('[DEBUG] checkAuth - Token JWT:', token ? 'Présent' : 'Absent');

    const authButtons = document.getElementById('authButtons');
    const userMenu = document.getElementById('userMenu');
    const username = document.getElementById('username');

    if (token) {
        console.log('[DEBUG] checkAuth - Utilisateur connecté');
        try {
            // Décoder le token JWT pour obtenir le nom d'utilisateur
            const payload = JSON.parse(atob(token.split('.')[1]));
            console.log('[DEBUG] checkAuth - Payload du token:', payload);

            if (authButtons) authButtons.style.display = 'none';
            if (userMenu) userMenu.style.display = 'flex';
            if (username) username.textContent = payload.username;
        } catch (error) {
            console.error('[DEBUG] checkAuth - Erreur de décodage du token:', {
                message: error.message,
                stack: error.stack
            });
            // Token invalide, déconnexion
            logout();
        }
    } else {
        console.log('[DEBUG] checkAuth - Utilisateur non connecté');
        if (authButtons) authButtons.style.display = 'flex';
        if (userMenu) userMenu.style.display = 'none';
    }
}

// Fonction pour la déconnexion
function logout() {
    console.log('[DEBUG] logout - Début de la déconnexion');
    localStorage.removeItem('jwt_token');
    console.log('[DEBUG] logout - Token supprimé, redirection vers login');
    window.location.href = '/login';
}

// Fonction pour afficher un message
function showMessage(message, type = 'error') {
    const errorDiv = document.getElementById('errorMessage');
    const successDiv = document.getElementById('successMessage');
    
    if (errorDiv) errorDiv.style.display = 'none';
    if (successDiv) successDiv.style.display = 'none';
    
    if (type === 'error' && errorDiv) {
        errorDiv.textContent = message;
        errorDiv.style.display = 'block';
    } else if (successDiv) {
        successDiv.textContent = message;
        successDiv.style.display = 'block';
    }
}

// Fonction pour afficher/masquer le chargement
function showLoading(show = true) {
    const loading = document.getElementById('loading');
    const submitBtn = document.querySelector('.auth-submit');
    
    if (loading && submitBtn) {
        if (show) {
            loading.style.display = 'block';
            submitBtn.style.display = 'none';
        } else {
            loading.style.display = 'none';
            submitBtn.style.display = 'block';
        }
    }
}

// Vérifier l'état de l'authentification au chargement de la page
document.addEventListener('DOMContentLoaded', function() {
    console.log('[DEBUG] DOMContentLoaded - Initialisation de l\'authentification');
    checkAuth();
});

// Fonction pour vérifier si l'utilisateur est authentifié
function isAuthenticated() {
    const token = localStorage.getItem('jwt_token');
    return token !== null;
}

// Exporter les fonctions nécessaires
window.auth = {
    isAuthenticated,
    checkAuth,
    logout,
    login: handleLogin
}; 