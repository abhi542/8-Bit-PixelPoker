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

    // Helper: Draw Chip Stack (Pixel Art)
    drawChips(x, y, amount) {
        if (amount <= 0) return;

        // Denominations: 1000=Gold, 500=Purple, 100=Black, 25=Green, 5=Red, 1=White
        // INLINED COLORS to avoid Config Caching issues
        const denoms = [
            { val: 1000, color: '#f1c40f' }, // Gold
            { val: 500, color: '#9b59b6' }, // Purple
            { val: 100, color: '#444444' }, // Light Black
            { val: 25, color: '#6abe30' }, // Green
            { val: 10, color: '#5b6ee1' }, // Blue
            { val: 5, color: '#d95763' }, // Red
            { val: 1, color: '#ffffff' }  // White
        ];

        let remaining = amount;
        let stack = [];

        // Limit total chips visually
        const MAX_VISIBLE_CHIPS = 12;
        let totalChips = 0;

        for (let d of denoms) {
            const count = Math.floor(remaining / d.val);
            if (count > 0) {
                for (let i = 0; i < count; i++) {
                    stack.push(d.color);
                    totalChips++;
                    if (totalChips >= MAX_VISIBLE_CHIPS) break;
                }
                remaining %= d.val;
            }
            if (totalChips >= MAX_VISIBLE_CHIPS) break;
        }

        // Draw Stack (Bottom Up)
        const cw = 14 * 1.5;
        const ch = 4 * 1.5;
        const step = 3 * 1.5;

        let curY = y;
        for (let color of stack) {
            // Side (Darker)
            const darker = this.shadeColor(color, -20);
            Draw.ctx.fillStyle = darker;
            Draw.ctx.fillRect(x - cw / 2, curY - ch, cw, ch);

            // Top (Surface)
            Draw.ctx.fillStyle = color;
            Draw.ctx.fillRect(x - cw / 2, curY - ch - 2, cw, 3);

            // Highlight strip on top
            Draw.ctx.fillStyle = "rgba(255,255,255,0.3)";
            Draw.ctx.fillRect(x - cw / 2 + 2, curY - ch - 2, cw - 4, 1);

            // Side Stripes (White) - aesthetic noise
            if (color !== CONFIG.colors.chipWhite) {
                Draw.ctx.fillStyle = "rgba(255,255,255,0.6)";
                Draw.ctx.fillRect(x - cw / 2 + 2, curY - 2, 2, 2);
                Draw.ctx.fillRect(x + cw / 2 - 4, curY - 2, 2, 2);
            }

            curY -= step;
        }

        Draw.text(`${amount}`, x, y + 20, 8, '#fff');
    },

    shadeColor(color, percent) {
        if (!color) return "#000";
        let f = parseInt(color.slice(1), 16), t = percent < 0 ? 0 : 255, p = percent < 0 ? percent * -1 : percent, R = f >> 16, G = f >> 8 & 0x00FF, B = f & 0x0000FF;
        return "#" + (0x1000000 + (Math.round((t - R) * p) + R) * 0x10000 + (Math.round((t - G) * p) + G) * 0x100 + (Math.round((t - B) * p) + B)).toString(16).slice(1);
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
