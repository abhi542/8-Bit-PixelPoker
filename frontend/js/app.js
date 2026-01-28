const App = {
    init() {
        Draw.init('gameCanvas');
        WS.connect();

        this.bindEvents();

        // Start Loop
        Renderer.loop();
    },

    bindEvents() {
        // Lobby
        document.getElementById('btn-create').onclick = () => {
            const user = document.getElementById('inp-name').value || "Host";
            localStorage.setItem('poker_username', user);
            WS.sendCreateLobby(user);
            this.savedUsername = user;
        };

        document.getElementById('btn-join').onclick = () => {
            const code = document.getElementById('inp-code').value.toUpperCase();
            const user = document.getElementById('inp-name').value || "Player";
            if (code.length === 4) {
                WS.sendJoinLobby(code, user);
                this.savedUsername = user;
            } else {
                alert("Invalid Code");
            }
        };

        document.getElementById('btn-start').onclick = () => {
            WS.sendStartGame();
        };

        // Dev / Testing: Reset Identity
        document.getElementById('btn-reset').onclick = () => {
            localStorage.removeItem('poker_player_id');
            localStorage.removeItem('poker_username');
            location.reload();
        };

        document.getElementById('btn-home').onclick = () => {
            location.reload();
        };

        document.getElementById('btn-home').onclick = () => {
            location.reload();
        };

        // Actions
        document.getElementById('act-fold').onclick = () => WS.sendAction('FOLD', 0);
        document.getElementById('act-check').onclick = () => WS.sendAction('CHECK', 0);
        document.getElementById('act-call').onclick = () => WS.sendAction('CALL', 0);
        document.getElementById('act-bet').onclick = () => {
            const val = parseInt(document.getElementById('bet-slider').value);
            WS.sendAction('BET', val);
        };

        document.getElementById('bet-slider').oninput = (e) => {
            document.getElementById('bet-val').innerText = e.target.value;
        };
    },

    onLobbyCreated(code) {
        // Auto join
        WS.sendJoinLobby(code, this.savedUsername || "Host");
    },

    showWaitingRoom(code) {
        document.getElementById('lobby-screen').classList.add('hidden');
        document.getElementById('waiting-screen').classList.remove('hidden');
        document.getElementById('lobby-code-disp').innerText = code;
        // In full impl, we list players here using WS updates.
        // For now, simple "Start" button becomes available.
        document.getElementById('btn-start').classList.remove('hidden');
    },

    showGame() {
        document.getElementById('waiting-screen').classList.add('hidden');
        document.getElementById('hud-screen').classList.remove('hidden');

        this.updateSidebar();
    },

    updateSidebar() {
        if (!State.gameState || State.playerSeat === -1) return;

        const mySeat = State.gameState.Seats[State.playerSeat];
        if (mySeat && mySeat.hand_desc) { // Note: JSON tag is hand_desc
            const rank = mySeat.hand_desc; // e.g. "One Pair"
            document.getElementById('my-hand-rank').innerText = rank;

            // Highlight list
            // Fix mismatch: Go returns "Pair", UI expects "One Pair"
            let visualRank = rank;
            if (rank === "Pair") visualRank = "One Pair";

            document.querySelectorAll('#rank-list li').forEach(li => {
                li.classList.remove('active');
                if (li.dataset.rank === visualRank) {
                    li.classList.add('active');
                }
            });
        }
    },

    checkWinner() {
        // Called by WS on state update
        if (!State.gameState) return;
        const g = State.gameState;

        // If Showdown and Winners exist
        // State 9 = Showdown. Note: g.State is an object { State: int, Name: string }
        if (g.State && g.State.State === 9 && g.winners && g.winners.length > 0) {
            const winnerIdx = g.winners[0]; // Just take first for now
            const seat = g.Seats[winnerIdx];
            if (seat) {
                document.getElementById('winner-name').innerText = seat.Username;
                document.getElementById('winner-hand').innerText = seat.hand_desc || "Winning Hand";

                // Show Overlay with delay to let canvas render showdown first
                setTimeout(() => {
                    document.getElementById('winner-screen').classList.remove('hidden');
                }, 1000);
            }
        } else {
            document.getElementById('winner-screen').classList.add('hidden');
        }
    }
};

window.onload = () => App.init();
