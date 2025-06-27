// Variables globales
let currentUser = null;
let currentTab = 'threads';

// Initialisation de la page
document.addEventListener('DOMContentLoaded', function() {
    loadUserProfile();
    setupEventListeners();
});

// Configuration des écouteurs d'événements
function setupEventListeners() {
    // Formulaire de mise à jour du profil
    const updateForm = document.getElementById('updateProfileForm');
    if (updateForm) {
        updateForm.addEventListener('submit', handleProfileUpdate);
    }

    // Upload de photo de profil
    const profilePictureInput = document.getElementById('profile_picture');
    if (profilePictureInput) {
        profilePictureInput.addEventListener('change', handleProfilePictureUpload);
    }
}

// Charger le profil utilisateur
async function loadUserProfile() {
    try {
        const response = await fetch('/api/users/me', {
            method: 'GET',
            headers: {
                'Content-Type': 'application/json',
                'Authorization': 'Bearer ' + getToken()
            }
        });

        if (response.ok) {
            const data = await response.json();
            currentUser = data.data;
            displayUserProfile(currentUser);
            showProfileActions();
            loadUserThreads();
        } else if (response.status === 401) {
            // Token expiré, rediriger vers la connexion
            window.location.href = '/login';
        } else {
            showMessage('Erreur lors du chargement du profil', 'error');
        }
    } catch (error) {
        console.error('Erreur:', error);
        showMessage('Erreur de connexion', 'error');
    }
}

// Afficher les informations du profil
function displayUserProfile(user) {
    // Informations de base
    document.getElementById('usernameProfile').textContent = user.username || 'Utilisateur';
    document.getElementById('username').textContent = user.username || 'Utilisateur';
    
    // Biographie
    const biographyElement = document.getElementById('biography');
    if (biographyElement) {
        biographyElement.textContent = user.biography || 'Aucune biographie disponible.';
    }

    // Photo de profil
    const profilePicture = document.getElementById('profilePicture');
    if (profilePicture) {
        profilePicture.src = user.profile_picture || '/static/images/default-avatar.png';
    }

    // Statistiques
    document.getElementById('threadCount').textContent = user.thread_count || 0;
    document.getElementById('messageCount').textContent = user.message_count || 0;
    
    // Dernière connexion
    const lastConnection = document.getElementById('lastConnection');
    if (lastConnection && user.last_login) {
        const date = new Date(user.last_login);
        lastConnection.textContent = date.toLocaleDateString('fr-FR');
    } else if (lastConnection) {
        lastConnection.textContent = 'Jamais connecté';
    }
}

// Afficher les actions du profil pour le propriétaire
function showProfileActions() {
    const profileActions = document.getElementById('profileActions');
    if (profileActions) {
        profileActions.style.display = 'block';
    }
}

// Basculer entre les onglets
function showTab(tabName) {
    // Masquer tous les contenus d'onglets
    const tabContents = document.querySelectorAll('.tab-content');
    tabContents.forEach(content => {
        content.style.display = 'none';
    });

    // Retirer la classe active de tous les boutons
    const tabButtons = document.querySelectorAll('.tab-btn');
    tabButtons.forEach(btn => {
        btn.classList.remove('active');
    });

    // Afficher le contenu sélectionné
    const selectedTab = document.getElementById(tabName + 'Tab');
    if (selectedTab) {
        selectedTab.style.display = 'block';
    }

    // Ajouter la classe active au bouton sélectionné
    event.target.classList.add('active');

    currentTab = tabName;

    // Charger le contenu approprié
    if (tabName === 'threads') {
        loadUserThreads();
    } else if (tabName === 'messages') {
        loadUserMessages();
    }
}

// Charger les discussions de l'utilisateur
async function loadUserThreads() {
    try {
        const response = await fetch(`/api/users/${currentUser.id}/threads`, {
            method: 'GET',
            headers: {
                'Content-Type': 'application/json',
                'Authorization': 'Bearer ' + getToken()
            }
        });

        if (response.ok) {
            const data = await response.json();
            displayUserThreads(data.data || []);
        } else {
            console.error('Erreur lors du chargement des discussions');
        }
    } catch (error) {
        console.error('Erreur:', error);
    }
}

// Afficher les discussions de l'utilisateur
function displayUserThreads(threads) {
    const container = document.getElementById('userThreads');
    if (!container) return;

    if (threads.length === 0) {
        container.innerHTML = '<p class="no-content">Aucune discussion créée.</p>';
        return;
    }

    const threadsHTML = threads.map(thread => `
        <div class="thread-item">
            <h3><a href="/threads/show/${thread.id}">${thread.title}</a></h3>
            <p class="thread-meta">
                Créé le ${new Date(thread.created_at).toLocaleDateString('fr-FR')} 
                - ${thread.message_count || 0} message(s)
            </p>
            <p class="thread-excerpt">${truncateText(thread.content, 150)}</p>
        </div>
    `).join('');

    container.innerHTML = threadsHTML;
}

// Charger les messages de l'utilisateur
async function loadUserMessages() {
    try {
        const response = await fetch(`/api/users/${currentUser.id}/messages`, {
            method: 'GET',
            headers: {
                'Content-Type': 'application/json',
                'Authorization': 'Bearer ' + getToken()
            }
        });

        if (response.ok) {
            const data = await response.json();
            displayUserMessages(data.data || []);
        } else {
            console.error('Erreur lors du chargement des messages');
        }
    } catch (error) {
        console.error('Erreur:', error);
    }
}

// Afficher les messages de l'utilisateur
function displayUserMessages(messages) {
    const container = document.getElementById('userMessages');
    if (!container) return;

    if (messages.length === 0) {
        container.innerHTML = '<p class="no-content">Aucun message posté.</p>';
        return;
    }

    const messagesHTML = messages.map(message => `
        <div class="message-item">
            <div class="message-header">
                <a href="/threads/show/${message.thread_id}" class="thread-link">
                    ${message.thread_title}
                </a>
                <span class="message-date">
                    ${new Date(message.created_at).toLocaleDateString('fr-FR')}
                </span>
            </div>
            <p class="message-content">${truncateText(message.content, 200)}</p>
        </div>
    `).join('');

    container.innerHTML = messagesHTML;
}

// Modifier le profil
function editProfile() {
    const form = document.getElementById('updateProfileForm');
    if (form) {
        // Pré-remplir le formulaire avec les données actuelles
        const biographyField = form.querySelector('#biography');
        if (biographyField && currentUser) {
            biographyField.value = currentUser.biography || '';
        }
        
        form.style.display = form.style.display === 'none' ? 'block' : 'none';
    }
}

// Gérer la mise à jour du profil
async function handleProfileUpdate(event) {
    event.preventDefault();
    
    const formData = new FormData(event.target);
    const updateData = {
        biography: formData.get('biography')
    };

    try {
        const response = await fetch('/api/users/me', {
            method: 'PUT',
            headers: {
                'Content-Type': 'application/json',
                'Authorization': 'Bearer ' + getToken()
            },
            body: JSON.stringify(updateData)
        });

        if (response.ok) {
            const data = await response.json();
            showMessage('Profil mis à jour avec succès', 'success');
            loadUserProfile(); // Recharger le profil
            document.getElementById('updateProfileForm').style.display = 'none';
        } else {
            const errorData = await response.json();
            showMessage(errorData.message || 'Erreur lors de la mise à jour', 'error');
        }
    } catch (error) {
        console.error('Erreur:', error);
        showMessage('Erreur de connexion', 'error');
    }
}

// Gérer l'upload de photo de profil
async function handleProfilePictureUpload(event) {
    const file = event.target.files[0];
    if (!file) return;

    // Vérifier le type de fichier
    if (!file.type.startsWith('image/')) {
        showMessage('Veuillez sélectionner une image valide', 'error');
        return;
    }

    // Vérifier la taille (max 5MB)
    if (file.size > 5 * 1024 * 1024) {
        showMessage('L\'image doit faire moins de 5MB', 'error');
        return;
    }

    const formData = new FormData();
    formData.append('profile_picture', file);

    try {
        const response = await fetch('/api/users/me/avatar', {
            method: 'POST',
            headers: {
                'Authorization': 'Bearer ' + getToken()
            },
            body: formData
        });

        if (response.ok) {
            const data = await response.json();
            showMessage('Photo de profil mise à jour avec succès', 'success');
            loadUserProfile(); // Recharger le profil
        } else {
            const errorData = await response.json();
            showMessage(errorData.message || 'Erreur lors de l\'upload', 'error');
        }
    } catch (error) {
        console.error('Erreur:', error);
        showMessage('Erreur de connexion', 'error');
    }
}

// Fonction utilitaire pour tronquer le texte
function truncateText(text, maxLength) {
    if (!text) return '';
    if (text.length <= maxLength) return text;
    return text.substring(0, maxLength) + '...';
}

// Obtenir le token JWT
function getToken() {
    return localStorage.getItem('jwt_token') || '';
}

// Afficher un message
function showMessage(message, type) {
    const messageElement = document.getElementById('message');
    if (messageElement) {
        messageElement.textContent = message;
        messageElement.className = `message ${type}`;
        messageElement.style.display = 'block';
        
        setTimeout(() => {
            messageElement.style.display = 'none';
        }, 5000);
    }
}