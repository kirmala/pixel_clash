import { GameTypes } from './types/Game.js';

export default class UIComponents {
    constructor(app) {
        this.app = app;
    }

    init() {
        // Initial render if needed
        this.bindGlobalEvents();
    }

    bindGlobalEvents() {
        // Events that persist across pages can go here
    }

    renderPage(pageName, data = {}) {
        const pages = {
            'home': this.renderHomePage,
            'menu': this.renderMenuPage,
            'search': this.renderSearchPage,
            'game': this.renderGamePage,
            'finished': this.renderFinishedPage
        };

        if (pages[pageName]) {
            pages[pageName].call(this, data);
            this.bindPageEvents(pageName);
        }
    }

    renderHomePage() {
        document.getElementById('app').innerHTML = `
            <div class="page active" id="home-page">
                <h1>Pixel Clash</h1>
                <button id="find-game-btn">Find Game</button>
            </div>
        `;
    }

    renderMenuPage() {
        const currentNick = this.app.state.nickname || '';
        document.getElementById('app').innerHTML = `
            <div class="page active" id="menu-page">
                <h1>Enter Your Nickname</h1>
                <input type="text" id="nickname-input" 
                       placeholder="Your nickname" 
                       value="${currentNick}">
                <h1>Select Game Mode</h1>
                <button id="2players-btn">2 Players</button>
                <button id="4players-btn">4 Players</button>
            </div>
        `;
    }

    renderSearchPage() {
        const playersNeeded = this.app.state.gameType.size
        const playersWaiting = 1
        document.getElementById('app').innerHTML = `
            <div class="page active" id="search-page">
                <h1>Searching for Game</h1>
                <div class="search-info">
                    Players waiting:<span id="players-waiting">${playersWaiting}</span> / 
                    <span id="players-needed">${playersNeeded}</span>
                </div>
                <button ID="cancel-search">Cancel</button>
            </div>
        `;
    }

    renderGamePage() {
        document.getElementById('app').innerHTML = `
            <div class="page active" id="game-page">
                <div class="game-header">
                    <div>Time: <span id="game-timer">00:00</span></div>
                </div>
                <div class="game-container">
                    <div class="game-field" id="game-field"></div>
                    <div class="game-scores" id="current-scores"></div>
                </div>
            </div>
        `;
        this.renderField(this.app.state.field)
        this.renderScores(this.app.state.participants, "current-scores")
    }

    renderFinishedPage() {
        document.getElementById('app').innerHTML = `
            <div class="page active" id="finished-page">
                <h1>Game Finished</h1>
                <div id="final-scores"></div>
                <button id="return-home">Return to Home</button>
            </div>
        `;
        this.renderScores(this.app.state.participants, "final-scores")
    }

    renderField(field) {
        const gameField = document.getElementById('game-field');
        gameField.innerHTML = '';

        // Set grID dimensions based on the field data
        const rows = field.data.length;
        const cols = rows > 0 ? field.data[0].length : 0;
        gameField.style.gridTemplateColumns = `repeat(${cols}, 1fr)`;
        gameField.style.gridTemplateRows = `repeat(${rows}, 1fr)`;

        // Render each cell
        field.data.forEach((row, y) => {
            row.forEach((cell, x) => {
                const cellElement = document.createElement('div');
                cellElement.className = 'cell';
                
                // Set cell color and size
                cellElement.style.backgroundColor = this.app.ColorService.getParticipantColor(cell.participant_id);
                cellElement.dataset.color = this.app.ColorService.getParticipantColor(cell.participant_id);
                
                // Create component size display
                if (cell.comp_size > 0) {
                    const sizeBadge = document.createElement('div');
                    sizeBadge.className = 'size-badge';
                    sizeBadge.textContent = cell.comp_size;
                    cellElement.appendChild(sizeBadge);
                }
                
                // Add click handler with coordinates
                cellElement.addEventListener('click', () => {
                    this.app.gameService.handleCellClick(cell, y, x);
                });
                
                // Store coordinates for easy reference
                cellElement.dataset.x = x;
                cellElement.dataset.y = y;
                
                gameField.appendChild(cellElement);
            });
        });
    }

    renderScores(participants, containerID) {
        const container = document.getElementById(containerID);
        if (!container) return;

        // Convert to array and sort by score (descending)
        const sortedParticipants = Object.values(participants)
            .sort((a, b) => b.score - a.score);

        container.innerHTML = `
            <div class="scoreboard">
                <h2>Player Scores</h2>
                <div class="score-list">
                    ${sortedParticipants.map((participant, index) => {
                        const color = this.app.ColorService.getParticipantColor(participant.data);
                        return `
                            <div class="score-entry ${participant.id === this.app.state.participantID ? 'current-player' : ''}"
                                style="border-left: 4px solid ${color}">
                                <span class="rank">${index + 1}.</span>
                                <span class="player-color" style="background: ${color}"></span>
                                <span class="player-name">${participant.nickname || participant.id}</span>
                                <span class="score-value">${participant.score}</span>
                            </div>
                        `;
                    }).join('')}
                </div>
            </div>
        `;
    }


    bindPageEvents(pageName) {
        switch (pageName) {
            case 'home':
                document.getElementById('find-game-btn').addEventListener('click', () => {
                    this.app.ws.connect()
                    this.app.router.navigate('menu');
                });
                break;
                
            case 'menu':
                document.getElementById('2players-btn').addEventListener('click', () => {
                    const nickname = document.getElementById('nickname-input').value.trim();
                    this.app.gameService.startGameSearch(nickname, GameTypes.TWO_PLAYERS);
                });
                document.getElementById('4players-btn').addEventListener('click', () => {
                    const nickname = document.getElementById('nickname-input').value.trim();
                    this.app.gameService.startGameSearch(nickname, GameTypes.FOUR_PLAYERS);
                });
                break;
            case 'search':
                document.getElementById('cancel-search').addEventListener('click', () => {
                    this.app.gameService.cancelSearch();
                });
                break;
                
            case 'finished':
                document.getElementById('return-home').addEventListener('click', () => {
                    this.app.router.navigate('home');
                });
                break;
        }
    }

    updateSearchStatus(waiting) {
        const waitingEl = document.getElementById('players-waiting');
        if (waitingEl) waitingEl.textContent = waiting;
    }

    updateTimer(seconds) {
        const timerEl = document.getElementById('game-timer');
        if (timerEl) {
            const minutes = Math.floor(seconds / 60).toString().padStart(2, '0');
            const secs = (seconds % 60).toString().padStart(2, '0');
            timerEl.textContent = `${minutes}:${secs}`;
        }
    }

    updateField(field) {
        if (!field?.data) return;
        const rows = field.data.length;
        const cols = rows > 0 ? field.data[0].length : 0;

        field.data.forEach((row, y) => {
            row.forEach((cell, x) => {

                const cellElement = document.querySelector(`[data-x="${x}"][data-y="${y}"]`);
                if (!cellElement) return;
                
                // Update color
                const color = this.app.ColorService.getParticipantColor(cell.participant_id);
                cellElement.style.backgroundColor = color;
                cellElement.dataset.color = color;
                
                // Update size badge
                let sizeBadge = cellElement.querySelector('.size-badge');
                if (cell.comp_size > 0) {
                    if (!sizeBadge) {
                        sizeBadge = document.createElement('div');
                        sizeBadge.className = 'size-badge';
                        cellElement.appendChild(sizeBadge);
                    }
                    sizeBadge.textContent = cell.comp_size;
                } else if (sizeBadge) {
                    cellElement.removeChild(sizeBadge);
                }
            });
        });
    }

    updateScores(participants) {
        const container = document.getElementById("current-scores");
        if (!container) return;

        if (!container || !participants) return;

        // Convert to array and sort by score (descending)
        const sortedParticipants = Object.values(participants)
            .sort((a, b) => b.score - a.score);

        container.innerHTML = `
            <div class="scoreboard">
                <h2>Player Scores</h2>
                <div class="score-list">
                    ${sortedParticipants.map((participant, index) => {
                        const color = this.app.ColorService.getParticipantColor(participant.data);
                        return `
                            <div class="score-entry ${participant.id === this.app.state.participantID ? 'current-player' : ''}"
                                style="border-left: 4px solid ${color}">
                                <span class="rank">${index + 1}.</span>
                                <span class="player-color" style="background: ${color}"></span>
                                <span class="player-name">${participant.nickname || participant.id}</span>
                                <span class="score-value">${participant.score}</span>
                            </div>
                        `;
                    }).join('')}
                </div>
            </div>
        `;
    }
}