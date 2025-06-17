// Fonction pour charger le profil de l'utilisateur
async function loadProfile() {
    try {
        console.log('[DEBUG] loadProfile - Début du chargement du profil');
        const token = localStorage.getItem('jwt_token');
        console.log('[DEBUG] loadProfile - Token JWT:', token ? 'Présent' : 'Absent');

        console.log('[DEBUG] loadProfile - Appel de l\'API pour récupérer le profil');
        const response = await apiCall('/api/users/me');
        console.log('[DEBUG] loadProfile - Réponse reçue:', response);

        if (response.status === 'success' && response.data) {
            console.log('[DEBUG] loadProfile - Mise à jour de l\'interface avec les données du profil');
            const user = response.data;
            document.getElementById('username').textContent = user.username;
            document.getElementById('email').textContent = user.email;
            document.getElementById('role').textContent = user.role;
            document.getElementById('messageCount').textContent = user.message_count;
            document.getElementById('threadCount').textContent = user.thread_count;
            document.getElementById('lastConnection').textContent = new Date(user.last_connection).toLocaleString();
            
            if (user.profile_picture) {
                document.getElementById('profilePicture').src = user.profile_picture;
            }
            
            if (user.biography) {
                document.getElementById('biography').textContent = user.biography;
            }
        }
    } catch (error) {
        console.error('[DEBUG] loadProfile - Erreur:', {
            message: error.message,
            stack: error.stack
        });
        showMessage('Erreur lors du chargement du profil: ' + error.message);
    }
}

// Fonction pour mettre à jour le profil
async function updateProfile(event) {
    event.preventDefault();
    try {
        console.log('[DEBUG] updateProfile - Début de la mise à jour du profil');
        const formData = new FormData(event.target);
        const data = {
            biography: formData.get('biography'),
            profile_picture: formData.get('profile_picture')
        };
        console.log('[DEBUG] updateProfile - Données à envoyer:', data);

        console.log('[DEBUG] updateProfile - Appel de l\'API pour mettre à jour le profil');
        const response = await apiCall('/api/users/me', 'PUT', data);
        console.log('[DEBUG] updateProfile - Réponse reçue:', response);

        if (response.status === 'success') {
            console.log('[DEBUG] updateProfile - Mise à jour réussie');
            showMessage('Profil mis à jour avec succès', 'success');
            loadProfile(); // Recharger le profil
        }
    } catch (error) {
        console.error('[DEBUG] updateProfile - Erreur:', {
            message: error.message,
            stack: error.stack
        });
        showMessage('Erreur lors de la mise à jour du profil: ' + error.message);
    }
}

// Fonction pour afficher un message
function showMessage(message, type = 'error') {
    console.log('[DEBUG] showMessage - Affichage du message:', {
        message,
        type
    });
    const messageDiv = document.getElementById('message');
    if (messageDiv) {
        messageDiv.textContent = message;
        messageDiv.className = `message ${type}`;
        messageDiv.style.display = 'block';
        setTimeout(() => {
            messageDiv.style.display = 'none';
        }, 5000);
    }
}

// Charger le profil au chargement de la page
document.addEventListener('DOMContentLoaded', function() {
    console.log('[DEBUG] DOMContentLoaded - Initialisation de la page de profil');
    loadProfile();

    // Ajouter l'écouteur d'événement pour le formulaire de mise à jour
    const updateForm = document.getElementById('updateProfileForm');
    if (updateForm) {
        console.log('[DEBUG] DOMContentLoaded - Ajout de l\'écouteur sur le formulaire de mise à jour');
        updateForm.addEventListener('submit', updateProfile);
    }
}); 