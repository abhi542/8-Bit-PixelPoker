const State = {
    lobbyCode: null,
    playerID: localStorage.getItem('poker_player_id') || null,
    playerSeat: -1,
    gameState: null,
    lobbyState: 'WAITING', // WAITING, JOINED, IN_GAME
    lobbyPlayers: [], // List of names in waiting room
    isConnected: false,

    reset() {
        this.gameState = null;
        this.lobbyState = 'WAITING';
        this.playerSeat = -1;
    },

    setPlayerID(id) {
        this.playerID = id;
        localStorage.setItem('poker_player_id', id);
    },

    updateGame(serverState) {
        this.gameState = serverState;

        // Find my seat
        if (this.playerID && serverState.Seats) {
            const seat = serverState.Seats.find(s => s && s.PlayerID === this.playerID);
            this.playerSeat = seat ? seat.Index : -1;
        }
    }
};
