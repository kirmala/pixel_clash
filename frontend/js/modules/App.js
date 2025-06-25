import Router from './Router.js';
import GameService from './GameService.js';
import UIComponents from './UIComponents.js';

export default class App {
    constructor() {
        this.router = new Router(this);
        this.gameService = new GameService(this);
        this.ui = new UIComponents(this);
        this.state = {
            nickname: '',
            playersCount: null,
            isPlaying: false
        };
    }

    init() {
        this.router.init();
        this.ui.init();
        this.gameService.init();
        
        // Initial page
        this.router.navigate('home');
    }

    updateState(newState) {
        this.state = { ...this.state, ...newState };
    }
}