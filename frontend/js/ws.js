const WS = {
    socket: null,
    url: 'ws://localhost:8080/ws',

    connect() {
        this.socket = new WebSocket(this.url);

        this.socket.onopen = () => {
            console.log("Connected to WS");
            State.isConnected = true;
            // Auto-reconnect if we have data
            if (State.playerID && State.lobbyCode) {
                this.sendReconnect();
            }
        };

        this.socket.onclose = () => {
            console.log("Disconnected. Retrying...");
            State.isConnected = false;
            setTimeout(() => this.connect(), 2000);
        };

        this.socket.onmessage = (evt) => {
            const msg = JSON.parse(evt.data);
            this.handleMessage(msg);
        };
    },

    send(type, payload) {
        if (this.socket && this.socket.readyState === WebSocket.OPEN) {
            this.socket.send(JSON.stringify({ type, payload }));
        }
    },

    sendCreateLobby(username) {
        this.send('CREATE_LOBBY', { max_players: 6 }); // Stub max
    },

    sendJoinLobby(code, username) {
        this.send('JOIN_LOBBY', { code, username });
    },

    sendStartGame() {
        this.send('START_GAME', {});
    },

    sendAction(type, amount) {
        this.send('ACTION', { action_type: type, amount: amount });
    },

    sendReconnect() {
        this.send('RECONNECT', {
            lobby_code: State.lobbyCode,
            player_id: State.playerID
        });
    },

    handleMessage(msg) {
        console.log("RX:", msg);
        switch (msg.type) {
            case 'CREATE_LOBBY':
                // Ack. Waiting for JOIN_LOBBY usually follows, but here we just get code.
                // Actually, server sends CREATE_LOBBY back with code.
                // We typically auto-join after creating.
                // My server implementation: CreateLobby returns Code. 
                // Client must then Join.
                if (msg.payload.code) {
                    State.lobbyCode = msg.payload.code;
                    // Auto-join as "Host"
                    // We need to capture the username from input, currently stored in DOM or passed.
                    // For MVP, app.js handles this flow.
                    App.onLobbyCreated(msg.payload.code);
                }
                break;

            case 'JOIN_LOBBY':
                State.lobbyCode = msg.payload.code;
                if (msg.payload.player_id) {
                    State.setPlayerID(msg.payload.player_id);
                }
                State.lobbyState = 'JOINED';
                App.showWaitingRoom(msg.payload.code);
                break;

            case 'STATE':
                State.updateGame(msg.payload);
                State.lobbyState = 'IN_GAME';
                App.showGame();
                if (App && App.checkWinner) App.checkWinner();
                break; // Trigger render loop automatically via RAF

            case 'ERROR':
                alert("Error: " + msg.payload.message);
                break;
        }
    }
};
