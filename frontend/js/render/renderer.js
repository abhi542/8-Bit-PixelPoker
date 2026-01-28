const Renderer = {
    loop() {
        Draw.clear();

        if (State.lobbyState === 'IN_GAME') {
            TableRenderer.draw();
        }

        requestAnimationFrame(() => this.loop());
    }
};
