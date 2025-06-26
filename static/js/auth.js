// Fonction de connexion
async function handleLogin(event) {
    if (event && event.preventDefault) {
        event.preventDefault();
    }
    console.log('[DEBUG] login - Début de la connexion');
    console.log('[DEBUG] login - URL actuelle:', window.location.href);

    const form = event.target;
    const username = form.username.value;
    const password = form.password.value;

    console.log('[DEBUG] login - Données du formulaire:', {
        username,
        password: '***' // Ne pas logger le mot de passe en clair
    });

    try {
        showLoading(true);
        console.log('[DEBUG] login - Appel de l\'API de connexion');
        const response = await apiCall('/api/auth/login', 'POST', {
            username,
            password
        });

        console.log('[DEBUG] login - Réponse complète:', {
            status: response.status,
            data: response.data,
            message: response.message
        });

        if (response.status === 'success') {
            console.log('[DEBUG] login - Connexion réussie, stockage du token');
            localStorage.setItem('jwt_token', response.data.token);
            
            // Vérifier que le token a bien été stocké
            const storedToken = localStorage.getItem('jwt_token');
            console.log('[DEBUG] login - Token stocké:', storedToken ? 'Présent' : 'Absent');
            
            console.log('[DEBUG] login - Redirection vers la page d\'accueil');
            showMessage('Connexion réussie ! Redirection...', 'success');
            
            setTimeout(() => {
                console.log('[DEBUG] login - Début de la redirection');
                window.location.href = '/';
            }, 1500);
        } else {
            console.error('[DEBUG] login - Erreur de connexion:', response.message);
            showMessage(response.message || 'Erreur lors de la connexion');
        }
    } catch (error) {
        console.error('[DEBUG] login - Erreur détaillée:', {
            message: error.message,
            stack: error.stack,
            name: error.name
        });
        showMessage(error.message || 'Erreur lors de la connexion');
    } finally {
        showLoading(false);
    }
    return false; // Empêcher la soumission par défaut
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
    console.log('[DEBUG] checkAuth - Début de la vérification');
    console.log('[DEBUG] checkAuth - URL actuelle:', window.location.href);
    
    const token = localStorage.getItem('jwt_token');
    console.log('[DEBUG] checkAuth - Token JWT:', token ? 'Présent' : 'Absent');

    const authButtons = document.getElementById('authButtons');
    const userMenu = document.getElementById('userMenu');
    const username = document.getElementById('username');

    console.log('[DEBUG] checkAuth - Éléments DOM:', {
        authButtons: authButtons ? 'Présent' : 'Absent',
        userMenu: userMenu ? 'Présent' : 'Absent',
        username: username ? 'Présent' : 'Absent'
    });

    if (token) {
        console.log('[DEBUG] checkAuth - Utilisateur connecté');
        try {
            // Décoder le token JWT pour obtenir le nom d'utilisateur
            const payload = JSON.parse(atob(token.split('.')[1]));
            console.log('[DEBUG] checkAuth - Payload du token:', payload);

            if (authButtons) authButtons.style.display = 'none';
            if (userMenu) userMenu.style.display = 'flex';
            if (username) username.textContent = payload.username;

            // Si on est sur la page de login et qu'on est déjà connecté, rediriger vers l'accueil
            if (window.location.pathname === '/auth/login') {
                console.log('[DEBUG] checkAuth - Redirection depuis la page login');
                showMessage('Vous êtes déjà connecté. Redirection...', 'success');
                setTimeout(() => {
                    window.location.href = '/';
                }, 1500);
            }
        } catch (error) {
            console.error('[DEBUG] checkAuth - Erreur de décodage du token:', {
                message: error.message,
                stack: error.stack,
                token: token.substring(0, 20) + '...' // Afficher le début du token pour debug
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
    console.log('[DEBUG] logout - URL actuelle:', window.location.href);
    
    localStorage.removeItem('jwt_token');
    console.log('[DEBUG] logout - Token supprimé');
    
    console.log('[DEBUG] logout - Redirection vers login');
    window.location.href = '/auth/login';
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

// Fonction pour vérifier l'authentification au chargement de la page
async function checkAuthOnLoad() {
    console.log('[DEBUG] checkAuthOnLoad - Vérification de l\'authentification');
    const token = localStorage.getItem('jwt_token');
    
    if (!token) {
        console.log('[DEBUG] checkAuthOnLoad - Pas de token, redirection vers login');
        window.location.href = '/login';
        return false;
    }
    
    try {
        // Vérifier si le token est valide en appelant une API protégée
        const response = await apiCall('/api/users/me');
        if (response.status !== 'success') {
            console.log('[DEBUG] checkAuthOnLoad - Token invalide, suppression et redirection');
            localStorage.removeItem('jwt_token');
            window.location.href = '/login';
            return false;
        }
        
        console.log('[DEBUG] checkAuthOnLoad - Authentification valide');
        return true;
    } catch (error) {
        console.error('[DEBUG] checkAuthOnLoad - Erreur lors de la vérification:', error);
        localStorage.removeItem('jwt_token');
        window.location.href = '/login';
        return false;
    }
}

// Fonction pour vérifier si l'utilisateur est authentifié
function isAuthenticated() {
    const token = localStorage.getItem('jwt_token');
    console.log('[DEBUG] isAuthenticated - Token:', token ? 'Présent' : 'Absent');
    return token !== null;
}

// Exporter les fonctions nécessaires
window.auth = {
    isAuthenticated,
    checkAuth,
    logout,
    login: handleLogin
};