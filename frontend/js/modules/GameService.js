export default class GameService {
    constructor(app) {
        this.app = app;
        this.timerInterval = null;
        this.searchInterval = null;
    }

    init() {
        
    }

    startGameSearch(playersCount) {
        this.app.updateState({ playersCount, isSearching: true });
        this.app.router.navigate('search');
        
        // Simulate searching for players
        let playersWaiting = 1;
        this.app.ui.updateSearchStatus(playersWaiting, playersCount);
        
        this.searchInterval = setInterval(() => {
            playersWaiting++;
            this.app.ui.updateSearchStatus(playersWaiting, playersCount);
            
            if (playersWaiting >= playersCount) {
                clearInterval(this.searchInterval);
                this.startGame();
            }
        }, 1500);
    }

    startGame() {
        this.app.updateState({ isPlaying: true, timer: 0, score: 0 });
        this.app.router.navigate('game');
        this.app.ui.renderGameField(25, 25);
        
        // Start timer
        this.timerInterval = setInterval(() => {
            const newTimer = this.app.state.timer + 1;
            this.app.updateState({ timer: newTimer });
            this.app.ui.updateTimer(newTimer);
            
            // Simulate game end after 30 seconds
            if (newTimer >= 30) {
                this.endGame();
            }
        }, 1000);
    }

    endGame() {
        clearInterval(this.timerInterval);
        this.app.updateState({ isPlaying: false });
        
        // Generate random scores for other players
        const scores = [
            { name: this.app.state.nickname, score: this.app.state.score }
        ];
        
        for (let i = 1; i < this.app.state.playersCount; i++) {
            scores.push({ name: `Player ${i + 1}`, score: Math.floor(Math.random() * 100) });
        }
        
        this.app.router.navigate('finished');
    }

    handleCellClick(cell) {
        if (!this.app.state.isPlaying) return;
        
        cell.classList.toggle('active');
        const scoreChange = cell.classList.contains('active') ? 1 : -1;
        const newScore = this.app.state.score + scoreChange;
        this.app.updateState({ score: newScore });
        this.app.ui.updateScore(newScore);
    }
}