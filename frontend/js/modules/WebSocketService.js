export default class WebSocketService {
    constructor(app) {
        this.app = app;
        this.socket = null;
        this.pendingRequests = new Map();
        // this.reconnectAttempts = 0;
        // this.maxReconnectAttempts = 5;
    }

    connect() {
        // Adjust the WebSocket URL to your server
        URL = "ws://localhost:8080/user/join"
        this.socket = new WebSocket(URL); 

        this.socket.onopen = () => {
            console.log('WebSocket connected');
            //this.reconnectAttempts = 0;
            this.app.updateState({ isConnected: true });
            
            // // Rejoin game if connection was lost during gameplay
            // if (this.app.state.gameID) {
            //     this.send('rejoin', { 
            //         gameID: this.app.state.gameID,
            //         playerID: this.app.state.playerID 
            //     });
            // }
        };

        this.socket.onmessage = (event) => {
            const message = JSON.parse(event.data);
            
            if (message.type === "event") {
                this.handleEvent(message.data);
            } else if (message.type === "response") {
                this.handleResponse(message.data);
            }
        };

        this.socket.onclose = () => {
            console.log('WebSocket disconnected');
            this.app.updateState({ isConnected: false });
            //this.handleReconnect();
            this.pendingRequests.forEach((request, ID) => {
                this.handleResponse({
                    type: request.type,
                    id: ID,
                    status: "error",
                    data: "Connection closed"
                });
            });
            this.app.ws.disconnect()
            this.app.ColorService.reset()
            this.app.router.navigate('home')
            clearInterval(this.app.gameService);
        };

        this.socket.onerror = (error) => {
            console.error('WebSocket error:', error);
        };
    }

    handleEvent(event) {
        switch (event.type) {
            case 'game_start':
                this.app.gameService.handleGameStart(event.data);
                break;
            case 'game_finish':
                this.app.gameService.handleGameFinish(event.data);
                break;
            case 'player_move':
                this.app.gameService.handleOpponentMove(event.data);
                break;
            case 'waiting_change':
                this.app.ui.updateSearchStatus(event.data.waiting);
                break;
        }
    }

    handleResponse(response) {
        const callback = this.pendingRequests.get(response.id);
        if (callback) {
            callback(response);
            this.pendingRequests.delete(response.id);
        }
    }


    sendRequest(type, data, callback) {
        // // Validate WebSocket connection
        // if (this.socket?.readyState !== WebSocket.OPEN) {
        //     callback?.({
        //         Type: "type",
        //         Status: "error",
        //         Data: { Message: "WebSocket not connected" }
        //     });
        //     return;
        // }

        const requestId = crypto.randomUUID();
        const request = {
            type: type,
            id: requestId,
            data: data
        };

        // Set timeout handler
        // const timeout = setTimeout(() => {
        //     this.pendingRequests.delete(requestId);
        //     callback?.({
        //         type: request.Type,
        //         id: request.Id,
        //         status: "error",
        //         data: "Request timed out"
        //     });
        // }, 10000); // 10 seconds

        // Store callback
        this.pendingRequests.set(requestId, (response) => {
            //clearTimeout(timeout);
            callback?.(response);
        });

        // ACTUALLY SEND THE REQUEST
        try {
            this.socket.send(JSON.stringify(request));
        } catch (err) {
            //clearTimeout(timeout);
            this.pendingRequests.delete(requestId);
            callback?.({
                type: request.type,
                id: request.id,
                Status: "error",
                Data: "Failed to send request"
            });
        }
    }

    findGame(nickname, gameType) {
        return new Promise((resolve) => {
            this.sendRequest('find_game', { 
                nickname: nickname,
                game_type: gameType
            }, (response) => {
                resolve(response);
            });
        });
    }

    makeMove(y, x) {
        return new Promise((resolve) => {
            this.sendRequest('move', {
                x: x,
                y: y
            }, (response) => {
                resolve(response);
            });
        });
    }

    stopSearching() {
        return new Promise((resolve) => {
            this.sendRequest('stop_searching', {}, (response) => {
                resolve(response);
            });
        });
    }

    // handleReconnect() {
    //     if (this.reconnectAttempts < this.maxReconnectAttempts) {
    //         const delay = Math.min(1000 * (2 ** this.reconnectAttempts), 10000);
    //         console.log(`Reconnecting in ${delay}ms...`);
            
    //         setTimeout(() => {
    //             this.reconnectAttempts++;
    //             this.connect();
    //         }, delay);
    //     }
    // }

    disconnect() {
        if (this.socket) {
            this.socket.close();
        }
    }
}