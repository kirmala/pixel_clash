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
                <button id="submit-nickname">Submit</button>
                <h1>Select Game Mode</h1>
                <button id="2players-btn">2 Players</button>
                <button id="4players-btn">4 Players</button>
            </div>
        `;
    }

    renderSearchPage() {
        const playersNeeded = this.app.state.playersCount
        document.getElementById('app').innerHTML = `
            <div class="page active" id="search-page">
                <h1>Searching for Game</h1>
                <div class="search-info">
                    Players waiting: <span id="players-waiting">0</span> / 
                    <span id="players-needed">${playersNeeded}</span>
                </div>
                <button id="cancel-search">Cancel</button>
            </div>
        `;
    }

    renderGamePage() {
        document.getElementById('app').innerHTML = `
            <div class="page active" id="game-page">
                <div class="game-header">
                    <div>Time: <span id="game-timer">00:00</span></div>
                    <div>Score: <span id="game-score">0</span></div>
                </div>
                <div class="game-field" id="game-field"></div>
            </div>
        `;
        this.renderGameField(25, 25);
    }

    renderFinishedPage(data) {
        const scores = data.scores || [];
        document.getElementById('app').innerHTML = `
            <div class="page active" id="finished-page">
                <h1>Game Finished</h1>
                <div id="scores-container">
                    <h2>Final Scores</h2>
                    ${scores.map(player => `
                        <p>${player.name}: ${player.score}</p>
                    `).join('')}
                </div>
                <button id="return-home">Return to Home</button>
            </div>
        `;
    }

    renderGameField(rows, cols) {
        const gameField = document.getElementById('game-field');
        gameField.innerHTML = '';
        gameField.style.gridTemplateColumns = `repeat(${cols}, 1fr)`;
        gameField.style.gridTemplateRows = `repeat(${rows}, 1fr)`;
        
        for (let i = 0; i < rows * cols; i++) {
            const cell = document.createElement('div');
            cell.className = 'cell';
            cell.addEventListener('click', () => this.app.gameService.handleCellClick(cell));
            gameField.appendChild(cell);
        }
    }

    bindPageEvents(pageName) {
        switch (pageName) {
            case 'home':
                document.getElementById('find-game-btn').addEventListener('click', () => {
                    this.app.router.navigate('menu');
                });
                break;
                
            case 'menu':
                document.getElementById('submit-nickname').addEventListener('click', () => {
                    const nickname = document.getElementById('nickname-input').value.trim();
                    if (nickname) {
                        this.app.updateState({ nickname });
                    }
                });
                document.getElementById('2players-btn').addEventListener('click', () => {
                    this.app.gameService.startGameSearch(2);
                });
                document.getElementById('4players-btn').addEventListener('click', () => {
                    this.app.gameService.startGameSearch(4);
                });
                break;
            case 'search':
                document.getElementById('cancel-search').addEventListener('click', () => {
                    clearInterval(this.app.gameService.searchInterval);
                    this.app.router.navigate('home');
                });
                break;
                
            case 'finished':
                document.getElementById('return-home').addEventListener('click', () => {
                    this.app.router.navigate('home');
                });
                break;
        }
    }

    updateSearchStatus(waiting, needed) {
        const waitingEl = document.getElementById('players-waiting');
        const neededEl = document.getElementById('players-needed');
        if (waitingEl) waitingEl.textContent = waiting;
        if (neededEl) neededEl.textContent = needed;
    }

    updateTimer(seconds) {
        const timerEl = document.getElementById('game-timer');
        if (timerEl) {
            const minutes = Math.floor(seconds / 60).toString().padStart(2, '0');
            const secs = (seconds % 60).toString().padStart(2, '0');
            timerEl.textContent = `${minutes}:${secs}`;
        }
    }

    updateScore(score) {
        const scoreEl = document.getElementById('game-score');
        if (scoreEl) scoreEl.textContent = score;
    }
}