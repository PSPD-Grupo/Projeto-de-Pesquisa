const API_BASE = '/api';
let currentNickname = '';
let currentRoom = '';
let chatSocket = null;

const publishForm = document.getElementById('publish-form');
const desabafoText = document.getElementById('desabafo-text');
const charCount = document.getElementById('char-count');
const btnPublish = document.getElementById('btn-publish');
const btnRefresh = document.getElementById('btn-refresh');
const refreshIcon = document.getElementById('refresh-icon');
const feedContainer = document.getElementById('feed-container');
const feedStat = document.getElementById('feed-stat');

const chatJoinView = document.getElementById('chat-join-view');
const chatActiveView = document.getElementById('chat-active-view');
const chatNicknameInput = document.getElementById('chat-nickname');
const chatRoomSelect = document.getElementById('chat-room-select');
const btnJoinChat = document.getElementById('btn-join-chat');
const btnLeaveChat = document.getElementById('btn-leave-chat');
const chatMessagesContainer = document.getElementById('chat-messages-container');
const chatInputForm = document.getElementById('chat-input-form');
const chatMessageText = document.getElementById('chat-message-text');
const chatCurrentRoom = document.getElementById('chat-current-room');
const chatStatus = document.getElementById('chat-status');

window.addEventListener('DOMContentLoaded', () => {
    const randomId = Math.floor(1000 + Math.random() * 9000);
    chatNicknameInput.placeholder = `Anonimo_${randomId}`;
    fetchFeed();
});

desabafoText.addEventListener('input', () => {
    const length = desabafoText.value.length;
    charCount.textContent = length;
    
    if (length >= 280) {
        charCount.style.color = 'var(--heart-color)';
    } else if (length > 250) {
        charCount.style.color = 'orange';
    } else {
        charCount.style.color = 'var(--text-muted)';
    }
});

async function fetchFeed() {
    refreshIcon.classList.add('spin');
    feedStat.textContent = 'Buscando desabafos no mural...';
    
    try {
        const response = await fetch(`${API_BASE}/feed?quant=100`);
        if (!response.ok) {
            throw new Error(`Erro HTTP: ${response.status}`);
        }
        
        const feedItems = await response.json();
        renderFeed(feedItems);
    } catch (error) {
        console.error('Erro ao buscar feed:', error);
        feedStat.textContent = 'Erro ao carregar desabafos. Clique para tentar novamente.';
        feedContainer.innerHTML = `
            <div class="glass-card animate-fade-in" style="text-align: center; border-color: rgba(255, 62, 108, 0.2)">
                <i class="fa-solid fa-triangle-exclamation" style="font-size: 32px; color: var(--heart-color); margin-bottom: 12px;"></i>
                <p>Não foi possível conectar ao Servidor A (gRPC).</p>
                <p style="font-size: 12px; color: var(--text-muted); margin-top: 8px;">${error.message}</p>
            </div>
        `;
    } finally {
        refreshIcon.classList.remove('spin');
    }
}

function renderFeed(items) {
    feedContainer.innerHTML = '';
    
    if (items.length === 0) {
        feedStat.textContent = 'Nenhum desabafo publicado ainda.';
        feedContainer.innerHTML = `
            <div class="glass-card animate-fade-in" style="text-align: center; color: var(--text-muted)">
                <i class="fa-regular fa-comment-dots" style="font-size: 32px; margin-bottom: 12px; display: block; color: var(--primary)"></i>
                <span>O mural está silencioso. Seja o primeiro a desabafar!</span>
            </div>
        `;
        return;
    }
    
    feedStat.textContent = `${items.length} desabafo(s) encontrados.`;
    
    items.forEach((item, index) => {
        const card = document.createElement('div');
        card.classList.add('desabafo-card');
        card.style.animationDelay = `${index * 0.05}s`;
        
        const date = new Date(item.created_at);
        const formattedDate = date.toLocaleDateString('pt-BR', {
            day: '2-digit',
            month: '2-digit',
            year: 'numeric',
            hour: '2-digit',
            minute: '2-digit'
        });
        
        card.innerHTML = `
            <div class="desabafo-content">${escapeHtml(item.texto)}</div>
            <div class="desabafo-footer">
                <span class="time-stamp">
                    <i class="fa-regular fa-clock"></i> ${formattedDate}
                </span>
                <button class="btn-reaction" data-id="${item.id}">
                    <i class="fa-regular fa-heart"></i>
                    <span class="react-count">${item.reactions || 0}</span>
                </button>
            </div>
        `;
        
        const reactBtn = card.querySelector('.btn-reaction');
        reactBtn.addEventListener('click', (e) => {
            handleReact(item.id, reactBtn, e);
        });
        
        feedContainer.appendChild(card);
    });
}

publishForm.addEventListener('submit', async (e) => {
    e.preventDefault();
    const texto = desabafoText.value.trim();
    if (!texto) return;
    
    btnPublish.disabled = true;
    const btnSpan = btnPublish.querySelector('span');
    const originalText = btnSpan.textContent;
    btnSpan.textContent = 'Enviando...';
    
    try {
        const response = await fetch(`${API_BASE}/publicar`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ texto })
        });
        
        if (!response.ok) {
            const errData = await response.json();
            throw new Error(errData.detail || `Erro HTTP: ${response.status}`);
        }
        
        desabafoText.value = '';
        charCount.textContent = '0';
        charCount.style.color = 'var(--text-muted)';
        await fetchFeed();
        
    } catch (error) {
        console.error('Erro ao publicar:', error);
        alert(`Não foi possível publicar desabafo: ${error.message}`);
    } finally {
        btnPublish.disabled = false;
        btnSpan.textContent = originalText;
    }
});

async function handleReact(id, button, event) {
    button.classList.add('active');
    const heartIcon = button.querySelector('i');
    heartIcon.className = 'fa-solid fa-heart';
    
    const countSpan = button.querySelector('.react-count');
    const currentCount = parseInt(countSpan.textContent);
    countSpan.textContent = currentCount + 1;
    
    createFlyingHearts(event);
    
    try {
        const response = await fetch(`${API_BASE}/reagir`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ desabafo_id: String(id) })
        });
        
        if (!response.ok) {
            throw new Error('Falha ao salvar reação');
        }
    } catch (error) {
        console.error('Erro ao registrar reação:', error);
    }
}

btnRefresh.addEventListener('click', fetchFeed);

btnJoinChat.addEventListener('click', () => {
    let nickname = chatNicknameInput.value.trim();
    if (!nickname) {
        nickname = chatNicknameInput.placeholder;
    }
    
    const room = chatRoomSelect.value;
    currentNickname = nickname;
    currentRoom = room;
    
    joinChatRoom(room, nickname);
});

btnLeaveChat.addEventListener('click', leaveChatRoom);

function joinChatRoom(room, nickname) {
    chatMessagesContainer.innerHTML = '';
    
    // Define protocolo e endereço do WebSocket dinamicamente com base na URL atual
    const protocol = window.location.protocol === 'https:' ? 'wss' : 'ws';
    const host = window.location.host;
    const wsUrl = `${protocol}://${host}/sala/${room}?nickname=${encodeURIComponent(nickname)}`;
    
    chatStatus.textContent = 'Conectando...';
    chatStatus.className = 'desc-tag bg-live';
    
    try {
        chatSocket = new WebSocket(wsUrl);
        
        chatSocket.onopen = () => {
            chatStatus.textContent = 'Ao Vivo';
            chatStatus.classList.add('connected');
            
            chatJoinView.classList.add('hidden');
            chatActiveView.classList.remove('hidden');
            chatCurrentRoom.textContent = `#${room}`;
            chatMessageText.focus();
        };
        
        chatSocket.onmessage = (event) => {
            const msg = JSON.parse(event.data);
            if (msg.error) {
                appendSystemMessage(`Erro: ${msg.error}`);
                leaveChatRoom();
                return;
            }
            appendChatMessage(msg);
        };
        
        chatSocket.onclose = () => {
            chatStatus.textContent = 'Desconectado';
            chatStatus.className = 'desc-tag bg-live';
            
            chatJoinView.classList.remove('hidden');
            chatActiveView.classList.add('hidden');
            chatSocket = null;
        };
        
        chatSocket.onerror = (error) => {
            console.error('Erro no WebSocket:', error);
            appendSystemMessage('Erro de conexão com o chat.');
        };
        
    } catch (e) {
        console.error('Falha ao abrir WebSocket:', e);
        chatStatus.textContent = 'Erro';
        alert('Erro ao conectar ao chat.');
    }
}

function leaveChatRoom() {
    if (chatSocket) {
        chatSocket.close();
    }
}

chatInputForm.addEventListener('submit', (e) => {
    e.preventDefault();
    const text = chatMessageText.value.trim();
    if (!text || !chatSocket || chatSocket.readyState !== WebSocket.OPEN) return;
    
    chatSocket.send(JSON.stringify({ text }));
    chatMessageText.value = '';
    chatMessageText.focus();
});

function appendChatMessage(msg) {
    const bubble = document.createElement('div');
    bubble.classList.add('chat-bubble');
    const isSelf = msg.sender_nickname === currentNickname;
    
    if (msg.text === 'entrou na sala') {
        bubble.classList.add('chat-bubble-system');
        bubble.textContent = `${msg.sender_nickname} entrou na sala`;
        chatMessagesContainer.appendChild(bubble);
        chatMessagesContainer.scrollTop = chatMessagesContainer.scrollHeight;
        return;
    }
    
    if (isSelf) {
        bubble.classList.add('chat-bubble-self');
    } else {
        bubble.classList.add('chat-bubble-other');
    }
    
    const time = new Date(msg.send_at);
    const formattedTime = time.toLocaleTimeString('pt-BR', { hour: '2-digit', minute: '2-digit' });
    
    bubble.innerHTML = `
        <div class="bubble-meta">
            <span class="bubble-author">${escapeHtml(msg.sender_nickname)}</span>
            <span class="bubble-time">${formattedTime}</span>
        </div>
        <div class="bubble-text">${escapeHtml(msg.text)}</div>
    `;
    
    chatMessagesContainer.appendChild(bubble);
    chatMessagesContainer.scrollTop = chatMessagesContainer.scrollHeight;
}

function appendSystemMessage(text) {
    const bubble = document.createElement('div');
    bubble.classList.add('chat-bubble', 'chat-bubble-system');
    bubble.textContent = text;
    chatMessagesContainer.appendChild(bubble);
    chatMessagesContainer.scrollTop = chatMessagesContainer.scrollHeight;
}

// Cria efeito de corações flutuantes ao clicar no botão de reação
function createFlyingHearts(e) {
    const container = document.getElementById('hearts-container');
    const x = e.clientX;
    const y = e.clientY;
    
    for (let i = 0; i < 6; i++) {
        const heart = document.createElement('span');
        heart.classList.add('flying-heart');
        heart.innerHTML = '<i class="fa-solid fa-heart"></i>';
        
        const dx = (Math.random() - 0.5) * 120;
        const rot = (Math.random() - 0.5) * 60;
        
        heart.style.left = `${x}px`;
        heart.style.top = `${y}px`;
        heart.style.setProperty('--dx', `${dx}px`);
        heart.style.setProperty('--rot', `${rot}deg`);
        heart.style.animationDuration = `${0.6 + Math.random() * 0.4}s`;
        
        container.appendChild(heart);
        setTimeout(() => {
            heart.remove();
        }, 1000);
    }
}

function escapeHtml(text) {
    const map = {
        '&': '&amp;',
        '<': '&lt;',
        '>': '&gt;',
        '"': '&quot;',
        "'": '&#039;'
    };
    return text.replace(/[&<>"']/g, function(m) { return map[m]; });
}
