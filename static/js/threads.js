// Fonction pour charger les détails d'une discussion
async function loadThread() {
    console.log('[DEBUG] loadThread - Starting to load thread details');
    try {
        const threadId = window.location.pathname.split('/').pop();
        console.log('[DEBUG] loadThread - Thread ID:', threadId);
        
        const response = await api.threads.getById(threadId);
        console.log('[DEBUG] loadThread - API Response:', response);
        
        if (response.status === 'success') {
            const thread = response.data;
            const container = document.getElementById('threadContainer');
            
            container.innerHTML = `
                <div class="thread-details">
                    <h1>${thread.title}</h1>
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

            // Afficher le formulaire de réponse si l'utilisateur est connecté
            if (window.auth.isAuthenticated()) {
                document.getElementById('messageForm').style.display = 'block';
            }

            // Charger les messages
            loadMessages(threadId);
        }
    } catch (error) {
        console.error('[DEBUG] loadThread - Error:', error);
        document.getElementById('threadContainer').innerHTML = '<p class="error-message">Erreur lors du chargement de la discussion.</p>';
    }
}

// Fonction pour charger les messages d'une discussion
async function loadMessages(threadId) {
    console.log('[DEBUG] loadMessages - Loading messages for thread:', threadId);
    try {
        const response = await api.messages.getByThread(threadId);
        console.log('[DEBUG] loadMessages - API Response:', response);
        
        const container = document.getElementById('messagesContainer');
        
        if (response.status === 'success') {
            const messages = response.data.messages || [];
            
            if (messages.length === 0) {
                container.innerHTML = '<p class="no-messages">Aucun message pour le moment.</p>';
                return;
            }

            container.innerHTML = messages.map(message => `
                <div class="message" data-id="${message.id}">
                    <div class="message-header">
                        <span class="author">${message.author ? message.author.username : 'Anonyme'}</span>
                        <span class="date">Le ${new Date(message.created_at).toLocaleDateString()}</span>
                    </div>
                    <div class="message-content">${message.content}</div>
                    <div class="message-actions">
                        <button onclick="likeMessage(${message.id})" class="like-btn">
                            👍 ${message.likes}
                        </button>
                        <button onclick="dislikeMessage(${message.id})" class="dislike-btn">
                            👎 ${message.dislikes}
                        </button>
                    </div>
                </div>
            `).join('');
        }
    } catch (error) {
        console.error('[DEBUG] loadMessages - Error:', error);
        document.getElementById('messagesContainer').innerHTML = '<p class="error-message">Erreur lors du chargement des messages.</p>';
    }
}

// Fonction pour liker un message
async function likeMessage(messageId) {
    console.log('[DEBUG] likeMessage - Liking message:', messageId);
    if (!window.auth.isAuthenticated()) {
        alert('Vous devez être connecté pour liker un message');
        return;
    }

    try {
        const response = await api.messages.like(messageId);
        console.log('[DEBUG] likeMessage - API Response:', response);
        if (response.status === 'success') {
            loadMessages(window.location.pathname.split('/').pop());
        }
    } catch (error) {
        console.error('[DEBUG] likeMessage - Error:', error);
    }
}

// Fonction pour disliker un message
async function dislikeMessage(messageId) {
    console.log('[DEBUG] dislikeMessage - Disliking message:', messageId);
    if (!window.auth.isAuthenticated()) {
        alert('Vous devez être connecté pour disliker un message');
        return;
    }

    try {
        const response = await api.messages.dislike(messageId);
        console.log('[DEBUG] dislikeMessage - API Response:', response);
        if (response.status === 'success') {
            loadMessages(window.location.pathname.split('/').pop());
        }
    } catch (error) {
        console.error('[DEBUG] dislikeMessage - Error:', error);
    }
}

// Gérer la soumission d'un nouveau message
document.addEventListener('DOMContentLoaded', () => {
    const messageForm = document.getElementById('newMessageForm');
    if (messageForm) {
        messageForm.addEventListener('submit', async function(e) {
            e.preventDefault();
            console.log('[DEBUG] newMessageForm - Form submitted');

            const threadId = window.location.pathname.split('/').pop();
            const content = document.getElementById('messageContent').value;

            try {
                console.log('[DEBUG] newMessageForm - Creating message:', { threadId, content });
                const response = await api.messages.create(threadId, content);
                console.log('[DEBUG] newMessageForm - API Response:', response);
                
                if (response.status === 'success') {
                    document.getElementById('messageContent').value = '';
                    loadMessages(threadId);
                } else {
                    alert('Erreur lors de l\'envoi du message');
                }
            } catch (error) {
                console.error('[DEBUG] newMessageForm - Error:', error);
                alert('Erreur lors de l\'envoi du message');
            }
        });
    }

    // Charger la discussion au chargement de la page
    loadThread();
}); 