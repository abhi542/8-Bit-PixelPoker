const CardRenderer = {
    drawCard(card, x, y, scale = 1) {
        const w = 40 * scale;
        const h = 56 * scale;
        const pixelSize = 3 * scale; // Increased to 3 for thicker look

        // Shadow
        this.drawRect(x + pixelSize, y + pixelSize, w, h, '#000');

        // Main Border
        this.drawRect(x, y, w, h, CONFIG.colors.cardRed);

        // Inner White
        const border = pixelSize;
        this.drawRect(x + border, y + border, w - border * 2, h - border * 2, CONFIG.colors.cardWhite);

        if (!card) {
            // BACK OF CARD (Checkered)
            const innerX = x + border;
            const innerY = y + border;
            const innerW = w - border * 2;
            const innerH = h - border * 2;

            // Base Red
            Draw.ctx.fillStyle = CONFIG.colors.cardBackMain;
            Draw.ctx.fillRect(innerX, innerY, innerW, innerH);

            // Checkered Pattern
            Draw.ctx.fillStyle = CONFIG.colors.cardBackCheck;
            const checkSize = 4 * scale;
            for (let py = 0; py < innerH; py += checkSize) {
                for (let px = 0; px < innerW; px += checkSize) {
                    if ((Math.floor(px / checkSize) + Math.floor(py / checkSize)) % 2 === 0) {
                        // Check bounds to ensure we don't draw outside inner card
                        const dw = Math.min(checkSize, innerW - px);
                        const dh = Math.min(checkSize, innerH - py);
                        Draw.ctx.fillRect(innerX + px, innerY + py, dw, dh);
                    }
                }
            }
            return;
        }

        // FRONT OF CARD
        // Suits: 0=Clubs, 1=Diamonds, 2=Hearts, 3=Spades
        const colors = [CONFIG.colors.cardBlack, CONFIG.colors.cardRed, CONFIG.colors.cardRed, CONFIG.colors.cardBlack];
        const ranks = { 11: 'J', 12: 'Q', 13: 'K', 14: 'A', 10: '10' };

        const sIdx = card.suit;
        const rVal = card.rank;

        let rankStr = rVal;
        // Loose equality to catch string "10" or int 10
        if (rVal == 10 || rVal === 'T') {
            rankStr = "10";
        } else if (ranks[rVal]) {
            rankStr = ranks[rVal];
        }

        const color = colors[sIdx];

        // Draw Rank (Top Left) using Font
        // Move inward to avoid 3px border
        Draw.text(rankStr, x + 10 * scale, y + 16 * scale, 10 * scale, color, 'center');

        // Draw Suit (Middle) using Pixel Sprite
        // Center of card: x + w/2, y + h/2 + offset
        // Sprite is 7x7 pixels.
        const sprite = PixelSprites.suits[sIdx];
        const spriteScale = 2 * scale; // Make suit chunky
        const spriteSize = 7 * spriteScale;
        const sx = x + (w - spriteSize) / 2;
        const sy = y + (h - spriteSize) / 2 + 5 * scale; // Lower slightly

        Draw.pixelSprite(sprite, sx, sy, spriteScale, color);

        // Draw Upside Down Code? Maybe just simplifed for now.
    },

    // Helper for crisp rects
    drawRect(x, y, w, h, color) {
        Draw.ctx.fillStyle = color;
        Draw.ctx.fillRect(Math.floor(x), Math.floor(y), Math.floor(w), Math.floor(h));
    }
};
