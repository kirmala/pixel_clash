import Router from './Router.js';
import GameService from './GameService.js';
import UIComponents from './UIComponents.js';
import WebSocketService from './WebSocketService.js';
import ColorService from './ColorService.js';


export default class App {
    constructor() {
        this.router = new Router(this);
        this.gameService = new GameService(this);
        this.ui = new UIComponents(this);
        this.ws = new WebSocketService(this);
        this.ColorService = new ColorService()
        this.state = {
            nickname: '',
            status: 'menu', //serching, playing, menu
            gameType: null,
            feild: null,
            participants: null,
            LastMoveTime: 0,
            participantID: null,
            color: null,
            timer: null,
        };
    }

    init() {
        this.router.init();
        this.ui.init();
        this.gameService.init();
        this.router.navigate('home');
    }

    updateState(newState) {
        this.state = { ...this.state, ...newState };
    }
}