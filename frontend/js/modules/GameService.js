export default class GameService {
    constructor(app) {
        this.app = app;
        this.timerInterval = null;
    }

    init() {
        
    }


    isOnCooldown() {
        const { lastMoveTime } = this.app.state;
        return Date.now() - lastMoveTime < this.app.state.gameType.Cooldown;
    }

    async startGameSearch(nickname, gameType) {
        try {
            const response = await this.app.ws.findGame(nickname, gameType);
            if (response.status === "error") {
                console.error("start game search failed:", response.data);
            } else if (response.status === "success") {
                const color = this.app.ColorService.getParticipantColor(response.data.participant_id)
                this.app.updateState({
                    status: "searching",
                    nickname: nickname,
                    gameType: gameType,
                    participantID: response.data.participant_id,
                    color: color,
                });
                this.app.router.navigate('search');
            }
        } catch (error) {
            console.error("Find game error:", error);
        }
    }

    async handleCellClick(cell, y, x) {
        if (this.app.state.status != "playing") return;

        if (this.isOnCooldown()) {
            return;
        }

        if (cell.participant_id === this.app.state.participantID) {
            return;
        }
        
        try {
            const response = await this.app.ws.makeMove(y, x);
            if (response.status === "error") {
                console.error("Move failed:", response.data);
            } else {
                this.app.updateState({ lastMoveTime: Date.now() });
            }
        } catch (error) {
            console.error("Move error:", error);
        }
    }

    async cancelSearch() {
        try {
            const response = await this.app.ws.stopSearching();
            if (response.status === "error") {
                if (response.data === "Connection canceled") {
                    this.app.router.navigate("home")
                }
                console.error("Stop searching failed:", response.data);
            }
        } catch (error) {
            console.error("Stop searching error:", error);
        }
    }

    handleGameStart(game_data) {
        this.app.updateState({
            status: "playing",
            field: game_data.field,
            participants: game_data.participants.data,
            timer: this.app.state.gameType.time
        });
        this.timerInterval = setInterval(() => {
            const newTimer = this.app.state.timer - 1;
            this.app.updateState({ timer: newTimer });
            this.app.ui.updateTimer(newTimer);
        }, 1000);
        this.app.router.navigate('game');
    }

    handleGameFinish(finish_data) {
        clearInterval(this.timerInterval);
        this.app.updateState({
            status: "menu",
            participants: finish_data.participants.data
        });
        this.app.router.navigate("finished")
    }

    handleOpponentMove(move_data) {
        this.app.updateState({
            field: move_data.field,
            participants: move_data.participants.data
        });
        this.app.ui.updateField(move_data.field);
        this.app.ui.updateScores(move_data.participants.data);
    }
}