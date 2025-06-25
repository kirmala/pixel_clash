export default class Router {
    constructor(app) {
        this.app = app;
        this.pages = {
            home: 'home-page',
            menu: 'menu-page',
            search: 'search-page',
            game: 'game-page',
            finished: 'finished-page'
        };
    }

    init() {
        window.addEventListener('popstate', this.handlePopState.bind(this));
    }

    handlePopState(event) {
        const path = window.location.pathname.substr(1) || 'home';
        this.showPage(path);
    }

    navigate(pageName, data = {}) {
        if (!this.pages[pageName]) return;
        
        history.pushState(data, '', `/${pageName}`);
        this.showPage(pageName);
    }

    showPage(pageName) {
        this.app.ui.renderPage(pageName, this.app.state);
    }
}