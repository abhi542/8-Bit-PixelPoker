const TableRenderer = {
    draw() {
        if (!State.gameState) return;
        const g = State.gameState;
        const ctx = Draw.ctx;

        // 1. Table Outline
        const cx = CONFIG.canvas.width / 2;
        const cy = CONFIG.canvas.height / 2 - 80; // Shift up MORE (-80)
        const tw = 600;
        const th = 300;

        Draw.ctx.beginPath();
        Draw.ctx.ellipse(cx, cy, tw / 2, th / 2, 0, 0, Math.PI * 2);
        Draw.ctx.fillStyle = CONFIG.colors.table;
        Draw.ctx.fill();
        Draw.ctx.lineWidth = 10;
        Draw.ctx.strokeStyle = CONFIG.colors.tableBorder;
        Draw.ctx.stroke();

        // 2. Community Cards
        // g.Community is Array of Card objects
        if (g.Community) {
            const cw = 40;
            const startX = cx - (g.Community.length * 50) / 2;
            g.Community.forEach((c, i) => {
                CardRenderer.drawCard(c, startX + i * 50, cy - 28, 1.5);
            });
        }

        // 3. Pot
        // Pot info in g.Betting.MainPot?
        // Or g.Pot? Wait, Table struct has Seats, Community. BettingEngine has Pots.
        // `g.Betting` might be null if game not started?
        // Check `table.go`: BetState is linked. 
        // Table JSON has `Betting`.
        // `Betting.MainPot` is int64.
        if (g.Betting) {
            let pot = g.Betting.MainPot;
            // Add side pots?
            if (g.Betting.SidePots) {
                g.Betting.SidePots.forEach(p => pot += p.Amount);
            }
            Draw.text(`POT: ${pot}`, cx, cy + 130, 16, '#fff');

            // Turn Indicator Text
            if (g.Betting.Players && g.Betting.Players[g.Betting.ActionIndex]) {
                const actorSeatIdx = g.Betting.Players[g.Betting.ActionIndex].seat_index;
                const actor = g.Seats[actorSeatIdx];
                if (actor) {
                    // Draw Bubble
                    const bx = cx;
                    const by = 80;
                    Draw.ctx.fillStyle = "rgba(0, 0, 0, 0.7)";
                    Draw.ctx.roundRect(bx - 100, by - 20, 200, 40, 10);
                    Draw.ctx.fill();
                    Draw.text(`🟢 ${actor.Username}'S TURN`, bx, by + 5, 14, '#2ecc71');
                }
            }
        }

        // SHOWDOWN ANNOUNCEMENT
        // StateShowdown = 9
        if (g.State && g.State.State === 9) {
            Draw.text(`★ SHOWDOWN ★`, cx, cy + 160, 24, '#f1c40f');
            Draw.text(`CHECK LOGS FOR WINNER`, cx, cy + 185, 12, '#95a5a6');
        }

        // 4. Seats
        // Fixed 9 positions ellipse
        for (let i = 0; i < 9; i++) {
            this.drawSeat(i, cx, cy, tw, th, g.Seats[i]);
        }
    },

    // Helper: Draw Chip Stack
    drawChips(x, y, amount) {
        if (amount <= 0) return;
        const stackHeight = Math.min(Math.floor(amount / 10), 5); // Visual limit

        for (let i = 0; i < stackHeight + 1; i++) {
            const dy = y - (i * 4);
            // Main body
            Draw.ctx.fillStyle = "#3498db";
            Draw.ctx.beginPath();
            Draw.ctx.ellipse(x, dy, 12, 6, 0, 0, Math.PI * 2);
            Draw.ctx.fill();
            Draw.ctx.strokeStyle = "#2980b9";
            Draw.ctx.lineWidth = 1;
            Draw.ctx.stroke();

            // Side/Thickness
            Draw.ctx.fillStyle = "#2980b9";
            Draw.ctx.beginPath();
            Draw.ctx.moveTo(x - 12, dy);
            Draw.ctx.lineTo(x - 12, dy + 4);
            Draw.ctx.ellipse(x, dy + 4, 12, 6, 0, 0, Math.PI);
            Draw.ctx.lineTo(x + 12, dy);
            Draw.ctx.fill();
        }

        Draw.text(`${amount}`, x, y + 20, 10, '#fff');
    },

    drawSeat(idx, cx, cy, tw, th, seat) {
        // Position logic
        const angle = (idx / 9) * Math.PI * 2 + Math.PI / 2;
        const rx = tw / 2 + 60;
        const ry = th / 2 + 60;
        const x = cx + Math.cos(angle) * rx;
        const y = cy + Math.sin(angle) * ry;

        // Base Avatar
        let color = CONFIG.colors.seat;
        let stroke = null;

        // Check current turn
        if (State.gameState.Betting && State.gameState.Betting.ActionIndex === idx) { // Note: ActionIndex refers to Players slice, NOT seat index directly usually.
            // Wait, BettingEngine has `ActionIndex`. `Players` slice.
            // Players slice is sorted?
            // `internal/poker/table.go` populates `PlayersForBet` from Seats.
            // If seats are sparse, `Players` slice [0] might be Seat 2.
            // We need to match seat.Index.
            // Actually, we highlight based on `seat.Status == Active` and if it's their turn.
            // The frontend doesn't easily know who `ActionIndex` points to unless we map it.
            // BUT: simple visual check - if `seat` exists and `seat.BetState` exists, compare logic?
            // Better: `Betting.ActionIndex` is index in `Betting.Players`. 
            // `Betting.Players` contains `SeatIndex`.
            // So: find P in Betting.Players at ActionIndex.

            const be = State.gameState.Betting;
            if (be && be.Players && be.Players[be.ActionIndex]) {
                if (be.Players[be.ActionIndex].seat_index === idx) {
                    color = CONFIG.colors.seatTurn; // Turn Highlight
                }
            }
        }

        Draw.circle(x, y, 40, color, stroke);

        if (seat) {
            // Name
            Draw.text(seat.Username.substring(0, 8), x, y - 50, 10, '#fff');
            // Chips
            Draw.text(`$${seat.Chips}`, x, y + 55, 10, '#f1c40f');

            // CURRENT BET DISPLAY
            // Where to find it? seat.BetState is linked in Go, so json structure is:
            // seat.BetState.bet_this_round (if we exported it correctly).
            // Wait, Table struct has `BetState *PlayerBetState`. When marshaled, it uses the struct tags.
            // I just added tags `bet_this_round`.

            if (seat.BetState && seat.BetState.bet_this_round > 0) {
                // Draw just inside the table
                // Vector towards center
                const dx = cx - x;
                const dy = cy - y;
                const dist = Math.sqrt(dx * dx + dy * dy);
                // Move 70px towards center
                const udx = dx / dist;
                const udy = dy / dist;
                const bx = x + udx * 90;
                const by = y + udy * 90;

                this.drawChips(bx, by, seat.BetState.bet_this_round);
            }

            // Hole Cards
            // If it's ME (State.playerSeat == idx), show face up.
            // If Showdown, show face up? (Server sends them).
            // Else face down.
            if (seat.HoleCards && seat.HoleCards.length > 0) {
                // Draw 2 cards
                CardRenderer.drawCard(seat.HoleCards[0], x - 25, y - 10, 1.0);
                CardRenderer.drawCard(seat.HoleCards[1], x + 5, y - 10, 1.0);
            } else if (seat.Status === 0) { // Active
                // Draw backs
                CardRenderer.drawCard(null, x - 25, y - 10, 1.0);
                CardRenderer.drawCard(null, x + 5, y - 10, 1.0);
            }

            // Hand Strength (if ME)
            if (idx === State.playerSeat && seat.hand_desc) {
                // Draw.text(seat.hand_desc, x, y - 70, 8, '#f39c12');
                // We update sidebar instead, but debug text here is fine.
            }

            // Turn Indicator (Already handled by circle color? Enhanced:)
            // If active ...
        } else {
            Draw.text(`OPEN`, x, y + 5, 8, '#666');
        }
    }
};
