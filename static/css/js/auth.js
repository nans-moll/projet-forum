

/**
 * Basculer la visibilité d'un champ mot de passe
 * @param {string} fieldId - ID du champ à basculer
 */
function togglePassword(fieldId) {
    const passwordInput = document.getElementById(fieldId);
    const toggle = passwordInput.nextElementSibling;
    
    if (passwordInput.type === 'password') {
        passwordInput.type = 'text';
        toggle.textContent = '🙈';
    } else {
        passwordInput.type = 'password';
        toggle.textContent = '👁️';
    }
}

/**
 * Afficher un message d'erreur ou de succès
 * @param {string} message - Le message à afficher
 * @param {string} type - Type de message ('error' ou 'success')
 */
function showMessage(message, type = 'error') {
    const errorDiv = document.getElementById('errorMessage');
    const successDiv = document.getElementById('successMessage');
    
    // Cacher tous les messages
    errorDiv.style.display = 'none';
    successDiv.style.display = 'none';
    
    // Afficher le message approprié
    if (type === 'error') {
        errorDiv.textContent = message;
        errorDiv.style.display = 'block';
        
        // Auto-hide après 5 secondes pour les erreurs
        setTimeout(() => {
            errorDiv.style.display = 'none';
        }, 5000);
    } else {
        successDiv.textContent = message;
        successDiv.style.display = 'block';
    }
    
    // Scroll vers le haut pour voir le message
    window.scrollTo({ top: 0, behavior: 'smooth' });
}

/**
 * Afficher/cacher l'indicateur de chargement
 * @param {boolean} show - Afficher ou cacher le loading
 */
function showLoading(show = true) {
    const loading = document.getElementById('loading');
    const submitBtn = document.getElementById('submitBtn');
    
    if (show) {
        loading.style.display = 'block';
        submitBtn.style.display = 'none';
    } else {
        loading.style.display = 'none';
        submitBtn.style.display = 'block';
    }
}



/**
 * Valider les champs du formulaire de connexion
 * @returns {boolean} - True si le formulaire est valide
 */
function validateForm() {
    const identifier = document.getElementById('identifier').value.trim();
    const password = document.getElementById('password').value;

    // Validation des champs
    const isIdentifierValid = identifier.length > 0;
    const isPasswordValid = password.length > 0;

    const isFormValid = isIdentifierValid && isPasswordValid;
    
    // Activer/désactiver le bouton de soumission
    document.getElementById('submitBtn').disabled = !isFormValid;
    
    return isFormValid;
}

/**
 * Valider un champ identifiant (email ou username)
 * @param {string} identifier - L'identifiant à valider
 * @returns {boolean} - True si valide
 */
function validateIdentifier(identifier) {
    if (identifier.length === 0) return false;
    
    // Vérifier si c'est un email ou un nom d'utilisateur
    const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
    const usernameRegex = /^[a-zA-Z0-9_]{3,20}$/;
    
    return emailRegex.test(identifier) || usernameRegex.test(identifier);
}


/**
 * Initialiser les event listeners pour la validation en temps réel
 */
function initEventListeners() {
    // Validation en temps réel pour l'identifiant
    document.getElementById('identifier').addEventListener('input', function() {
        const identifier = this.value.trim();
        const isValid = validateIdentifier(identifier);
        
        // Changer la couleur de la bordure selon la validité
        if (identifier.length > 0) {
            this.style.borderColor = isValid ? '#10b981' : '#ef4444';
        } else {
            this.style.borderColor = '';
        }
        
        validateForm();
    });

    // Validation en temps réel pour le mot de passe
    document.getElementById('password').addEventListener('input', function() {
        const password = this.value;
        const isValid = password.length > 0;
        
        // Changer la couleur de la bordure selon la validité
        if (password.length > 0) {
            this.style.borderColor = isValid ? '#10b981' : '#ef4444';
        } else {
            this.style.borderColor = '';
        }
        
        validateForm();
    });
}


/**
 * Gérer la soumission du formulaire de connexion
 * @param {Event} e - L'événement de soumission
 */
async function handleFormSubmission(e) {
    e.preventDefault();
    
    // Vérifier la validité du formulaire
    if (!validateForm()) {
        showMessage('Veuillez remplir tous les champs correctement', 'error');
        return;
    }
    
    // Préparer les données
    const formData = new FormData(e.target);
    const loginData = {
        identifier: formData.get('identifier').trim(),
        password: formData.get('password'),
        rememberMe: formData.get('rememberMe') ? true : false
    };
    
    // Afficher le loading
    showLoading(true);
    
    try {
        // Envoi de la requête au serveur
        const response = await fetch('/api/auth/login', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
                'Accept': 'application/json'
            },
            body: JSON.stringify(loginData)
        });
        
        const data = await response.json();
        
        if (response.ok && data.status === 'success') {
            // Succès de la connexion
            await handleLoginSuccess(data, loginData.rememberMe);
        } else {
            // Erreur de connexion
            handleLoginError(response, data);
        }
        
    } catch (error) {
        console.error('Erreur de connexion:', error);
        showMessage('Erreur de connexion. Vérifiez votre connexion internet et réessayez.', 'error');
    } finally {
        showLoading(false);
    }
}

/**
 * Gérer le succès de la connexion
 * @param {Object} data - Données de réponse du serveur
 * @param {boolean} rememberMe - Si l'utilisateur veut être remembered
 */
async function handleLoginSuccess(data, rememberMe) {
    // Stocker le token JWT selon la préférence de l'utilisateur
    if (data.token) {
        if (rememberMe) {
            localStorage.setItem('jwt_token', data.token);
            localStorage.setItem('user_data', JSON.stringify(data.user || {}));
        } else {
            sessionStorage.setItem('jwt_token', data.token);
            sessionStorage.setItem('user_data', JSON.stringify(data.user || {}));
        }
    }
    
    showMessage('Connexion réussie ! Redirection...', 'success');
    
    // Redirection intelligente
    const urlParams = new URLSearchParams(window.location.search);
    const redirectTo = urlParams.get('redirect') || '/';
    
    // Délai pour permettre à l'utilisateur de voir le message
    setTimeout(() => {
        window.location.href = redirectTo;
    }, 1500);
}

/**
 * Gérer les erreurs de connexion
 * @param {Response} response - Réponse HTTP
 * @param {Object} data - Données de réponse
 */
function handleLoginError(response, data) {
    let errorMessage = 'Erreur lors de la connexion';
    
    // Messages d'erreur spécifiques selon le code de statut
    switch (response.status) {
        case 400:
            errorMessage = 'Données invalides. Vérifiez vos informations.';
            break;
        case 401:
            errorMessage = 'Identifiants incorrects. Vérifiez votre email/nom d\'utilisateur et mot de passe.';
            break;
        case 403:
            errorMessage = 'Compte suspendu. Contactez l\'administration.';
            break;
        case 429:
            errorMessage = 'Trop de tentatives de connexion. Veuillez patienter avant de réessayer.';
            break;
        case 500:
            errorMessage = 'Erreur serveur. Veuillez réessayer plus tard.';
            break;
        default:
            if (data.message) {
                errorMessage = data.message;
            }
    }
    
    showMessage(errorMessage, 'error');
}



/**
 * Initialiser l'animation d'entrée de la page
 */
function initPageAnimation() {
    const card = document.querySelector('.auth-card');
    if (card) {
        card.style.opacity = '0';
        card.style.transform = 'translateY(20px)';
        
        setTimeout(() => {
            card.style.transition = 'all 0.6s ease';
            card.style.opacity = '1';
            card.style.transform = 'translateY(0)';
        }, 100);
    }
}

/**
 * Vérifier si l'utilisateur est déjà connecté
 */
function checkExistingAuth() {
    const token = localStorage.getItem('jwt_token') || sessionStorage.getItem('jwt_token');
    
    if (token) {
        showMessage('Vous êtes déjà connecté. Redirection...', 'success');
        setTimeout(() => {
            window.location.href = '/';
        }, 1500);
        return true;
    }
    
    return false;
}

/**
 * Gérer les messages depuis l'URL (paramètres GET)
 */
function handleUrlMessages() {
    const urlParams = new URLSearchParams(window.location.search);
    const message = urlParams.get('message');
    
    switch (message) {
        case 'account_created':
            const accountCreatedDiv = document.getElementById('accountCreatedMessage');
            if (accountCreatedDiv) {
                accountCreatedDiv.style.display = 'block';
            }
            break;
        case 'logout':
            showMessage('Vous avez été déconnecté avec succès.', 'success');
            break;
        case 'session_expired':
            showMessage('Votre session a expiré. Veuillez vous reconnecter.', 'error');
            break;
        case 'access_denied':
            showMessage('Accès refusé. Veuillez vous connecter.', 'error');
            break;
    }
    
    // Nettoyer l'URL après avoir traité le message
    if (message) {
        window.history.replaceState({}, document.title, window.location.pathname);
    }
}

/**
 * Initialiser les raccourcis clavier
 */
function initKeyboardShortcuts() {
    // Gestion de la touche Entrée pour soumettre le formulaire
    document.addEventListener('keypress', function(e) {
        if (e.key === 'Enter' && document.activeElement.tagName !== 'BUTTON') {
            const submitBtn = document.getElementById('submitBtn');
            if (submitBtn && !submitBtn.disabled) {
                submitBtn.click();
            }
        }
    });
    
    // Échap pour nettoyer les messages
    document.addEventListener('keydown', function(e) {
        if (e.key === 'Escape') {
            const errorDiv = document.getElementById('errorMessage');
            const successDiv = document.getElementById('successMessage');
            if (errorDiv) errorDiv.style.display = 'none';
            if (successDiv) successDiv.style.display = 'none';
        }
    });
}




 
function initLoginPage() {
    // Vérifier si l'utilisateur est déjà connecté
    if (checkExistingAuth()) {
        return; // Arrêter l'initialisation si déjà connecté
    }
    
    // Initialiser les composants de la page
    initPageAnimation();
    initEventListeners();
    handleUrlMessages();
    initKeyboardShortcuts();
    
    // Ajouter l'event listener pour la soumission du formulaire
    const loginForm = document.getElementById('loginForm');
    if (loginForm) {
        loginForm.addEventListener('submit', handleFormSubmission);
    }
    
    // Focus automatique sur le premier champ
    const identifierField = document.getElementById('identifier');
    if (identifierField) {
        identifierField.focus();
    }
    
    // Validation initiale
    validateForm();
    
    console.log('Page de connexion initialisée avec succès');
}


// Initialiser la page quand le DOM est prêt
document.addEventListener('DOMContentLoaded', initLoginPage);

// Exposer les fonctions globalement si nécessaire (pour les onclick dans le HTML)
window.togglePassword = togglePassword;